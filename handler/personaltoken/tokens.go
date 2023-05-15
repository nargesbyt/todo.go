package personaltoken

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
	"github.com/nargesbyt/todo.go/handler"
	"github.com/nargesbyt/todo.go/internal/dto"
	"github.com/nargesbyt/todo.go/repository"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"time"
)

type PersonalTokens struct {
	TokenRepository repository.PersonalAccessToken
}

func (t PersonalTokens) Create(c *gin.Context) {
	createTokenRequest := dto.CreateTokenRequest{}
	err := c.BindJSON(&createTokenRequest)
	if err != nil {
		log.Error().Stack().Err(err).Msg("unprocessable entity")
		c.AbortWithStatus(http.StatusUnprocessableEntity)

		return
	}
	userId, _ := c.Get("userId")
	expireTime, err := time.Parse(time.RFC3339, createTokenRequest.ExpiresAt)
	//fmt.Println(expireTime)
	token, err := t.TokenRepository.Add(createTokenRequest.Title, expireTime, userId.(int64))
	if err != nil {
		log.Error().Stack().Err(err).Msg("internal server error while inserting a record in tokens table")
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}
	resp := dto.Tokens{}
	resp.FromEntity(token)
	c.Header("Content-Type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayload(c.Writer, &resp); err != nil {
		log.Fatal().Err(err).Msg("can not respond")
	}
	//c.JSON(http.StatusCreated, token)
}
func (t PersonalTokens) GetToken(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Error().Stack().Err(err).Msg("can not assign id from url param")

		return
	}
	userId, _ := c.Get("userId")
	token, err := t.TokenRepository.Get(id)

	if err != nil {
		if err == repository.ErrTokenNotFound {
			log.Error().Stack().Err(err).Msg("token not found")
			c.AbortWithStatusJSON(http.StatusNotFound, handler.NewProblem(http.StatusNotFound, "Token not found"))

			return
		}

		log.Error().Stack().Err(err).Msg("internal server error")
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}
	if token.UserID != userId {
		log.Error().Stack().Err(err).Msg("unauthorized")
		c.AbortWithStatus(http.StatusUnauthorized)

		return
	}
	/*resp := dto.Tokens{}
	resp.FromEntity(token)
	c.Header("Content-Type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayload(c.Writer, &resp); err != nil {
		log.Fatal().Err(err).Msg("can not respond")
	}*/
	c.JSON(http.StatusOK, &token)
}
func (t PersonalTokens) List(c *gin.Context) {
	userId, _ := c.Get("userId")
	tokens, err := t.TokenRepository.List(c.Query("title"), userId.(int64))
	if err != nil {
		if err == repository.ErrTokenNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			log.Error().Stack().Err(err).Msg("token not found")

			return
		}
		c.AbortWithStatus(http.StatusBadRequest)
		log.Error().Stack().Err(err).Msg("can not fetch record from database")

		return
	}
	var dtoTokens []*dto.Tokens
	for _, token := range tokens {
		resp := dto.Tokens{}
		resp.FromEntity(*token)
		dtoTokens = append(dtoTokens, &resp)

	}
	c.Header("Content-Type", jsonapi.MediaType)
	if err := jsonapi.MarshalPayload(c.Writer, dtoTokens); err != nil {
		log.Fatal().Err(err).Msg("can not respond")
	}

}
func (t PersonalTokens) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	token, err := t.TokenRepository.Get(id)
	if err != nil {
		if err == repository.ErrTokenNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			log.Error().Stack().Err(err).Msg("token not found")

			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Error().Stack().Err(err).Msg("internal server error")

		return
	}
	userId, _ := c.Get("userId")
	if token.UserID != userId {
		log.Error().Stack().Err(errors.New("unauthorized ")).Msg("unauthorized")
		c.AbortWithStatus(http.StatusUnauthorized)

		return
	}

	uRequest := dto.UpdateRequest{}
	if err := c.BindJSON(&uRequest); err != nil {
		log.Error().Stack().Err(err).Msg("unprocessable entity")
		c.AbortWithStatus(http.StatusUnprocessableEntity)

		return
	}

	resp := dto.Tokens{}
	updateResult, err := t.TokenRepository.Update(id, uRequest.Title, uRequest.ExpiresAt)
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
		log.Fatal().Err(err).Msg("unsuccessful token update")
	}
}

func (t PersonalTokens) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Error().Stack().Err(err).Msg("can not assign id from url param")

		return
	}
	token, err := t.TokenRepository.Get(id)
	if err != nil {
		if err == repository.ErrTokenNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			log.Error().Stack().Err(err).Msg("token not found")

			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Error().Stack().Err(err).Msg("Internal server error")

		return
	}
	userId, _ := c.Get("userId")
	if userId != token.UserID {
		c.AbortWithStatus(http.StatusUnauthorized)
		log.Error().Stack().Err(errors.New("unauthorized")).Msg("unauthorized")

		return
	}
	err = t.TokenRepository.Delete(id)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Error().Stack().Err(err).Msg("internal server error")

		return
	}
	c.Status(http.StatusAccepted)

}
