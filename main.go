package main

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/nargesbyt/todo.go/database"
	"github.com/nargesbyt/todo.go/entity"
	"github.com/nargesbyt/todo.go/handler/task"
	"github.com/nargesbyt/todo.go/handler/token"
	"github.com/nargesbyt/todo.go/handler/user"
	"github.com/nargesbyt/todo.go/repository"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
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

				verifiedToken = t
				break
			}

			if verifiedToken == nil {
				c.AbortWithStatus(http.StatusUnauthorized)

				return
			}

			var ExpiredAt time.Time
			if verifiedToken.ExpiredAt.Valid {
				ExpiredAt = verifiedToken.ExpiredAt.Time
			}
			tokensRepository.Update(verifiedToken.ID, verifiedToken.Title, ExpiredAt, time.Now(), verifiedToken.Active)

			c.Set("userId", userEntity.ID)
			c.Next()
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
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = zerolog.New(os.Stderr).
		Level(zerolog.ErrorLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	db, err := database.NewSqlite("todo.db")

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

	th := task.Task{TasksRepository: repo}
	uh := user.User{UsersRepository: userRepository}
	toh := token.Token{TokenRepository: tRepository}

	r := gin.Default()
	r.GET("/tasks", BasicAuth(userRepository, tRepository), th.List)
	r.GET("/tasks/:id", BasicAuth(userRepository, tRepository), th.DisplayTasks)
	r.POST("/tasks", BasicAuth(userRepository, tRepository), th.AddTask)
	r.PATCH("/tasks/:id", BasicAuth(userRepository, tRepository), th.Update)
	r.DELETE("/tasks/:id", BasicAuth(userRepository, tRepository), th.DeleteTask)

	r.POST("/users", uh.Create)
	r.GET("/users", uh.ListUsers)
	r.GET("/users/:id", uh.Get)
	r.PATCH("/users/:id", uh.UpdateUsers)
	r.DELETE("/users/:id", uh.Delete)

	r.POST("/tokens", BasicAuth(userRepository, tRepository), toh.Create)
	r.GET("/tokens/:id", BasicAuth(userRepository, tRepository), toh.Get)
	r.GET("/tokens", BasicAuth(userRepository, tRepository), toh.List)
	r.PATCH("/tokens/:id", BasicAuth(userRepository, tRepository), toh.Update)
	r.DELETE("/tokens/:id", BasicAuth(userRepository, tRepository), toh.Delete)

	err = r.Run(":8080")
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to run HTTP server")
	}
}
