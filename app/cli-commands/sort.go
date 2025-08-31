package clicommands

import (
	"fmt"
	"log"
	"slices"
	"sort"
	"time"

	"github.com/ItsOnlyGame/my-spotify-playlist-sorter-go/app/playlists"
	"github.com/schollz/progressbar/v3"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/urfave/cli/v2"
	"github.com/zmb3/spotify/v2"
)

type SortCommand struct {
	Sp *spotify.Client
}

func (cmd *SortCommand) New() *cli.Command {
	return &cli.Command{
		Name:  "sort",
		Usage: "Sort a playlist",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "playlist",
				Usage:    "The playlist to sort",
				Aliases:  []string{"p"},
				Required: true,
			},
		},
		Action: cmd.Run,
	}
}

func (cmd *SortCommand) Run(ctx *cli.Context) error {
	fmt.Printf("Sort command\n\n")

	playlistUrl := ctx.String("playlist")
	playlistId := spotify.ID(playlists.ExtractSpotifyID(playlistUrl))

	tracks, err := playlists.GetPlaylistTracks(cmd.Sp, playlistId)

	if err != nil {
		log.Fatal(err)
	}

	// This is has been removed from the Spotify Web API
	// It is not possible to remove duplicates from a playlist the way that I would like it to be done
	// I would like to preserve the "added to playlist date" so this features will be on hold for now
	//
	// tracks, err = cmd.removeDuplicateTracksFromPlaylist(ctx, playlistId, tracks)
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }

	var sortedTracks []*spotify.PlaylistItem

	// Step 1: Group tracks by artist
	bar := progressbar.Default(int64(len(tracks)), "Grouping tracks by artist")
	artistTracksMap := orderedmap.NewOrderedMap[string, []*spotify.PlaylistItem]()
	for _, track := range tracks {
		bar.Add(1)

		artistName := track.Track.Track.Artists[0].Name
		artistTracksMap.Set(artistName, append(artistTracksMap.GetOrDefault(artistName, []*spotify.PlaylistItem{}), track))
	}

	// Step 2: Sort tracks for each artist by album, then by release date
	bar = progressbar.Default(int64(artistTracksMap.Len()), "Sort tracks by album and release date")

	for _, artistTracks := range artistTracksMap.Iterator() {
		// Create a map to group tracks by album with release dates
		albumTracksMap := orderedmap.NewOrderedMap[string, []*spotify.PlaylistItem]()
		albumReleaseDates := orderedmap.NewOrderedMap[string, time.Time]()

		for _, track := range artistTracks {
			albumName := track.Track.Track.Album.Name
			albumReleaseDate := track.Track.Track.Album.ReleaseDateTime()

			// Add to the album tracks map
			albumTracksMap.Set(albumName, append(albumTracksMap.GetOrDefault(albumName, []*spotify.PlaylistItem{}), track))

			// Store the release date for sorting
			if existingDate, exists := albumReleaseDates.Get(albumName); !exists || albumReleaseDate.Before(existingDate) {
				albumReleaseDates.Set(albumName, albumReleaseDate)
			}
		}

		// Step 3: Create a slice to sort the albums by release date
		type albumInfo struct {
			name string
			date time.Time
		}

		var albums []albumInfo
		for album, releaseDate := range albumReleaseDates.Iterator() {
			albums = append(albums, albumInfo{name: album, date: releaseDate})
		}

		// Sort albums by release date
		sort.SliceStable(albums, func(i, j int) bool {
			return albums[i].date.Before(albums[j].date)
		})

		// Step 4: Append tracks in the order of albums and their release dates
		for _, album := range albums {
			sortedTracks = append(sortedTracks, albumTracksMap.GetOrDefault(album.name, []*spotify.PlaylistItem{})...) // Append all tracks of this album
		}
		bar.Add(1)
	}

	// Step 5: Reorder the playlist
	err = cmd.reorderPlaylistTracks(ctx, playlistId, tracks, sortedTracks)
	if err != nil {
		return fmt.Errorf("failed to reorder tracks: %w", err)
	}

	fmt.Println("Successfully reordered the playlist.")
	return nil
}

func (cmd *SortCommand) reorderPlaylistTracks(ctx *cli.Context, playlistID spotify.ID, unsortedTracks, sortedTracks []*spotify.PlaylistItem) error {
	bar := progressbar.Default(int64(len(sortedTracks)), "Applying changes to the playlist")

	var snapshotID string

	// Track which positions are out of place
	i := 0
	for i < len(unsortedTracks) {
		if unsortedTracks[i].Track.Track.URI == sortedTracks[i].Track.Track.URI {
			i++
			bar.Add(1)
			continue
		}

		targetIndex := cmd.findIndex(unsortedTracks, sortedTracks[i])

		if targetIndex == -1 {
			log.Panicf("Track %s not found in the current playlist", sortedTracks[i].Track.Track.URI)
		}

		// Create reorder options for a single track
		reorderOptions := spotify.PlaylistReorderOptions{
			RangeStart:   spotify.Numeric(targetIndex),
			RangeLength:  spotify.Numeric(1),
			InsertBefore: spotify.Numeric(i),
		}

		if snapshotID != "" {
			reorderOptions.SnapshotID = snapshotID
		}

		newSnapshotID, err := cmd.Sp.ReorderPlaylistTracks(ctx.Context, playlistID, reorderOptions)
		if err != nil {
			return fmt.Errorf("failed to reorder track %s from position %d: %w", sortedTracks[i].Track.Track.URI, targetIndex, err)
		}

		temp := unsortedTracks[targetIndex]
		unsortedTracks = slices.Delete(unsortedTracks, targetIndex, targetIndex+1)
		unsortedTracks = slices.Insert(unsortedTracks, i, temp)

		// Update the snapshotID for the next reordering operation
		snapshotID = newSnapshotID

		bar.Add(1)
		i++
	}

	return nil
}

func (cmd *SortCommand) removeDuplicateTracksFromPlaylist(ctx *cli.Context, playlistID spotify.ID, tracks []*spotify.PlaylistItem) ([]*spotify.PlaylistItem, error) {
	bar := progressbar.Default(int64(len(tracks)*2), "Identify and remove duplicate tracks from playlist")

	// Step 1: Identify duplicate tracks
	trackSet := make(map[spotify.ID]struct{}, len(tracks))
	duplicates := []spotify.TrackToRemove{}
	uniqueTracks := make([]*spotify.PlaylistItem, 0, len(tracks))

	for index, track := range tracks {
		if track.Track.Track == nil {
			continue
		}

		uri := spotify.ID(track.Track.Track.URI)

		if _, exists := trackSet[uri]; exists {
			duplicates = append(duplicates, spotify.TrackToRemove{URI: uri.String(), Positions: []int{index}})
		} else {
			trackSet[uri] = struct{}{}
			uniqueTracks = append(uniqueTracks, track)
		}

		bar.Add(1)
	}

	bar.Add(len(uniqueTracks))

	// Step 2: Remove duplicate tracks from the playlist
	if len(duplicates) > 0 {

		for _, track := range duplicates {
			fmt.Printf("%s %d \n\n", track.URI, track.Positions[0])
		}

		_, err := cmd.Sp.RemoveTracksFromPlaylistOpt(ctx.Context, playlistID, duplicates, "")
		if err != nil {
			return nil, fmt.Errorf("failed to remove duplicate tracks: %w", err)
		}

		bar.Add(len(duplicates))
		fmt.Printf("Removed %d duplicate tracks from the playlist \n", len(duplicates))
	}

	return uniqueTracks, nil
}

// Helper function to find the index of an element in an array
func (cmd *SortCommand) findIndex(arr []*spotify.PlaylistItem, value *spotify.PlaylistItem) int {
	for i, v := range arr {
		if v.Track.Track.URI == value.Track.Track.URI {
			return i
		}
	}
	return -1
}
