package sorter

import (
	"sort"
)

// SortAlbumsByReleaseDate sorts the provided AlbumGroup slices by their release dates in ascending order.
// It updates the provided AlbumGroup slices in place.
func SortAlbumsByReleaseDate(
	groups []*AlbumGroup,
) {
	sort.SliceStable(groups, func(i, j int) bool {
		return groups[i].ReleaseDate.Before(groups[j].ReleaseDate)
	})
}
