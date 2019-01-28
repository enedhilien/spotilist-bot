package telegram

import (
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

type PlaylistBot struct {
}

func NewPlaylistBot(token string) *tb.Bot {
	bot, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Error(err)
		panic(err)
	}

	bot.Handle(tb.OnText, func(m *tb.Message) {
		log.Info(m.Text)
	})

	return bot
}
