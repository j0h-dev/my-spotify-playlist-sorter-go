package sorter

import (
	"sort"
)

// SortTracksInsideAlbums sorts the tracks inside each album group by disc number and track number.
// It updates the provided AlbumGroup slices in place.
func SortTracksInsideAlbums(
	groups []*AlbumGroup,
) {
	for idx := range groups {
		tracks := groups[idx].Tracks

		sort.SliceStable(tracks, func(i, j int) bool {
			a := tracks[i].Track.Track
			b := tracks[j].Track.Track

			// First: disc number (for multi-disc albums)
			if a.DiscNumber != b.DiscNumber {
				return a.DiscNumber < b.DiscNumber
			}

			// Then: track number
			return a.TrackNumber < b.TrackNumber
		})

		groups[idx].Tracks = tracks
	}
}
