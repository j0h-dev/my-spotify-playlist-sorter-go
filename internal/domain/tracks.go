package domain

import "github.com/zmb3/spotify/v2"

func FilterNonLocalTracks(
	tracks []*spotify.PlaylistItem,
) []*spotify.PlaylistItem {

	out := make([]*spotify.PlaylistItem, 0, len(tracks))
	for _, t := range tracks {
		if !t.IsLocal {
			out = append(out, t)
		}
	}

	return out
}
