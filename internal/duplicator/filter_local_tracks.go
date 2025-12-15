package duplicator

import (
	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/ui"
	"github.com/zmb3/spotify/v2"
)

func FilterLocalTracks(
	tracks []*spotify.PlaylistItem,
	progress ui.Progress,
) []*spotify.PlaylistItem {
	tracksLength := len(tracks)

	if progress != nil {
		progress.Start(tracksLength, "Filtering local tracks")
		defer progress.Finish()
	}

	filteredTracks := make([]*spotify.PlaylistItem, 0)

	for _, track := range tracks {
		if progress != nil {
			progress.Advance(1)
		}

		if track.IsLocal {
			continue
		}

		filteredTracks = append(filteredTracks, track)
	}

	return filteredTracks
}
