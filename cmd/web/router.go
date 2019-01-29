package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
	"net/http"
	"spotilist/cmd/spotifyClient"
	"spotilist/cmd/spotifyClient/auth"
	"spotilist/cmd/spotifyClient/playlists"
	"strconv"
)

func SetupRouter(authenticator spotify.Authenticator, keeper auth.TokenKeeper, repository playlists.PlaylistSinkRepository) *gin.Engine {
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

	r.Handle(http.MethodPost, "/playlist", func(c *gin.Context) {
		userIdQP := c.Query("userId")
		chatIdQP := c.Query("chatId")
		playlistIdQP := c.Query("playlistId")
		if userIdQP == "" || chatIdQP == "" || playlistIdQP == ""{
			c.Error(errors.New("Wrong request"))
			return
		}

		chatId, _ := strconv.ParseInt(chatIdQP, 10, 64)
		userId, _ := strconv.ParseInt(userIdQP, 10, 32)
		logrus.Info(fmt.Sprintf("Adding playlist %v %v %v", chatId, userId, playlistIdQP))
		repository.AddPlaylistForUserAndChat(chatId, int(userId), playlistIdQP)

	})
	r.Handle(http.MethodGet, "/version", func(c *gin.Context) {
		c.String(200, "Hey Space Cowboy")
	})

	return r
}
