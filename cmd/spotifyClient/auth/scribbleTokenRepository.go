package auth

import (
	"fmt"
	scribble "github.com/nanobox-io/golang-scribble"
	"golang.org/x/oauth2"
	"strconv"
)

type scribbleTokenRepository struct {
	tokenStore *scribble.Driver
}

func NewScribbleTokenRepository() TokenManager{
	db, err := scribble.New("./db", nil)
	if err != nil {
		fmt.Println("Error", err)
	}
	return &scribbleTokenRepository{tokenStore:db}
}

func (this scribbleTokenRepository) GetToken(telegramUserId int) (*oauth2.Token, error) {
	token := oauth2.Token{}
	this.tokenStore.Read("tokens", strconv.Itoa(telegramUserId), &token)

	return &token, nil
}

func (this *scribbleTokenRepository) Store(token oauth2.Token, telegramUserId string) {
	this.tokenStore.Write("tokens", telegramUserId, &token)
}

func (this scribbleTokenRepository) PrintState() string {
	return "Not implemented"
}
