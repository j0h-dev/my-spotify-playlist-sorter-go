package playlistutils

import (
	"context"
	"math"
	"os"

	"github.com/schollz/progressbar/v3"
	"github.com/zmb3/spotify/v2"
)

func GetPlaylist(sp *spotify.Client, playlistId spotify.ID) (*spotify.FullPlaylist, error) {
	ctx := context.Background()

	country := os.Getenv("COUNTRY")
	playlist, err := sp.GetPlaylist(ctx, playlistId, spotify.RequestOption(spotify.Country(country)))

	if err != nil {
		return nil, err
	}

	return playlist, nil
}

func GetPlaylistTracks(sp *spotify.Client, playlistId spotify.ID) ([]*spotify.PlaylistItem, error) {
	ctx := context.Background()

	country := os.Getenv("COUNTRY")
	playlistDetails, err := sp.GetPlaylist(ctx, playlistId, spotify.RequestOption(spotify.Country(country)))

	if err != nil {
		return nil, err
	}

	bar := progressbar.Default(int64(playlistDetails.Tracks.Total), "Fetching playlist tracks")

	tracks := make([]*spotify.PlaylistItem, playlistDetails.Tracks.Total)

	total := float64(playlistDetails.Tracks.Total)
	pages := math.Ceil(total / 100.0)

	for i := 0; i < int(pages); i++ {
		trackPage, err := sp.GetPlaylistItems(
			ctx,
			playlistId,
			spotify.RequestOption(spotify.Offset(i*100)),
			spotify.RequestOption(spotify.Limit(100)),
			spotify.RequestOption(spotify.Country(country)),
		)

		if err != nil {
			return nil, err
		}

		for itemIndex, item := range trackPage.Items {
			bar.Add(1)
			tracks[i*100+itemIndex] = &item
		}
	}

	return tracks, nil
}
