package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nargesbyt/todo.go/database"
	"github.com/nargesbyt/todo.go/handler/task"
	"github.com/nargesbyt/todo.go/handler/user"
	"github.com/nargesbyt/todo.go/repository"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strings"
)

func GenerateToken(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal().Err(err)
	}
	fmt.Println("Hash to store:", string(hash))

	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil))
}

func AccessTokenAuth(usersRepository repository.Users) gin.HandlerFunc {
	return func(c *gin.Context) {
		authz := c.GetHeader("Authorization")

		splits := strings.Split(authz, " ")

		if splits[0] != "Basic" {
			c.Next()
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
		//userId := user.ID
		if err := user.CheckPassword(splitedUserPass[1]); err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		c.Set("userId", user.ID)
		token := GenerateToken(splitedUserPass[1])
		c.JSON(http.StatusOK, token)

	}
}
func BasicAuth(usersRepository repository.Users) gin.HandlerFunc {
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
	th := task.Task{TasksRepository: repo}
	uh := user.User{UsersRepository: userRepository}

	r := gin.Default()
	r.GET("/tasks", BasicAuth(userRepository), th.List)
	r.GET("/tasks/:id", BasicAuth(userRepository), th.DisplayTasks)
	r.POST("/tasks", BasicAuth(userRepository), th.AddTask)
	r.PATCH("/tasks/:id", BasicAuth(userRepository), th.Update)
	r.DELETE("/tasks/:id", BasicAuth(userRepository), th.DeleteTask)

	r.POST("/users", uh.Create)
	r.GET("/users", uh.ListUsers)
	r.GET("/users/:id", uh.Get)
	r.PATCH("/users/:id", uh.UpdateUsers)
	r.DELETE("/users/:id", uh.Delete)
	r.Run(":8080")
}
