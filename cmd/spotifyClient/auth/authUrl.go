package auth

import "github.com/zmb3/spotify"

type AuthUrlPrinter interface {
	AuthUrl() string
}

type WebUrlPrinter struct {
	authenticator spotify.Authenticator
}

func (this WebUrlPrinter) AuthUrl() string{
	return this.authenticator.AuthURL("WebState")
}

type CallbackAuthUrlPrinter struct {
	authenticator spotify.Authenticator
	state func() string
}

func (this CallbackAuthUrlPrinter) AuthUrl() string{
	return this.authenticator.AuthURL(this.state())
}




