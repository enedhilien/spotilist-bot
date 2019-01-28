package telegram

import "testing"
import "github.com/stretchr/testify/assert"

func TestMatchTrack(t *testing.T) {
	ok, track := ParseTrack("https://open.spotify.com/track/6EztTMv6l0WEj7H8tIjn9i")
	assert.True(t, ok)
	assert.EqualValues(t, "6EztTMv6l0WEj7H8tIjn9i", track)
}

func TestMatchTrackWithSi(t *testing.T) {
	ok, track := ParseTrack("https://open.spotify.com/track/5p5L0CuLxfFUGp3dhZ6nZr?si=AUu4H77qQjeuZsMdveW-pQ")
	assert.True(t, ok)
	assert.EqualValues(t, "5p5L0CuLxfFUGp3dhZ6nZr", track)
}
