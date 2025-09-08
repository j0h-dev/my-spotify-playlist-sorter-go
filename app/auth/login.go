package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/ItsOnlyGame/my-spotify-playlist-sorter-go/app/config"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const (
	STATE = "main"
)

var (
	spotifyAuthenticator *spotifyauth.Authenticator
	spotifyClient        *spotify.Client
)

type LoginParams struct {
	Redirect string
}

func Login(config *config.Config) (*spotify.Client, error) {
	spotifyAuthenticator = spotifyauth.New(
		spotifyauth.WithRedirectURL(config.Spotify.Redirect),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate),
	)

	token, _ := getSpotifyToken()

	ctx := context.Background()

	if token != nil {
		spotifyClient = spotify.New(spotifyAuthenticator.Client(ctx, token))
		return spotifyClient, nil
	}

	var wg sync.WaitGroup
	wg.Add(1)

	url := spotifyAuthenticator.AuthURL("main")
	fmt.Println("Please log in to Spotify by visiting the following page: \n", url)

	port, err := extractPortFromURL(config.Spotify.Redirect)
	if err != nil {
		return nil, err
	}

	server := &http.Server{Addr: fmt.Sprintf(":%s", port)}

	http.HandleFunc("/api/auth", func(w http.ResponseWriter, r *http.Request) {
		success := "Log in successful! You can go back to the application :)"

		fmt.Fprintln(w, success)
		log.Println(success)

		token, err := spotifyAuthenticator.Token(r.Context(), STATE, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusNotFound)
			return
		}

		saveSpotifyToken(token)
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

	return spotifyClient, nil
}
