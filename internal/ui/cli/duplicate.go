package cli

import (
	"fmt"

	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/config"
	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/domain"
	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/duplicator"
	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/spotify"
	"github.com/urfave/cli/v2"
	spotifyApi "github.com/zmb3/spotify/v2"
)

const DUPLICATE_CMD_TITLE = `
     _____ _____ _____
    |   __|  _  |   __|
    |__   |   __|__   |
    |_____|__|  |_____|

 -- Playlist duplicating --
`

type DuplicateCommand struct{}

func (cmd *DuplicateCommand) New() *cli.Command {
	return &cli.Command{
		Name:  "duplicate",
		Usage: "Duplicate a playlist",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "playlist",
				Usage:    "URL of the playlist to duplicate",
				Aliases:  []string{"p"},
				Required: true,
			},
			&cli.StringFlag{
				Name:        "name",
				Usage:       "Name of the new playlist",
				Aliases:     []string{"n"},
				DefaultText: "Copy of <playlist name>",
				Required:    false,
			},
			&cli.StringFlag{
				Name:        "description",
				Usage:       "Description of the new playlist",
				Aliases:     []string{"d"},
				DefaultText: "Copy of playlist <playlist name>",
				Required:    false,
			},
			&cli.BoolFlag{
				Name:        "public",
				Usage:       "Whether the new playlist should be public",
				DefaultText: "false",
				Required:    false,
			},
			&cli.BoolFlag{
				Name:        "collaborative",
				Usage:       "Whether the new playlist should be collaborative",
				DefaultText: "false",
				Required:    false,
			},
			&cli.StringFlag{
				Name:     "country",
				Usage:    "The country code for market-specific track availability (e.g., US, GB, DE). If not provided, it defaults to the user's locale.",
				Aliases:  []string{"c"},
				Required: false,
			},
		},
		Action: cmd.Run,
	}
}

func (cmd *DuplicateCommand) Run(ctx *cli.Context) error {
	println(DUPLICATE_CMD_TITLE)

	config := config.ReadConfig()
	spotifyClient, err := spotify.Login(config)

	if err != nil {
		fmt.Printf("Failed to login to Spotify: %v\n", err)
	}

	country := ctx.String("country")
	if country == "" {
		country, err = domain.GetDeviceCountry()

		if err != nil {
			return fmt.Errorf("failed to determine device country: %w", err)
		}
	}

	playlistUrl := ctx.String("playlist")
	playlistID := spotifyApi.ID(spotify.ExtractSpotifyID(playlistUrl))

	playlist, err := spotify.GetPlaylist(spotifyClient, playlistID, country, nil)
	if err != nil {
		return fmt.Errorf("failed to get playlist: %w", err)
	}

	playlistName := ctx.String("name")
	if playlistName == "" {
		playlistName = fmt.Sprintf("Copy of %s", playlist.Name)
	}

	description := ctx.String("description")
	if description == "" {
		description = fmt.Sprintf("Copy of playlist: %s", playlist.Name)
	}

	isPlaylistPublic := ctx.Bool("public")
	isPlaylistCollaborative := ctx.Bool("collaborative")

	duplicator.Run(spotifyClient, playlistID, country, duplicator.PlaylistInfo{
		Name:          playlistName,
		Description:   description,
		Public:        isPlaylistPublic,
		Collaborative: isPlaylistCollaborative,
	})

	fmt.Printf("Successfully cloned playlist: %s\n", playlistName)

	return nil
}
