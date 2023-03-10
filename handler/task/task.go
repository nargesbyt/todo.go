package task

import (
	"awesomeProject/handler"
	"awesomeProject/repository"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type Task struct {
	TasksRepository repository.Tasks
}

func (t Task) List(c *gin.Context) {
	tasks, err := t.TasksRepository.Find(c.Query("title"), c.Query("status"))
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, handler.NewProblem(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
		return
	}
	c.JSON(http.StatusOK, tasks)

}
func (t Task) DisplayTasks(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, handler.NewProblem(http.StatusBadRequest, "invalid task id"))
		return
	}
	task, err := t.TasksRepository.DisplayTask(id)

	if err != nil {
		if err == repository.ErrTaskNotFound {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusNotFound, handler.NewProblem(http.StatusNotFound, "Task not found"))
			return
		}
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	resp := Response{}
	resp.FromEntity(task)
	c.JSON(http.StatusOK, resp)

}
func (t Task) AddTask(c *gin.Context) {
	//task := entity.Task{}
	cRequest := CreateRequest{}
	if err := c.BindJSON(&cRequest); err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	task, err := t.TasksRepository.Create(cRequest.Title)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(500)
		return
	}
	resp := Response{}
	resp.FromEntity(task)
	c.JSON(http.StatusCreated, resp)

}

func (t Task) DeleteTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = t.TasksRepository.Delete(id)
	if err != nil {
		if err == repository.ErrTaskNotFound {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusNotFound, handler.NewProblem(http.StatusNotFound, "Task not found"))
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusAccepted)
}

/*reg, err := regexp.Compile(`/tasks/(\d+)`)
if err != nil {
	log.Println(err)
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	return
}
matches := reg.FindStringSubmatch(r.URL.Path)
if len(matches) == 0 {
	log.Println(err)
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	return
}
id, err := strconv.ParseInt(matches[1], 10, 64)
delResult, err := t.TasksRepository.Delete(id)
if err != nil {
	log.Println(err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
w.Write(delResult)*/

func (t Task) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	uRequest := UpdateRequest{}
	if err := c.BindJSON(&uRequest); err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	resp := Response{}
	updateResult, err := t.TasksRepository.Update(id, uRequest.Title, uRequest.Status)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	resp.FromEntity(updateResult)
	c.JSON(http.StatusOK, resp)
}
