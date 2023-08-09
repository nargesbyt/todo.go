package main

import (
	"context"
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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"strings"
	"time"
)

func BasicAuth(usersRepository repository.Users, tokensRepository repository.Tokens) gin.HandlerFunc {
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
		} else {
			log.Fatal().Err(err).Msg("unexpected error")
			return
		}
	}

	provider, err := oidc.NewProvider(context.Background(), "https://gitlab.com")
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to initialize an OIDC provider")
		return
	}

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

	db, err := database.NewSqlite(viper.GetString("database.dsn"))
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to initialize database connection")
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

	ah := oauth.OAuth{OAuth2Config: oauth2Config}
	th := task.Task{TasksRepository: repo}
	uh := user.User{UsersRepository: userRepository}
	toh := token.Token{TokenRepository: tRepository}

	r := gin.Default()
	r.GET("/oauth", ah.Get)

	r.POST("/tasks", BasicAuth(userRepository, tRepository), th.Create)
	r.GET("/tasks", BasicAuth(userRepository, tRepository), th.List)
	r.GET("/tasks/:id", BasicAuth(userRepository, tRepository), th.Get)
	r.PATCH("/tasks/:id", BasicAuth(userRepository, tRepository), th.Update)
	r.DELETE("/tasks/:id", BasicAuth(userRepository, tRepository), th.Delete)

	r.POST("/users", uh.Create)
	r.GET("/users", uh.List)
	r.GET("/users/:id", uh.Get)
	r.PATCH("/users/:id", uh.Update)
	r.DELETE("/users/:id", uh.Delete)

	r.POST("/tokens", BasicAuth(userRepository, tRepository), toh.Create)
	r.GET("/tokens", BasicAuth(userRepository, tRepository), toh.List)
	r.GET("/tokens/:id", BasicAuth(userRepository, tRepository), toh.Get)
	r.PATCH("/tokens/:id", BasicAuth(userRepository, tRepository), toh.Update)
	r.DELETE("/tokens/:id", BasicAuth(userRepository, tRepository), toh.Delete)

	err = r.Run(viper.GetString("port"))
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to run HTTP server")
	}
}
