package main

import (
	"os"

	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/ui/cli"
	cliapp "github.com/urfave/cli/v2"
)

type CliCommands struct {
	sortCommand      *cli.SortCommand
	duplicateCommand *cli.DuplicateCommand
	loginCommand     *cli.LoginCommand
}

func main() {
	cliCommands := &CliCommands{
		sortCommand:      &cli.SortCommand{},
		duplicateCommand: &cli.DuplicateCommand{},
		loginCommand:     &cli.LoginCommand{},
	}

	app := &cliapp.App{
		Usage: "Spotify Playlist Sorter CLI",
		Commands: []*cliapp.Command{
			cliCommands.sortCommand.New(),
			cliCommands.duplicateCommand.New(),
			cliCommands.loginCommand.New(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
