package playlists

import (
	"fmt"
	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
	"spotilist/cmd/spotifyClient/auth"
)

type UserPlaylist interface {
	AddTrackOnPosition(trackId string, position int) error
	ID() spotify.ID
}

type SpotifyUserPlaylist struct {
	spotifyClient spotify.Client
	playlistId    spotify.ID
}

func (this SpotifyUserPlaylist) AddTrackOnPosition(trackId string, position int) error {
	_, err := this.spotifyClient.AddTracksToPlaylistOnPosition(this.playlistId, position, spotify.ID(trackId))
	return err
}

func (this SpotifyUserPlaylist) ID() spotify.ID {
	return this.playlistId
}

type PlaylistSinkRepository interface {
	GetPlaylistsForChat(chatId int64) []UserPlaylist
	AddPlaylistForUserAndChat(chatId int64, userId int, playlistId string)
	RemovePlaylistForChat(chatId int64, userId int)
}

func NewInMemoryPlaylistSinkRepository(manager auth.TokenManager, authFactory func() spotify.Authenticator) PlaylistSinkRepository {
	db, err := scribble.New("./db", nil)
	if err != nil {
		logrus.Error(err)
	}
	chatIdToUserId := map[int64][]UserPlaylistEntry{}
	db.Read("playlists", "map", &chatIdToUserId)
	return &inMemoryPlaylistSinkRepository{tokenManager: manager, chatIdToUserId: chatIdToUserId, authFactory: authFactory, db:db}
}

type inMemoryPlaylistSinkRepository struct {
	db 	*scribble.Driver
	tokenManager   auth.TokenManager
	chatIdToUserId map[int64][]UserPlaylistEntry
	authFactory    func() spotify.Authenticator
}

type UserPlaylistEntry struct {
	PlaylistId string
	UserId     int
}

func (this inMemoryPlaylistSinkRepository) GetPlaylistsForChat(chatId int64) []UserPlaylist {
	results := []UserPlaylist{}
	clientsCache := map[int]spotify.Client{}

	if userIds, exists := this.chatIdToUserId[chatId]; exists {
		for _, playlistEntry := range userIds {
			if _, isCached := clientsCache[playlistEntry.UserId]; !isCached { // First we update the cache
				if token, err := this.tokenManager.GetToken(playlistEntry.UserId); err != nil {
					logrus.Error(fmt.Sprintf("No token exists for user %v, altough he has registered playlists with chat %v", playlistEntry, chatId))
					panic(err) //TODO
				} else {
					auth := this.authFactory()
					clientsCache[playlistEntry.UserId] = auth.NewClient(token)
				}
			}
			// and then create <playlist, track> handler
			results = append(results, SpotifyUserPlaylist{spotifyClient: clientsCache[playlistEntry.UserId], playlistId: spotify.ID(playlistEntry.PlaylistId)})
		}
	}
	return results
}

func (this *inMemoryPlaylistSinkRepository) AddPlaylistForUserAndChat(chatId int64, userId int, playlistId string) {
	if _, ok := this.chatIdToUserId[chatId]; !ok {
		this.chatIdToUserId[chatId] = []UserPlaylistEntry{}
	}
	this.chatIdToUserId[chatId] = append(this.chatIdToUserId[chatId], UserPlaylistEntry{PlaylistId: playlistId, UserId: userId})
	this.db.Write("playlists", "map", this.chatIdToUserId)
}

func (this *inMemoryPlaylistSinkRepository) RemovePlaylistForChat(chatId int64, userId int) {
	if chatEntries, ok := this.chatIdToUserId[chatId]; ok {
		for i, chatEntry := range chatEntries {
			if chatEntry.UserId == userId {
				chatEntries[i] = chatEntries[len(chatEntries)-1]
				chatEntries = chatEntries[:len(chatEntries)-1]
			}
		}
		this.chatIdToUserId[chatId] = chatEntries
		this.db.Write("playlists", "map", this.chatIdToUserId)
	}
}
