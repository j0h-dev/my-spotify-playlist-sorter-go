package sorter

import (
	"time"

	spotifyApi "github.com/zmb3/spotify/v2"
)

type AlbumGroup struct {
	Name        string
	ReleaseDate time.Time
	Tracks      []*spotifyApi.PlaylistItem
}

type ArtistGroup struct {
	Name   string
	Tracks []*spotifyApi.PlaylistItem
}
