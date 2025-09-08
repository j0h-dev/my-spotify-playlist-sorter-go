package main

import (
	"log"
	"os"

	"github.com/j0h-dev/my-spotify-playlist-sorter-go/app/auth"
	clicommands "github.com/j0h-dev/my-spotify-playlist-sorter-go/app/cli-commands"
	"github.com/urfave/cli/v2"
	"github.com/zmb3/spotify/v2"
)

func main() {
	client, err := auth.Login()
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
		Usage: "Spotify Playlist Sorter CLI",
		Commands: []*cli.Command{
			cliCommands.sortCommand.New(),
			cliCommands.duplicateCommand.New(),
		},
	}

	return app.Run(os.Args)
}
