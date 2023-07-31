package oauth

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"net/http"
)

type OAuth struct {
	OAuth2Config oauth2.Config
}

func (a *OAuth) Get(c *gin.Context) {
	c.Redirect(http.StatusFound, a.OAuth2Config.AuthCodeURL("test"))
}

func (a *OAuth) Callback(c *gin.Context) {

}
