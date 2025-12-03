package duplicator

import (
	"context"
	"fmt"

	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/spotify"
	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/ui"
	spotifyApi "github.com/zmb3/spotify/v2"
)

func Run(client *spotifyApi.Client, playlistID spotifyApi.ID, country string, info PlaylistInfo) error {
	progress := ui.NewProgressBar()

	playlist, err := spotify.GetPlaylist(client, playlistID, country, progress)

	if err != nil {
		return fmt.Errorf("failed to get playlist: %w", err)
	}

	tracks, err := spotify.GetPlaylistTracks(client, playlist, country, progress)

	if err != nil {
		return fmt.Errorf("failed to get playlist tracks: %w", err)
	}

	user, err := client.CurrentUser(context.Background())

	createdPlaylist, err := spotify.CreatePlaylist(
		client,
		user.ID,
		info.Name,
		info.Description,
		info.Public,
		info.Collaborative,
	)

	if err != nil {
		return fmt.Errorf("failed to create new playlist for user: %w", err)
	}

	tracks = FilterLocalTracks(tracks, progress)

	err = spotify.AddTracksToPlaylist(client, tracks, createdPlaylist, progress)

	if err != nil {
		return fmt.Errorf("failed to add tracks to new playlist: %w", err)
	}

	return nil
}
