package spotify

import (
	"context"
	"fmt"

	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/ui"
	"github.com/zmb3/spotify/v2"
)

func RemoveDuplicateTracksFromPlaylist(
	client *spotify.Client,
	playlistID spotify.ID,
	tracks []*spotify.PlaylistItem,
	// Optional progress bar
	progress ui.Progress,
) ([]*spotify.PlaylistItem, error) {

	if progress != nil {
		progress.Start(len(tracks), "Identifying duplicate tracks")
	}

	// Track uniqueness by URI
	seen := make(map[spotify.URI]struct{}, len(tracks))
	var duplicates []spotify.TrackToRemove
	unique := make([]*spotify.PlaylistItem, 0, len(tracks))

	for index, track := range tracks {

		if track.Track.Track == nil {
			if progress != nil {
				progress.Advance(1)
			}
			continue
		}

		uri := track.Track.Track.URI

		if _, exists := seen[uri]; exists {
			duplicates = append(
				duplicates,
				spotify.TrackToRemove{
					URI:       string(uri),
					Positions: []int{index},
				},
			)
		} else {
			seen[uri] = struct{}{}
			unique = append(unique, track)
		}

		if progress != nil {
			progress.Advance(1)
		}
	}

	if progress != nil {
		progress.Finish()
	}

	// No duplicates detected
	if len(duplicates) == 0 {
		return unique, nil
	}

	// Remove duplicate tracks from playlist
	_, err := client.RemoveTracksFromPlaylistOpt(
		context.Background(),
		playlistID,
		duplicates,
		"",
	)

	if err != nil {
		return nil, fmt.Errorf("failed to remove duplicate tracks: %w", err)
	}

	return unique, nil
}
