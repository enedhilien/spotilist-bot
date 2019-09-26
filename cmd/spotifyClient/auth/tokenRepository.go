package auth

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"strings"
)

type TokenManager interface {
	OauthTokenRepository
	TokenStatePrinter
}

type OauthTokenRepository interface {
	GetToken(telegramUserId int) (*oauth2.Token, error)
	Store(token oauth2.Token, telegramUserId string)
}

type TokenStatePrinter interface {
	PrintState() string
}

func NewInMemoryTokenKeeper() TokenManager{
	return &inMemoryTokenKeeper{tokenStore: map[string]*oauth2.Token{}}
}

type inMemoryTokenKeeper struct {
	tokenStore map[string]*oauth2.Token
}

func (this inMemoryTokenKeeper) GetToken(telegramUserId int) (*oauth2.Token, error) {
	if val, ok := this.tokenStore[fmt.Sprintf("%v", telegramUserId)]; ok {
		return val, nil
	} else {
		return nil, errors.New(fmt.Sprintf("No token for user %v", telegramUserId))
	}
}

func (this *inMemoryTokenKeeper) Store(token oauth2.Token, telegramUserId string) {
	this.tokenStore[telegramUserId] = &token
}

func (this inMemoryTokenKeeper) PrintState() string {
	sb := strings.Builder{}
	for userId, token := range this.tokenStore {
		sb.WriteString(fmt.Sprintf("UserId: %v; Expires in %v; Valid: %v", userId, token.Expiry, token.Valid()))
	}
	return sb.String()
}
