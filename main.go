package main

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/nargesbyt/todo.go/database"
	"github.com/nargesbyt/todo.go/handler/task"
	"github.com/nargesbyt/todo.go/handler/user"
	"github.com/nargesbyt/todo.go/repository"
	"log"
	"net/http"
	"strings"
)

func BasicAuth(usersRepository repository.Users) gin.HandlerFunc {
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
		if err := user.CheckPassword(splitedUserPass[1]); err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		c.Set("userId", user.ID)
		//c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "user_id", user.ID))
		c.Next()

	}
}

func main() {
	db, err := database.NewSqlite("todo.db")

	if err != nil {
		log.Fatal(err)
	}
	repo, err := repository.NewTasks(db)
	if err != nil {
		log.Fatal("Init tasks table ", err)
	}
	//repo.Init()
	userRepository, err := repository.NewUser(db)
	if err != nil {
		log.Fatal("Init users table", err)
	}
	th := task.Task{TasksRepository: repo}
	uh := user.User{UsersRepository: userRepository}

	r := gin.Default()
	r.GET("/tasks", th.List)
	r.GET("/tasks/:id", th.DisplayTasks)
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
