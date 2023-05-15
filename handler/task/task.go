package task

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/nargesbyt/todo.go/handler"
	"github.com/nargesbyt/todo.go/internal/dto"
	"github.com/nargesbyt/todo.go/repository"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

type Task struct {
	TasksRepository repository.Tasks
}

func (t Task) List(c *gin.Context) {

	pageNumber, _ := strconv.Atoi(c.Query(jsonapi.QueryParamPageNumber))
	limit, _ := strconv.Atoi(c.Query(jsonapi.QueryParamPageLimit))
	userId, _ := c.Get("userId")
	tasks, err := t.TasksRepository.Find(c.Query("title"), c.Query("status"), userId.(int64), pageNumber, limit)
	fmt.Println("tasks are: ", tasks)
	if err != nil {
		log.Error().Stack().Err(err).Msg("internal server error")

		c.AbortWithStatusJSON(http.StatusInternalServerError, handler.NewProblem(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
		return
	}
	var dtoTasks []*dto.Task
	for _, task := range tasks {
		if task.UserID != userId {
			continue
		}
		resp := dto.Task{}
		resp.FromEntity(*task)
		dtoTasks = append(dtoTasks, &resp)

	}
	c.Header("Content-Type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayload(c.Writer, dtoTasks); err != nil {
		log.Fatal().Err(err).Msg("can not respond")
	}

}
func (t Task) DisplayTasks(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, handler.NewProblem(http.StatusBadRequest, "invalid task id"))
		return
	}
	userId, _ := c.Get("userId")
	task, err := t.TasksRepository.Get(id)

	if err != nil {
		if err == repository.ErrTaskNotFound {
			log.Error().Stack().Err(err).Msg("task not found")
			c.AbortWithStatusJSON(http.StatusNotFound, handler.NewProblem(http.StatusNotFound, "Task not found"))

			return
		}

		log.Error().Stack().Err(err).Msg("internal server error")
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}
	if task.UserID != userId {
		log.Error().Stack().Err(err).Msg("unauthorized")
		c.AbortWithStatus(http.StatusUnauthorized)

		return
	}

	resp := dto.Task{}
	resp.FromEntity(task)
	c.Header("Content-Type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayload(c.Writer, &resp); err != nil {
		log.Fatal().Err(err).Msg("can not respond")
	}
}

func (t Task) AddTask(c *gin.Context) {
	cRequest := dto.TaskCreateRequest{}
	if err := c.BindJSON(&cRequest); err != nil {
		log.Error().Stack().Err(err).Msg("unprocessable entity")
		c.AbortWithStatus(http.StatusUnprocessableEntity)

		return
	}
	userId, _ := c.Get("userId")
	task, err := t.TasksRepository.Create(cRequest.Title, userId.(int64))
	if err != nil {
		log.Error().Stack().Err(err).Msg("internal server error")
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}
	resp := dto.Task{}
	resp.FromEntity(task)
	c.Header("Content-Type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayload(c.Writer, &resp); err != nil {
		log.Fatal().Err(err).Msg("can not respond")
	}
}

func (t Task) DeleteTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	task, err := t.TasksRepository.Get(id)
	if err != nil {
		if err == repository.ErrTaskNotFound {
			log.Error().Stack().Err(err).Msg("task not found")
			c.AbortWithStatusJSON(http.StatusNotFound, handler.NewProblem(http.StatusNotFound, "Task not found"))

			return
		}

		log.Error().Stack().Err(err).Msg("internal server error")
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	userId, _ := c.Get("userId")
	if task.UserID != userId {
		log.Error().Stack().Err(err).Msg("unauthorized")
		c.AbortWithStatus(http.StatusUnauthorized)

		return
	}
	err = t.TasksRepository.Delete(id)
	if err != nil {
		/*if err == repository.ErrTaskNotFound {
			log.Error().Stack().Err(err).Msg("task not found")
			c.AbortWithStatusJSON(http.StatusNotFound, handler.NewProblem(http.StatusNotFound, "Task not found"))

			return
		}*/

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusAccepted)
}

func (t Task) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	task, err := t.TasksRepository.Get(id)
	if err != nil {
		if err == repository.ErrTaskNotFound {
			log.Error().Stack().Err(err).Msg("task not found")
			c.AbortWithStatusJSON(http.StatusNotFound, handler.NewProblem(http.StatusNotFound, "Task not found"))

			return
		}

		log.Error().Stack().Err(err).Msg("internal server error")
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	userId, _ := c.Get("userId")
	if task.UserID != userId {
		log.Error().Stack().Err(err).Msg("unauthorized")
		c.AbortWithStatus(http.StatusUnauthorized)

		return
	}

	uRequest := dto.TaskUpdateRequest{}
	if err := c.BindJSON(&uRequest); err != nil {
		log.Error().Stack().Err(err).Msg("unprocessable entity")
		c.AbortWithStatus(http.StatusUnprocessableEntity)

		return
	}

	resp := dto.Task{}
	updateResult, err := t.TasksRepository.Update(id, uRequest.Title, uRequest.Status)
	if err != nil {
		if err == repository.ErrUnauthorized {
			log.Error().Stack().Err(err).Msg("unauthorized")
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		log.Error().Stack().Err(err).Msg("internal server error")
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}
	resp.FromEntity(updateResult)
	c.Header("Content-Type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayload(c.Writer, &resp); err != nil {
		log.Fatal().Err(err).Msg("can not respond")
	}
}
