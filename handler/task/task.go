package task

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/nargesbyt/todo.go/handler"
	"github.com/nargesbyt/todo.go/internal/dto"
	"github.com/nargesbyt/todo.go/repository"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"net/http"
	"os"
	"strconv"
)

type Task struct {
	TasksRepository repository.Tasks
}

/*func jsonapiMeta() *jsonapi.Meta{
	return &jsonapi.Meta{
		"details":
			"totalPage": totalPages
	}


}*/

func (t Task) List(c *gin.Context) {

	pageNumber, _ := strconv.Atoi(c.Query(jsonapi.QueryParamPageNumber))
	limit, _ := strconv.Atoi(c.Query(jsonapi.QueryParamPageLimit))
	userId, _ := c.Get("userId")
	tasks, err := t.TasksRepository.Find(c.Query("title"), c.Query("status"), userId.(int64), pageNumber, limit)
	fmt.Println("tasks are: ", tasks)
	if err != nil {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		logger := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
		logger.Error().Stack().Err(err).Msg("internal server error")
		//log.Println(err)
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
	task, err := t.TasksRepository.DisplayTask(id)

	if err != nil {
		if err == repository.ErrTaskNotFound {
			zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
			logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
			logger.Error().Stack().Err(err).Msg("task not found")
			//log.Println(err)
			c.AbortWithStatusJSON(http.StatusNotFound, handler.NewProblem(http.StatusNotFound, "Task not found"))
			return
		}
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger.Error().Stack().Err(err).Msg("internal server error")
		//log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if task.UserID != userId {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger.Error().Stack().Err(err).Msg("unauthorized")
		//log.Println("Unauthorized")
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
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger.Error().Stack().Err(err).Msg("unprocessable entity")
		//log.Println(err)
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	userId, _ := c.Get("userId")
	task, err := t.TasksRepository.Create(cRequest.Title, userId.(int64))
	if err != nil {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger.Error().Stack().Err(err).Msg("internal server error")
		//log.Println(err)
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
	task, err := t.TasksRepository.DisplayTask(id)
	if err != nil {
		if err == repository.ErrTaskNotFound {
			zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
			logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
			logger.Error().Stack().Err(err).Msg("task not found")
			//log.Println(err)
			c.AbortWithStatusJSON(http.StatusNotFound, handler.NewProblem(http.StatusNotFound, "Task not found"))
			return
		}
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger.Error().Stack().Err(err).Msg("internal server error")
		//log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userId, _ := c.Get("userId")
	if task.UserID != userId {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger.Error().Stack().Err(err).Msg("unauthorized")
		//log.Println("Unauthorized")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err = t.TasksRepository.Delete(id)
	if err != nil {
		if err == repository.ErrTaskNotFound {
			zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
			logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
			logger.Error().Stack().Err(err).Msg("task not found")
			//log.Println(err)
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
	task, err := t.TasksRepository.DisplayTask(id)
	if err != nil {
		if err == repository.ErrTaskNotFound {
			zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
			logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
			logger.Error().Stack().Err(err).Msg("task not found")
			//log.Println(err)
			c.AbortWithStatusJSON(http.StatusNotFound, handler.NewProblem(http.StatusNotFound, "Task not found"))
			return
		}
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger.Error().Stack().Err(err).Msg("internal server error")
		//log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userId, _ := c.Get("userId")
	if task.UserID != userId {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger.Error().Stack().Err(err).Msg("unauthorized")
		//log.Println("Unauthorized")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	uRequest := dto.TaskUpdateRequest{}
	if err := c.BindJSON(&uRequest); err != nil {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger.Error().Stack().Err(err).Msg("unprocessable entity")
		//log.Println(err)
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	resp := dto.Task{}
	updateResult, err := t.TasksRepository.Update(id, uRequest.Title, uRequest.Status)
	if err != nil {
		if err == repository.ErrUnauthorized {
			zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
			logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
			logger.Error().Stack().Err(err).Msg("unauthorized")
			//log.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger.Error().Stack().Err(err).Msg("internal server error")
		//log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	resp.FromEntity(updateResult)
	c.Header("Content-Type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayload(c.Writer, &resp); err != nil {
		log.Fatal().Err(err).Msg("can not respond")
	}
}
