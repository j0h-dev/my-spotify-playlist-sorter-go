package spotify

import (
	"context"
	"fmt"

	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/ui"
	"github.com/zmb3/spotify/v2"
)

func CreatePlaylist(
	sp *spotify.Client,
	userID string,
	playlistName string,
	description string,
	isPlaylistPublic bool,
	isPlaylistCollaborative bool,
) (*spotify.FullPlaylist, error) {
	playlist, err := sp.CreatePlaylistForUser(
		context.Background(),
		userID,
		playlistName,
		description,
		isPlaylistPublic,
		isPlaylistCollaborative,
	)

	if err != nil {
		return nil, err
	}

	return playlist, nil
}

func GetPlaylist(
	sp *spotify.Client,
	playlistId spotify.ID,
	country string,
	progress ui.Progress,
) (*spotify.FullPlaylist, error) {
	ctx := context.Background()

	progress.Start(1, "Fetching playlist details")

	playlist, err := sp.GetPlaylist(ctx, playlistId, spotify.RequestOption(spotify.Country(country)))

	if err != nil {
		return nil, err
	}

	progress.Finish()
	return playlist, nil
}

func GetPlaylistTracks(
	sp *spotify.Client,
	playlist *spotify.FullPlaylist,
	country string,
	progress ui.Progress,
) ([]*spotify.PlaylistItem, error) {
	ctx := context.Background()

	if progress != nil {
		progress.Start(int(playlist.Tracks.Total), "Fetching playlist tracks")
		defer progress.Finish()
	}

	var tracks []*spotify.PlaylistItem
	offset := 0
	limit := 100

	for {
		page, err := sp.GetPlaylistItems(
			ctx,
			playlist.ID,
			spotify.RequestOption(spotify.Offset(offset)),
			spotify.RequestOption(spotify.Limit(limit)),
			spotify.RequestOption(spotify.Country(country)),
		)

		if err != nil {
			return nil, err
		}

		for i := range page.Items {
			tracks = append(tracks, &page.Items[i])
		}

		// Stop when we've fetched all items
		offset += len(page.Items)
		progress.Advance(len(page.Items))
		if offset >= int(page.Total) {
			break
		}
	}

	return tracks, nil
}

func AddTracksToPlaylist(
	client *spotify.Client,
	tracks []*spotify.PlaylistItem,
	playlist *spotify.FullPlaylist,
	progress ui.Progress,
) error {
	if progress != nil {
		progress.Start(len(tracks), "Adding tracks to new playlist")
		defer progress.Finish()
	}

	// Collect track IDs
	// We need to add tracks in batches of 100 due to Spotify API limitations

	batchSize := 100

	for i := 0; i < len(tracks); i += batchSize {
		end := min(i+batchSize, len(tracks))

		slice := make([]spotify.ID, len(tracks[i:end]))
		for j, track := range tracks[i:end] {
			slice[j] = track.Track.Track.ID
		}

		_, err := client.AddTracksToPlaylist(context.Background(), playlist.ID, slice...)

		if err != nil {
			return err
		}

		progress.Advance(len(slice))
	}

	return nil
}

func ReorderPlaylistTracks(
	client *spotify.Client,
	playlistID spotify.ID,
	original []*spotify.PlaylistItem,
	sorted []*spotify.PlaylistItem,
	progress ui.Progress,
) error {
	if progress != nil {
		progress.Start(len(original), "Reordering playlist tracks")
		defer progress.Finish()
	}

	// Build lookup map: URI → current index
	indexMap := make(map[spotify.URI]int, len(original))
	for i, t := range original {
		indexMap[t.Track.Track.URI] = i
	}

	var snapshotID string

	for i := range len(sorted) {
		if progress != nil {
			progress.Advance(1)
		}

		currentURI := original[i].Track.Track.URI
		targetURI := sorted[i].Track.Track.URI

		if currentURI == targetURI {
			continue
		}

		targetIndex, ok := indexMap[targetURI]
		if !ok {
			return fmt.Errorf("track %s not found in current playlist", targetURI)
		}

		// Reorder single item using Spotify API
		opts := spotify.PlaylistReorderOptions{
			RangeStart:   spotify.Numeric(targetIndex),
			RangeLength:  spotify.Numeric(1),
			InsertBefore: spotify.Numeric(i),
		}

		if snapshotID != "" {
			opts.SnapshotID = snapshotID
		}

		newSnapshotID, err := client.ReorderPlaylistTracks(context.Background(), playlistID, opts)
		if err != nil {
			return fmt.Errorf(
				"failed to reorder track %s from position %d: %w",
				targetURI,
				targetIndex,
				err,
			)
		}

		// Locally re-order slice to match Spotify state
		item := original[targetIndex]

		for j := targetIndex; j > i; j-- {
			original[j] = original[j-1]
			indexMap[original[j].Track.Track.URI] = j
		}

		original[i] = item
		indexMap[item.Track.Track.URI] = i

		snapshotID = newSnapshotID
	}

	return nil
}
