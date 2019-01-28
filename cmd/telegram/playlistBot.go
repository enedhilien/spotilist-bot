package telegram

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"spotilist/cmd/spotifyClient/auth"
	"time"
)

type PlaylistBot struct {
}

func NewPlaylistBot(token string,
	authUrlPrinter func(string) string,
	tokenStatePrinter auth.TokenStatePrinter) *tb.Bot {
	bot, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Error(err)
		panic(err)
	}

	bot.Handle("/authUrl", func(m *tb.Message) {
		if !m.Private() {
			return
		}
		bot.Send(m.Sender, authUrlPrinter(fmt.Sprintf("%v", m.Sender.ID)))
	})

	bot.Handle("/tokenStatus", func(m *tb.Message){
		if !m.Private() { // TODO admin rights authentication
			return
		}
		log.Debug(tokenStatePrinter.PrintState())
		bot.Send(m.Sender, tokenStatePrinter.PrintState())
	})

	bot.Handle(tb.OnText, func(m *tb.Message) {
		log.Info(m.Text)
	})

	return bot
}
