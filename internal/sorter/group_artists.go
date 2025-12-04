package sorter

import (
	spotifyApi "github.com/zmb3/spotify/v2"
)

func GroupTracksByArtist(
	tracks []*spotifyApi.PlaylistItem,
) []*ArtistGroup {

	// Fast lookup by artist name
	artistMap := make(map[string]*ArtistGroup, len(tracks))

	// Slice to preserve artist order of first appearance
	order := make([]string, 0)

	for _, track := range tracks {
		artistName := track.Track.Track.Artists[0].Name

		group, exists := artistMap[artistName]
		if !exists {
			group = &ArtistGroup{
				Name:   artistName,
				Tracks: make([]*spotifyApi.PlaylistItem, 0),
			}
			artistMap[artistName] = group
			order = append(order, artistName)
		}

		group.Tracks = append(group.Tracks, track)
	}

	// Convert to ordered slice
	result := make([]*ArtistGroup, 0, len(order))
	for _, name := range order {
		result = append(result, artistMap[name])
	}

	return result
}
