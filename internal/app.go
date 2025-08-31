package internal

import (
	"log"
	"os"

	"github.com/ItsOnlyGame/my-spotify-playlist-sorter-go/internal/auth"
	clicommands "github.com/ItsOnlyGame/my-spotify-playlist-sorter-go/internal/cli-commands"
	"github.com/ItsOnlyGame/my-spotify-playlist-sorter-go/internal/config"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"github.com/zmb3/spotify/v2"
)

func Run() {
	// Load environment variables from .env file
	godotenv.Load()

	config := config.Load()

	client, err := auth.Login(config)
	if err != nil {
		log.Fatalf("Failed to login to Spotify: %v", err)
	}

	err = createCli(client)
	if err != nil {
		log.Fatalf("Failed to run CLI: %v", err)
	}
}

type CliCommands struct {
	sortCommand      *clicommands.SortCommand
	duplicateCommand *clicommands.DuplicateCommand
}

func createCli(client *spotify.Client) error {
	cliCommands := &CliCommands{
		sortCommand: &clicommands.SortCommand{
			Sp: client,
		},
		duplicateCommand: &clicommands.DuplicateCommand{
			Sp: client,
		},
	}

	app := &cli.App{
		Name:  "sps",
		Usage: "Spotify Playlist Sorter CLI",
		Commands: []*cli.Command{
			cliCommands.sortCommand.New(),
			cliCommands.duplicateCommand.New(),
		},
	}

	return app.Run(os.Args)
}
