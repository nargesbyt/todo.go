package user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/nargesbyt/todo.go/entity"
	"github.com/nargesbyt/todo.go/handler"
	"github.com/nargesbyt/todo.go/internal/dto"
	"github.com/nargesbyt/todo.go/repository"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"os"

	"net/http"
	"strconv"
)

type User struct {
	UsersRepository repository.Users
}

func (u User) Create(c *gin.Context) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	user := entity.User{}
	err := c.BindJSON(&user)
	if err != nil {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		logger := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
		logger.Error().Stack().Err(err).Msg("can not read body request")
		//log.Error().Err(err).Msg("can not read body request")
		//log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return

	}
	createResponse := dto.User{}
	user, err = u.UsersRepository.Create(user.Email, user.Password, user.Username)
	createResponse.FromEntity(user)
	if err != nil {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		logger := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
		logger.Error().Stack().Err(err).Msg("can not save new user in repository")
		//log.Error().Err(err).Msg("new user can not save in repository")
		//log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Header("Content-Type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayload(c.Writer, &createResponse); err != nil {
		log.Fatal().Err(err).Msg("cannot write response")
		//log.Fatal(err)
	}
	//c.JSON(http.StatusCreated, user)

}
func (u User) ListUsers(c *gin.Context) {
	users, err := u.UsersRepository.GetUsers(c.Query("email"), c.Query("username"))
	if err != nil {
		if err == repository.ErrUserNotFound {
			zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
			logger := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
			logger.Error().Stack().Err(err).Msg("user not found")
			//log.Error().Err(err).Msg("user not found")
			//log.Println(err)
			c.AbortWithStatusJSON(http.StatusNotFound, handler.NewProblem(http.StatusNotFound, "User not found"))
			return
		}
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		logger := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
		logger.Error().Stack().Err(err).Msg("internal server error")
		//log.Error().Err(err).Msg("internal server error")
		//log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return

	}
	var dtoUsers []*dto.User
	for _, user := range users {
		resp := dto.User{}
		resp.FromEntity(*user)
		dtoUsers = append(dtoUsers, &resp)
	}

	c.Header("Content-Type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayload(c.Writer, dtoUsers); err != nil {
		log.Fatal().Err(err).Msg("cannot write response")
	}
	//c.JSON(http.StatusOK, user)
}
func (u User) UpdateUsers(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		logger := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
		logger.Error().Stack().Err(err).Msg("can not convert string to int")
		//log.Error().Err(err).Msg("can not convert string to integer")
		//log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	uRequest := dto.UserUpdateRequest{}
	if err := c.BindJSON(&uRequest); err != nil {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		logger := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
		logger.Error().Stack().Err(err).Msg("can not read request body")
		//log.Error().Err(err).Msg("can not read request body")
		//log.Println(err)
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	resp := dto.User{}
	updateResult, err := u.UsersRepository.UpdateUsers(id, uRequest.Username, uRequest.Email, uRequest.Password)
	if err != nil {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		logger := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
		logger.Error().Stack().Err(err).Msg("can not save changes to repository")
		//log.Error().Err(err).Msg("can not save changes in repositort")
		//log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	resp.FromEntity(updateResult)
	c.Header("Content-Type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayload(c.Writer, &resp); err != nil {
		log.Fatal().Err(err).Msg("can not write response")
	}
	//c.JSON(http.StatusOK, resp)
}
func (u User) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = u.UsersRepository.DeleteUsers(id)
	if err != nil {
		if err == repository.ErrTaskNotFound {
			zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
			logger := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
			logger.Error().Stack().Err(err).Msg("user not found")
			//log.Error().Err(err).Msg("user not found")
			//log.Println(err)
			c.AbortWithStatusJSON(http.StatusNotFound, handler.NewProblem(http.StatusNotFound, "User not found"))
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusAccepted)

}
func (u User) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, handler.NewProblem(http.StatusBadRequest, "invalid user id"))
		return
	}
	user, err := u.UsersRepository.GetUserByID(id)

	if err != nil {
		if err == repository.ErrUserNotFound {
			zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
			logger := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
			logger.Error().Stack().Err(err).Msg("user not found")
			//log.Println(err)
			c.AbortWithStatusJSON(http.StatusNotFound, handler.NewProblem(http.StatusNotFound, "User not found"))
			return
		}
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		logger := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
		logger.Error().Stack().Err(err).Msg("internal server error")
		//log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	resp := dto.User{}
	resp.FromEntity(user)
	c.Header("Content-Type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayload(c.Writer, &resp); err != nil {
		log.Fatal().Err(err).Msg("can not response")
	}

}
