package main

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/nargesbyt/todo.go/database"
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

func AfterAuth(tokenRepository repository.Tokens) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenId, _ := c.Get("tokenId")
		token, err := tokenRepository.Get(tokenId.(int64))
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var ExpiredAt time.Time
		if token.ExpiredAt.Valid {
			ExpiredAt = token.ExpiredAt.Time
		}
		tokenRepository.Update(tokenId.(int64), token.Title, ExpiredAt, time.Now(), token.Active)
	}
}

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
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		splitedUserPass := strings.Split(string(userPass), ":")
		user, err := usersRepository.GetUserByUsername(splitedUserPass[0])
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		if strings.HasPrefix(splitedUserPass[1], "todo_pat") {
			tokens, err := tokensRepository.List("", user.ID)
			if err != nil {
				c.AbortWithError(http.StatusUnauthorized, err)
				return
			}
			numTokens := len(tokens)
			for _, token := range tokens {
				err = token.VerifyToken(splitedUserPass[1])
				numTokens--
				if err != nil {
					if numTokens != 0 {
						continue
					}
					c.AbortWithError(http.StatusUnauthorized, err)
					return
				}
				/*if time.Now().After(token.ExpiredAt) {
					c.AbortWithError(http.StatusUnauthorized, err)
					return
				}*/
				if token.Active == 0 {
					c.AbortWithError(http.StatusUnauthorized, err)
					return
				}
				c.Set("tokenId", token.ID)
				c.Set("userId", user.ID)
				c.Next()

				return
			}

		}
		if err := user.CheckPassword(splitedUserPass[1]); err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		c.Set("userId", user.ID)
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

	userRepository, err := repository.NewUser(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to initialize the users repository")
	}

	tRepository, err := repository.NewToken(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to initialize the tokens repository")
	}

	th := task.Task{TasksRepository: repo}
	uh := user.User{UsersRepository: userRepository}
	toh := token.Token{TokenRepository: tRepository}

	r := gin.Default()
	r.GET("/tasks", BasicAuth(userRepository, tRepository), th.List, AfterAuth(tRepository))
	r.GET("/tasks/:id", BasicAuth(userRepository, tRepository), th.DisplayTasks, AfterAuth(tRepository))
	r.POST("/tasks", BasicAuth(userRepository, tRepository), th.AddTask, AfterAuth(tRepository))
	r.PATCH("/tasks/:id", BasicAuth(userRepository, tRepository), th.Update, AfterAuth(tRepository))
	r.DELETE("/tasks/:id", BasicAuth(userRepository, tRepository), th.DeleteTask, AfterAuth(tRepository))

	r.POST("/users", uh.Create)
	r.GET("/users", uh.ListUsers)
	r.GET("/users/:id", uh.Get)
	r.PATCH("/users/:id", uh.UpdateUsers)
	r.DELETE("/users/:id", uh.Delete)

	r.POST("/tokens", BasicAuth(userRepository, tRepository), toh.Create, AfterAuth(tRepository))
	r.GET("/tokens/:id", BasicAuth(userRepository, tRepository), toh.Get, AfterAuth(tRepository))
	r.GET("/tokens", BasicAuth(userRepository, tRepository), toh.List, AfterAuth(tRepository))
	r.PATCH("/tokens/:id", BasicAuth(userRepository, tRepository), toh.Update, AfterAuth(tRepository))
	r.DELETE("/tokens/:id", BasicAuth(userRepository, tRepository), toh.Delete, AfterAuth(tRepository))

	err = r.Run(":8080")
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to run HTTP server")
	}
}
