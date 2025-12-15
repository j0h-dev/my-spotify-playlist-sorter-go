package cli

import (
	"log"
	"strings"
	"time"

	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/config"
	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/sorter"
	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/spotify"
	"github.com/mrz1836/go-countries"
	"github.com/urfave/cli/v2"
	spotifyApi "github.com/zmb3/spotify/v2"
)

const SORT_CMD_TITLE = `
   _____ _____ _____
  |   __|  _  |   __|
  |__   |   __|__   |
  |_____|__|  |_____|

 -- Playlist sorting --
`

type SortCommand struct{}

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

func (cmd *SortCommand) Run(ctx *cli.Context) error {
	println(SORT_CMD_TITLE)

	config := config.ReadConfig()
	spotifyClient, err := spotify.Login(config)

	if err != nil {
		log.Fatalf("Failed to login to Spotify: %v", err)
	}

	playlistUrl := ctx.String("playlist")
	playlistId := spotifyApi.ID(spotify.ExtractSpotifyID(playlistUrl))

	country := ctx.String("country")
	if country == "" {
		loc := time.Now().Location().String()
		capital := strings.Split(loc, "/")[1]
		country = countries.GetByCapital(capital).Alpha2
	}

	progress := NewProgressBar()
	sorter.Run(spotifyClient, playlistId, country, progress)

	return nil
}
