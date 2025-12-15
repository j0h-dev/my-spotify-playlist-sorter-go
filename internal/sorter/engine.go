package sorter

import (
	"fmt"

	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/spotify"
	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/ui"
	spotifyApi "github.com/zmb3/spotify/v2"
)

func Run(client *spotifyApi.Client, playlistID spotifyApi.ID, country string, progress ui.Progress) error {
	playlist, err := spotify.GetPlaylist(
		client,
		playlistID,
		country,
		progress,
	)
	if err != nil {
		return fmt.Errorf("failed to get playlist: %w", err)
	}

	originalTracks, err := spotify.GetPlaylistTracks(client, playlist, country, progress)
	if err != nil {
		return fmt.Errorf("failed to fetch playlist: %w", err)
	}

	var finalSorted []*spotifyApi.PlaylistItem

	if progress != nil {
		progress.Start(len(originalTracks), "Sorting playlist")
	}

	artistGroups := GroupTracksByArtist(originalTracks)

	// For each artist → group albums → sort albums → sort tracks → flatten
	for _, artist := range artistGroups {
		// Group albums
		albumGroups := GroupTracksByAlbum(artist.Tracks)

		// Sort albums by release date
		SortAlbumsByReleaseDate(albumGroups)

		// Sort tracks inside each album (disc + track number)
		SortTracksInsideAlbums(albumGroups)

		// Flatten album groups into final slice
		for _, album := range albumGroups {
			finalSorted = append(finalSorted, album.Tracks...)
			if progress != nil {
				progress.Advance(len(album.Tracks))
			}
		}
	}

	if progress != nil {
		progress.Finish()
	}

	// Reorder playlist in Spotify
	err = spotify.ReorderPlaylistTracks(client, playlistID, originalTracks, finalSorted, progress)
	if err != nil {
		return fmt.Errorf("failed to reorder playlist: %w", err)
	}

	return nil
}
