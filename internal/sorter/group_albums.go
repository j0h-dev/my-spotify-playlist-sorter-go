package sorter

import (
	"time"

	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/spotify"
	spotifyApi "github.com/zmb3/spotify/v2"
)

// GroupTracksByAlbum groups the provided tracks by their album IDs.
// It returns an ordered map where the keys are album IDs and the values are AlbumGroup structs containing album details and associated tracks.
// An optional progress bar can be provided to track the grouping process.
func GroupTracksByAlbum(
	tracks []*spotifyApi.PlaylistItem,
) []*AlbumGroup {

	// Map for fast lookup
	albumMap := make(map[spotifyApi.ID]*AlbumGroup, len(tracks))

	// Slice to preserve album order of first appearance
	orderedAlbums := make([]spotifyApi.ID, 0)

	for _, track := range tracks {
		album := track.Track.Track.Album
		albumID := album.ID

		group, exists := albumMap[albumID]
		if !exists {
			// Parse the album release date –
			// Spotify sometimes returns:
			// - YYYY
			// - YYYY-MM
			// - YYYY-MM-DD
			releaseDate, err := spotify.ParseSpotifyReleaseDate(album.ReleaseDate, album.ReleaseDatePrecision)

			if err != nil {
				releaseDate = time.Time{}
			}

			group = &AlbumGroup{
				Name:        album.Name,
				ReleaseDate: releaseDate,
				Tracks:      make([]*spotifyApi.PlaylistItem, 0),
			}

			albumMap[albumID] = group
			orderedAlbums = append(orderedAlbums, albumID)
		}

		group.Tracks = append(group.Tracks, track)
	}

	// Build final ordered slice
	result := make([]*AlbumGroup, 0, len(orderedAlbums))
	for _, id := range orderedAlbums {
		result = append(result, albumMap[id])
	}

	return result
}
