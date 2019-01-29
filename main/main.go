package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
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

func main() {
	parseFlags()

	//Config things
	spotifyAuthenticator := spotifyClient.NewSpotifyAuthenticator(*spotifyClientId, *spotifySecret, *redirectUri)
	tokenManager := auth.NewInMemoryTokenKeeper()
	authFactory := func() spotify.Authenticator {
		return spotifyAuthenticator
	}
	playlistRepository := playlists.NewInMemoryPlaylistSinkRepository(tokenManager, authFactory)
	playlistRepository.AddPlaylistForUserAndChat(98025430, 98025430, "7coO8tS7abtBgV51jPju3n")
	playlistRepository.AddPlaylistForUserAndChat(-305749618, 98025430, "7coO8tS7abtBgV51jPju3n")

	bot := telegram.NewPlaylistBot(*botToken, func(s string) string {
		return spotifyAuthenticator.AuthURL(s)
	},
		tokenManager,
		authFactory,
		playlistRepository,
	)


	router := web.SetupRouter(spotifyAuthenticator, tokenManager)

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
