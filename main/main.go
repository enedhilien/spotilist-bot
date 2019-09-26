package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
	"os"
	"spotilist/cmd/spotifyClient"
	"spotilist/cmd/spotifyClient/auth"
	"spotilist/cmd/spotifyClient/playlists"
	"spotilist/cmd/telegram"
	"spotilist/cmd/web"
)

var botToken = flag.String("botToken", "", "Telegram API bot token")

var spotifyClientId = flag.String("spotifyClientId", "", "Spotify API client ID")
var spotifySecret = flag.String("spotifySecret", "", "spotifySecret API secret")
var redirectUri = flag.String("redirectUri", "", "spotifySecret redirect URI")

var contextPath = flag.String("contextPath", "", "server context path")
var routerSecret = flag.String("routerSecret", "", "")

func main() {
	//setup logging
	logrus.SetFormatter(&logrus.JSONFormatter{})
	file, err := os.OpenFile("spotilist.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
	 logrus.SetOutput(file)
	} else {
	 logrus.Info("Failed to log to file, using default stderr")
	}
	parseFlags()

	//Config things
	spotifyAuthenticator := spotifyClient.NewSpotifyAuthenticator(*spotifyClientId, *spotifySecret, *redirectUri)
	//tokenManager := auth.NewInMemoryTokenKeeper()
	tokenManager := auth.NewScribbleTokenRepository()
	authFactory := func() spotify.Authenticator {
		return spotifyAuthenticator
	}
	playlistRepository := playlists.NewInMemoryPlaylistSinkRepository(tokenManager, authFactory)

	bot := telegram.NewPlaylistBot(*botToken, func(s string) string {
		return spotifyAuthenticator.AuthURL(s)
	},
		tokenManager,
		authFactory,
		playlistRepository,
	)

	router := web.SetupRouter(spotifyAuthenticator, tokenManager, playlistRepository, *routerSecret)

	// Run things
	go func() {
		bot.Start()
	}()

	router.Run(*contextPath)

}

func parseFlags() {
	flag.Parse()
	logrus.Info(*botToken)
}
