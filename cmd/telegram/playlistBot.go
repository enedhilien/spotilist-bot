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

	bot.Handle("/spotify_status", func(m *tb.Message) {
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

	bot.Handle("/chat_id", func(m *tb.Message) {
		bot.Send(m.Sender, fmt.Sprintf("Chatid: %v, Userid: %v", m.Chat.ID, m.Sender.ID))
	})

	bot.Handle("/subscribe", func(m *tb.Message) {
		payload := m.Payload
		log.Info(payload)
		if isSpotifyConnectionValid(m, tokenManager, authFactory) {
			client, _ := createSpotifyClient(m, tokenManager, authFactory)
			playlist, playlistErr := client.GetPlaylist(spotify.ID(payload))
			user, _ := client.CurrentUser() // TODO
			if playlistErr != nil {
				bot.Send(m.Chat, fmt.Sprintf("Can't subscribe to playlist %v, error: %v", payload, playlistErr))
			} else {
				log.Info("Playlist owner", playlist.Owner.ID)
				log.Info("Spotify token owner", user.ID)
				if playlist.Owner.ID != user.ID{
					bot.Send(m.Chat, "Sorry, this isn't your playlist. Create your own one so I can mess it with (;")
				}else{
					playlistRepository.AddPlaylistForUserAndChat(m.Chat.ID, m.Sender.ID, payload)
					bot.Send(m.Chat, "Subscription is active! I will be adding tracks to this playlist as they appear here.")
				}
			}
		} else {
			bot.Send(m.Sender, fmt.Sprintf("It seems that you tried to subscribe playlist %v from chat %v - but I can't access Spotify."+
				" Please authenticate using following link and try again. \n%v", m.Payload, m.Chat.Title, authUrlPrinter(stringifyUserId(m))))
		}
	})

	bot.Handle(tb.OnAddedToGroup, func(m *tb.Message) {
		if isSpotifyConnectionValid(m, tokenManager, authFactory) {
			bot.Send(m.Sender, fmt.Sprintf("Hey! You just added me to group %v. It seems that I'm already authenticated with your "+
				"Spotify account. Please use /subscribe playlistId command on a chat that you want me to check for new songs."))
		} else {
			bot.Send(m.Sender, fmt.Sprintf("Hey! You just added me to group %v. It seems that you are not authenticated me with Spotify. "+
				"Please follow this link and complete the process before issuing me any commands (;\n%v",
				m.Chat.Username, authUrlPrinter(stringifyUserId(m))))
		}

	})

	bot.Handle(tb.OnText, func(m *tb.Message) {
		log.Info(m.Chat.ID, m.Text)
		if m.Text == "Say hi!" {
			bot.Send(m.Chat, "Hello! ^^")
			return
		}
		if match, trackId := ParseTrack(m.Text); match {
			for _, playlist := range playlistRepository.GetPlaylistsForChat(m.Chat.ID) {
				log.WithFields(log.Fields{"operation": "track_add"}).Info(fmt.Sprintf("%v;%v;%v;%v;%v;%v", m.Sender.Username, m.Sender.ID, trackId, m.Chat.Username, m.Chat.ID, playlist.ID()))
				err := playlist.AddTrackOnPosition(trackId, 0)
				if err != nil{
					log.Error(err)
				}
			}
		}
	})

	return bot
}

func isSpotifyConnectionValid(m *tb.Message,
	tokenManager auth.TokenManager,
	authFactory func() spotify.Authenticator) bool {
	spotifyClient, err := createSpotifyClient(m, tokenManager, authFactory)
	if err == nil {
		_, err = spotifyClient.CurrentUsersPlaylists()
		if err != nil {
			return false
		} else {
			return true
		}
	} else {
		return false
	}

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
