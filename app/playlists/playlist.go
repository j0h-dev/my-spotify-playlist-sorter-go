package playlists

import (
	"context"
	"math"
	"strings"

	"github.com/schollz/progressbar/v3"
	"github.com/zmb3/spotify/v2"
)

func GetPlaylist(
	sp *spotify.Client,
	playlistId spotify.ID,
	country string,
) (*spotify.FullPlaylist, error) {
	ctx := context.Background()

	playlist, err := sp.GetPlaylist(ctx, playlistId, spotify.RequestOption(spotify.Country(country)))

	if err != nil {
		return nil, err
	}

	return playlist, nil
}

func GetPlaylistTracks(
	sp *spotify.Client,
	playlistId spotify.ID,
	country string,
) ([]*spotify.PlaylistItem, error) {
	ctx := context.Background()

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

func ExtractSpotifyID(url string) string {
	parts := strings.Split(url, "/")

	if len(parts) >= 5 {
		idWithQuery := parts[len(parts)-1]
		id := strings.Split(idWithQuery, "?")[0]
		return id
	}

	return url
}
