package spotifyClient

import "github.com/zmb3/spotify"

const State = "STATE"

func NewSpotifyAuthenticator(clientId, secret, redirectUri string) spotify.Authenticator{
	auth := spotify.NewAuthenticator(redirectUri, spotify.ScopePlaylistModifyPublic, spotify.ScopeUserLibraryRead)
	auth.SetAuthInfo(clientId, secret)

	return auth
}
