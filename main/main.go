package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"spotilist/cmd/telegram"
)

var botToken = flag.String("botToken", "", "Telegram API bot token")

func main(){
	parseFlags()

	bot := telegram.NewPlaylistBot(*botToken)

	bot.Start()

}

func parseFlags(){
	flag.Parse()
	logrus.Info(*botToken)
}
