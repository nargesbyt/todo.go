package task

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nargesbyt/todo.go/entity"
	"github.com/nargesbyt/todo.go/handler"
	"github.com/nargesbyt/todo.go/internal/dto"
	"github.com/nargesbyt/todo.go/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDisplayTasks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Success", func(t *testing.T) {
		var id int64 = 5
		response := dto.Task{}

		mockTaskResp := entity.Task{
			ID:        id,
			Title:     "New task",
			Status:    "pending",
			CreatedAt: time.Now(),
		}
		response.FromEntity(mockTaskResp)
		mockTaskRepository := new(repository.MockTaskRepository)
		mockTaskRepository.On("DisplayTask", id).Return(mockTaskResp, nil)
		taskrepository := Task{mockTaskRepository}
		// a response recorder for getting written http response
		rr := httptest.NewRecorder()
		//creating a fake server and context
		c, router := gin.CreateTestContext(rr)
		router.GET("/tasks/:id", taskrepository.DisplayTasks)
		var err error
		c.Request, err = http.NewRequest(http.MethodGet, "/tasks/5", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, c.Request)

		respBody, err := json.Marshal(response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockTaskRepository.AssertExpectations(t) // assert that UserService.Get was called
	})

	t.Run("BadRequest", func(t *testing.T) {
		var id string = "abc"
		mockTaskRepository := new(repository.MockTaskRepository)
		mockTaskRepository.On("DisplayTask", id).Return(entity.Task{}, errors.New("invalid task id"))
		taskrepository := Task{mockTaskRepository}
		// a response recorder for getting written http response
		rr := httptest.NewRecorder()
		c, router := gin.CreateTestContext(rr)
		router.GET("/tasks/:id", taskrepository.DisplayTasks)
		var err error
		c.Request, err = http.NewRequest(http.MethodGet, "/tasks/\"abc\"", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, c.Request)

		respBody, err := json.Marshal(handler.NewProblem(http.StatusBadRequest, "invalid task id"))

		assert.NoError(t, err)
		assert.Equal(t, respBody, rr.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockTaskRepository.AssertNotCalled(t, "DisplayTask", mock.Anything)
	})

	t.Run("NotFound", func(t *testing.T) {
		var id int64 = 5
		mockTaskRepository := new(repository.MockTaskRepository)
		mockTaskRepository.On("DisplayTask", id).Return(entity.Task{}, repository.ErrTaskNotFound)
		rr := httptest.NewRecorder()

		c, router := gin.CreateTestContext(rr)

		taskrepository := Task{mockTaskRepository}

		router.GET("/tasks/:id", taskrepository.DisplayTasks)

		var err error
		c.Request, err = http.NewRequest(http.MethodGet, "/tasks/5", nil)

		assert.NoError(t, err)

		router.ServeHTTP(rr, c.Request)

		respBody, err := json.Marshal(handler.NewProblem(http.StatusNotFound, "Task not found"))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockTaskRepository.AssertExpectations(t)
	})
	t.Run("InternalServerError", func(t *testing.T) {

		mockTaskRepository := new(repository.MockTaskRepository)
		mockTaskRepository.On("DisplayTask", mock.Anything).Return(entity.Task{}, errors.New("db connection error"))
		taskrepository := Task{mockTaskRepository}

		rr := httptest.NewRecorder()

		//c, router := gin.CreateTestContext(rr)
		router := gin.Default()
		router.GET("/tasks/:id", taskrepository.DisplayTasks)
		var err error

		Request, err := http.NewRequest(http.MethodGet, "/tasks/7", nil)

		assert.NoError(t, err)

		router.ServeHTTP(rr, Request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockTaskRepository.AssertNotCalled(t, "DisplayTask", 5)

	})
}

func TestAddTask(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Success", func(t *testing.T) {
		//request := TaskCreateRequest{}
		var title string = "task jadid"
		response := dto.Task{}
		mockTaskResponse := entity.Task{
			ID:        2,
			Title:     title,
			Status:    "pending",
			CreatedAt: time.Now(),
			//FinishedAt: nil,

		}
		response.FromEntity(mockTaskResponse)
		mockTaskRepository := new(repository.MockTaskRepository)
		mockTaskRepository.On("Create", title).Return(mockTaskResponse, nil)
		taskRepository := Task{mockTaskRepository}
		rr := httptest.NewRecorder()
		c, router := gin.CreateTestContext(rr)
		router.POST("/tasks", taskRepository.AddTask)

		taskreq := entity.Task{Title: title}
		jsonValue, _ := json.Marshal(taskreq)
		var err error
		c.Request, err = http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonValue))

		assert.NoError(t, err)

		router.ServeHTTP(rr, c.Request)
		respBody, err := json.Marshal(response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockTaskRepository.AssertExpectations(t)

	})
	t.Run("ServerError", func(t *testing.T) {
		mockTaskRepository := new(repository.MockTaskRepository)

		mockTaskRepository.On("Create", mock.Anything).Return(entity.Task{}, errors.New("Internal Server Error"))

		taskRepository := Task{mockTaskRepository}

		rr := httptest.NewRecorder()
		c, router := gin.CreateTestContext(rr)

		var title string = "New task"
		taskreq := entity.Task{Title: title}

		jsonValue, _ := json.Marshal(taskreq)
		//var err error
		c.Request, _ = http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonValue))
		//assert.NoError(t, err)
		router.POST("/tasks", taskRepository.AddTask)
		router.ServeHTTP(rr, c.Request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		mockTaskRepository.AssertExpectations(t)

	})
}
func TestUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Success", func(t *testing.T) {
		var id int64 = 5
		var title string = "new task"
		var status string = "in progress"
		updateResponse := dto.Task{}
		mockTaskResponse := entity.Task{
			ID:        id,
			Title:     title,
			Status:    status,
			CreatedAt: time.Now(),
		}
		updateRequest := dto.TaskUpdateRequest{
			Title:  title,
			Status: status,
		}

		updateResponse.FromEntity(mockTaskResponse)

		mockTaskRepository := new(repository.MockTaskRepository)
		mockTaskRepository.On("Update", id, title, status).Return(mockTaskResponse, nil)
		taskRepository := Task{mockTaskRepository}

		rr := httptest.NewRecorder()
		c, router := gin.CreateTestContext(rr)

		router.PATCH("/tasks/:id", taskRepository.Update)

		//updateRequest := mockTask
		jsonValue, _ := json.Marshal(updateRequest)
		var err error
		c.Request, err = http.NewRequest("PATCH", "/tasks/5", bytes.NewBuffer(jsonValue))

		assert.NoError(t, err)

		router.ServeHTTP(rr, c.Request)

		respBody, err := json.Marshal(updateResponse)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		//mockTaskRepository.AssertCalled(t, "Update", id, title, status)

	})
	//t.Run("Error", func(t *testing.T) {
	//
	//})

}
func TestDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Success", func(t *testing.T) {
		var id int64
		id = 7
		mockTaskRepository := new(repository.MockTaskRepository)
		mockTaskRepository.On("Delete", id).Return(nil)
		taskRepository := Task{mockTaskRepository}

		rr := httptest.NewRecorder()
		c, router := gin.CreateTestContext(rr)

		router.DELETE("/tasks/:id", taskRepository.DeleteTask)
		var err error
		c.Request, err = http.NewRequest(http.MethodDelete, "/tasks/7", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, c.Request)
		assert.Equal(t, http.StatusAccepted, rr.Code)

	})
}
func TestFind(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Success", func(t *testing.T) {
		var title1 string = "task1"
		//var title2 string = "task two"
		var status string = "pending"
		mockTaskResponse := []entity.Task{
			{
				Title:  title1,
				Status: status,
			},
		}
		mockTaskRepository := new(repository.MockTaskRepository)
		mockTaskRepository.On("Find", title1, status).Return(mockTaskResponse, nil)

		task := Task{mockTaskRepository}
		rr := httptest.NewRecorder()

		c, router := gin.CreateTestContext(rr)
		router.GET("/tasks", task.List)
		var err error
		c.Request, err = http.NewRequest(http.MethodGet, "/tasks?title=task1&status=pending", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, c.Request)
		respBody, err := json.Marshal(mockTaskResponse)
		assert.NoError(t, err)

		assert.Equal(t, respBody, rr.Body.Bytes())
		assert.Equal(t, http.StatusOK, rr.Code)
		mockTaskRepository.AssertExpectations(t)

	})
	t.Run("Error", func(t *testing.T) {
		//var title string = "taskOne"
		//var status string = "pending"

		mockTaskRepository := new(repository.MockTaskRepository)
		mockTaskRepository.On("Find", mock.Anything, mock.Anything).Return([]entity.Task{}, errors.New("db connection error"))
		task := Task{mockTaskRepository}
		rr := httptest.NewRecorder()

		c, router := gin.CreateTestContext(rr)
		router.GET("/tasks", task.List)
		var err error
		c.Request, err = http.NewRequest(http.MethodGet, "/tasks?title=task1&status=pending", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, c.Request)
		respBody, err := json.Marshal(handler.NewProblem(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
		//respBody, err := json.Marshal(http.StatusInternalServerError)
		assert.NoError(t, err)

		assert.Equal(t, respBody, rr.Body.Bytes())
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockTaskRepository.AssertExpectations(t)

	})
}
