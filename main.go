package main

import (
	"awesomeProject/database"
	"awesomeProject/handler/task"
	"awesomeProject/repository"
	"github.com/gin-gonic/gin"
	"log"
)

/*func seq(result chan bool) {
	for i := 1; i < 100; i++ {
		println(i)
		//time.Sleep(time.Second * 1)

	}
	result <- true
}

func main() {
	result := make(chan bool, 1)
	go seq(result)
	<-result
	log.Println("success")

}*/

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

	th := task.Task{TasksRepository: repo}
	//http.Handle("/tasks/", th)

	//err = http.ListenAndServe(":8080", nil)
	//if err != nil {
	//	log.Fatal("ListenAndServe: ", err)
	//}
	r := gin.Default()
	r.GET("/tasks", th.List)
	r.GET("/tasks/:id", th.DisplayTasks)
	r.POST("/tasks", th.AddTask)
	r.PATCH("/tasks/:id", th.Update)
	r.DELETE("/tasks/:id", th.DeleteTask)
	r.Run(":8080")

}
