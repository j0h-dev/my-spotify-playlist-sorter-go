package domain

import "github.com/zmb3/spotify/v2"

// Helper function to find the index of an element in an array
func FindIndex(arr []*spotify.PlaylistItem, value *spotify.PlaylistItem) int {
	for i, v := range arr {
		if v.Track.Track.URI == value.Track.Track.URI {
			return i
		}
	}
	return -1
}
