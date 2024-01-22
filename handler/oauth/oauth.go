package oauth

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/nargesbyt/todo.go/internal/random"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

const stateKeyPrefix = "oidc:state:"

type OAuth struct {
	OAuth2Config oauth2.Config
	RedisClient  *redis.Client
}

func (a *OAuth) setState(state string) error {

	err := a.RedisClient.Set(context.Background(), stateKeyPrefix+state, true, time.Minute*5).Err()
	if err != nil {
		return err
	}
	return nil
}
func (a *OAuth) isValidState(state string) (bool, error) {

	exist, err := a.RedisClient.Exists(context.Background(), stateKeyPrefix+state).Result()
	if err != nil {
		return false, err
	}
	if exist == 1 {
		return true, nil
	}
	return false, nil
}

func (a *OAuth) Get(c *gin.Context) {
	state := random.Token(20)
	err := a.setState(state)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Redirect(http.StatusFound, a.OAuth2Config.AuthCodeURL(state))
}

func (a *OAuth) Callback(c *gin.Context) {
	state, err := a.isValidState(c.Query("state"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if !state {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	oauth2Token, err := a.OAuth2Config.Exchange(c, c.Query("code"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, oauth2Token.Extra("id_token").(string))
}
