package web

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
	"net/http"
	"spotilist/cmd/spotifyClient"
	"spotilist/cmd/spotifyClient/auth"
)

func SetupRouter(authenticator spotify.Authenticator, keeper auth.TokenKeeper) *gin.Engine {
	r := gin.Default()

	r.Handle(http.MethodGet, "/auth", func(c *gin.Context) {

		token, err := authenticator.Token(c.Query("state"), c.Request)
		if err != nil {
			c.Error(err)
			return
		}
		logrus.Info(token)
		keeper.Store(*token, c.Query("state"))
	})

	r.Handle(http.MethodGet, "/authUrl", func(c *gin.Context) {
		c.String(200, authenticator.AuthURL(spotifyClient.State))
	})

	return r
}
