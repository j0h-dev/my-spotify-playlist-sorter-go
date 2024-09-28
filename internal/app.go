package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	clicommands "github.com/ItsOnlyGame/my-spotify-playlist-sorter-go/internal/cli-commands"
	spotifytoken "github.com/ItsOnlyGame/my-spotify-playlist-sorter-go/internal/utils/spotify-token"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const (
	state       = "main"
	redirectURL = "http://localhost:3000/api/auth"
)

var (
	spotifyAuthenticator *spotifyauth.Authenticator
	spotifyClient        *spotify.Client
)

func Run() {
	godotenv.Load()
	handleSpotifyAuthentication()
	createCLIApp()
}

type CliCommands struct {
	sortCommand      *clicommands.SortCommand
	duplicateCommand *clicommands.DuplicateCommand
}

func createCLIApp() {
	cliCommands := &CliCommands{
		sortCommand: &clicommands.SortCommand{
			Sp: spotifyClient,
		},
		duplicateCommand: &clicommands.DuplicateCommand{
			Sp: spotifyClient,
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

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func handleSpotifyAuthentication() {
	spotifyAuthenticator = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURL), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate))

	token, _ := spotifytoken.GetSpotifyToken()
	ctx := context.Background()

	if token != nil {
		spotifyClient = spotify.New(spotifyAuthenticator.Client(ctx, token))
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)

	url := spotifyAuthenticator.AuthURL("main")
	fmt.Println("Please log in to Spotify by visiting the following page: \n", url)

	server := &http.Server{Addr: ":8080"}

	http.HandleFunc("/api/auth", func(w http.ResponseWriter, r *http.Request) {
		success := "Log in successful! You can go back to the application :)"

		fmt.Fprintln(w, success)
		log.Println(success)

		token, err := spotifyAuthenticator.Token(r.Context(), state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusNotFound)
			return
		}

		spotifytoken.SaveSpotifyToken(token)
		spotifyClient = spotify.New(spotifyAuthenticator.Client(r.Context(), token))

		go func() {
			defer wg.Done()
			if err := server.Close(); err != nil {
				log.Fatalf("Failed to close the server: %v", err)
			}
		}()
	})

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	wg.Wait()
}
