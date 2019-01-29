package telegram

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
	tb "gopkg.in/tucnak/telebot.v2"
	"spotilist/cmd/spotifyClient/auth"
	"spotilist/cmd/spotifyClient/playlists"
	"strings"
	"time"
)

type PlaylistBot struct {
}

func NewPlaylistBot(token string,
	authUrlPrinter func(string) string,
	tokenManager auth.TokenManager,
	authFactory func() spotify.Authenticator,
	playlistRepository playlists.PlaylistSinkRepository,
) *tb.Bot {
	bot, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Error(err)
		panic(err)
	}

	bot.Handle("/auth", func(m *tb.Message) {
		bot.Send(m.Sender, fmt.Sprintf("Use this url to authenticate with Spotify\n %v", authUrlPrinter(stringifyUserId(m))));
	})

	bot.Handle("/tokenStatus", func(m *tb.Message) {
		if !m.Private() { // TODO admin rights authentication
			return
		}
		log.Debug(tokenManager.PrintState())
		bot.Send(m.Sender, tokenManager.PrintState())
	})

	bot.Handle("/playlists", func(m *tb.Message) {
		if !m.Private() { // TODO admin rights authentication
			return
		}
		spotifyClient, err := createSpotifyClient(m, tokenManager, authFactory)
		if err != nil {
			log.Error(err)
			bot.Send(m.Sender, fmt.Sprintf("Couldn't create Spotify client. Error: %v", err))
		} else {
			playlists, err := spotifyClient.CurrentUsersPlaylists()
			if err != nil {
				log.Error(err)
				bot.Send(m.Sender, fmt.Sprintf("Couldn't retrieve your playlists. Error: %v", err))
			}
			bot.Send(m.Sender, printPlaylists(playlists))
		}
	})

	bot.Handle("/chatid", func(m *tb.Message){
		bot.Send(m.Sender, fmt.Sprintf("Chatid: %v, Userid: %v", m.Chat.ID, m.Sender.ID))
	})

	bot.Handle(tb.OnText, func(m *tb.Message) {
		log.Info(m.Chat.ID, m.Text)
		if m.Text == "Say hi!"{
			bot.Send(m.Chat, "Hello! ^^")
			return
		}
		if match, trackId := ParseTrack(m.Text); match {
			for _, playlist := range playlistRepository.GetPlaylistsForChat(m.Chat.ID) {
				log.Info(fmt.Sprintf("User %v(%v) wants to add track %v from chat %v(%v) to playlists %v", m.Sender.Username, m.Sender.ID, trackId, m.Chat.Username, m.Chat.ID, playlist.ID()))
				playlist.AddTrackOnPosition(trackId, 0)
			}
		}
	})

	return bot
}

func createSpotifyClient(m *tb.Message,
	tokenManager auth.TokenManager,
	authFactory func() spotify.Authenticator) (*spotify.Client, error) {
	token, err := tokenManager.GetToken(m.Sender.ID)
	if err != nil {
		return nil, err
	} else {
		auth := authFactory()
		client := auth.NewClient(token)
		return &client, nil
	}
}

func printPlaylists(page *spotify.SimplePlaylistPage) string {
	sb := strings.Builder{}
	sb.WriteString("Your playlists:\n")
	for _, playlist := range page.Playlists {
		sb.WriteString(fmt.Sprintf("%v %v %v\n", playlist.Name, playlist.ID, playlist.IsPublic))
	}
	return sb.String()
}

func stringifyUserId(m *tb.Message) string {
	return fmt.Sprintf("%v", m.Sender.ID)
}
