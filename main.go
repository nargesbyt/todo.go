package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/nargesbyt/todo.go/database"
	"github.com/nargesbyt/todo.go/entity"
	"github.com/nargesbyt/todo.go/handler/oauth"
	"github.com/nargesbyt/todo.go/handler/task"
	"github.com/nargesbyt/todo.go/handler/token"
	"github.com/nargesbyt/todo.go/handler/user"
	"github.com/nargesbyt/todo.go/repository"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strings"
	"time"
)
//BasicAuth authenticates users that want to send a request to server
func BasicAuth(usersRepository repository.Users, tokensRepository repository.Tokens, oidcProvider *oidc.Provider) gin.HandlerFunc {
	return func(c *gin.Context) {
		authz := c.GetHeader("Authorization")
		if authz == "" {
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		splits := strings.Split(authz, " ")
		if len(splits) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}
		if splits[0] == "Bearer" {
			var verifier = oidcProvider.Verifier(&oidc.Config{ClientID: viper.GetString("oauth.client_id")})
			_, err := verifier.Verify(context.Background(), splits[1])
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)

				return
			}
			c.Next()
			return
		}
		if splits[0] != "Basic" {
			c.AbortWithStatus(http.StatusUnauthorized)

			return

		}

		userPass, err := base64.StdEncoding.DecodeString(splits[1])
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		splitedUserPass := strings.Split(string(userPass), ":")
		userEntity, err := usersRepository.GetUserByUsername(splitedUserPass[0])
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		if strings.HasPrefix(splitedUserPass[1], "todo_pat_") {
			tokens, err := tokensRepository.GetTokensByUserID(userEntity.ID)
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)

				return
			}

			var verifiedToken *entity.Token
			for _, t := range tokens {
				err = t.VerifyToken(splitedUserPass[1])
				if err != nil {
					continue
				}

				if t.Active == 0 {
					break
				}
				if t.ExpiredAt.Valid {
					if t.ExpiredAt.Time.Before(time.Now()) {
						break
					}
				}

				verifiedToken = t
				break
			}

			if verifiedToken == nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				log.Fatal().Err(errors.New("token is expired")).Msg("token is expired")

				return
			}

			var ExpiredAt time.Time
			if verifiedToken.ExpiredAt.Valid {
				ExpiredAt = verifiedToken.ExpiredAt.Time
			}

			_, err = tokensRepository.Update(verifiedToken.ID, verifiedToken.Title, ExpiredAt, time.Now(), verifiedToken.Active)
			if err != nil {
				log.Fatal().Err(err).Msg("Unable to update the last_used column.")
			}

			c.Set("userId", userEntity.ID)
			c.Next()
			return
		}

		if err := userEntity.CheckPassword(splitedUserPass[1]); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		c.Set("userId", userEntity.ID)
		c.Next()
	}
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal().Err(err).Msg("unable to read config file")
			return
		} 
		log.Fatal().Err(err).Msg("unexpected error")
		return
	
	}

	provider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to initialize an OIDC provider")
		return
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})

	// Configure an OpenID Connect aware OAuth2 client.
	oauth2Config := oauth2.Config{
		ClientID:     viper.GetString("oauth.client_id"),
		ClientSecret: viper.GetString("oauth.client_secret"),
		RedirectURL:  "http://localhost:8080/oauth/callback",

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	logLevel, err := zerolog.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Unable to parse log level")

		return
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = zerolog.New(os.Stderr).
		Level(logLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	var db *gorm.DB

	switch viper.GetString("database.driver") {
	case "sqlite":
		db, err = database.NewSqlite(viper.GetString("database.dsn"))
		if err != nil {
			log.Fatal().Err(err).Msg("Unable to initialize database connection")
		}
	case "postgres":
		sqlDB, err := sql.Open("pgx", viper.GetString("database.dsn"))
		if err != nil {
			log.Fatal().Err(err).Msg("Unable to initialize database connection")
		}
		db, err = database.NewPostgres(sqlDB)
		if err != nil {
			log.Fatal().Err(err).Msg("Unable to initialize database connection")
		}
	default:
		log.Fatal().Msg("driver not found")

	}

	repo, err := repository.NewTasks(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to initialize the tasks repository")
	}

	userRepository, err := repository.NewUsers(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to initialize the users repository")
	}

	tRepository, err := repository.NewTokens(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to initialize the tokens repository")
	}

	ah := oauth.OAuth{OAuth2Config: oauth2Config, RedisClient: redisClient}

	th := task.Task{TasksRepository: repo}
	uh := user.User{UsersRepository: userRepository}
	toh := token.Token{TokenRepository: tRepository}

	r := gin.Default()
	r.GET("/oauth", ah.Get)
	r.GET("/oauth/callback", ah.Callback)

	r.POST("/tasks", BasicAuth(userRepository, tRepository, provider), th.Create)
	r.GET("/tasks", BasicAuth(userRepository, tRepository, provider), th.List)
	r.GET("/tasks/:id", BasicAuth(userRepository, tRepository, provider), th.Get)
	r.PATCH("/tasks/:id", BasicAuth(userRepository, tRepository, provider), th.Update)
	r.DELETE("/tasks/:id", BasicAuth(userRepository, tRepository, provider), th.Delete)

	r.POST("/users", uh.Create)
	r.GET("/users", uh.List)
	r.GET("/users/:id", uh.Get)
	r.PATCH("/users/:id", uh.Update)
	r.DELETE("/users/:id", uh.Delete)

	r.POST("/tokens", BasicAuth(userRepository, tRepository, provider), toh.Create)
	r.GET("/tokens", BasicAuth(userRepository, tRepository, provider), toh.List)
	r.GET("/tokens/:id", BasicAuth(userRepository, tRepository, provider), toh.Get)
	r.PATCH("/tokens/:id", BasicAuth(userRepository, tRepository, provider), toh.Update)
	r.DELETE("/tokens/:id", BasicAuth(userRepository, tRepository, provider), toh.Delete)

	err = r.Run(viper.GetString("port"))
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to run HTTP server")
	}
}
