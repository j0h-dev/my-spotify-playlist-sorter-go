package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	urlutils "github.com/ItsOnlyGame/my-spotify-playlist-sorter-go/internal/utils/url-utils"
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

func Login() (*spotify.Client, error) {
	spotifyAuthenticator = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURL),
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

	redirect_url, exists := os.LookupEnv("SPOTIFY_REDIRECT")
	if !exists {
		log.Fatal("SPOTIFY_REDIRECT environment variable not set")
	}

	port, err := urlutils.ExtractPortFromURL(redirect_url)
	if err != nil {
		return nil, err
	}

	server := &http.Server{Addr: fmt.Sprintf(":%s", port)}

	http.HandleFunc("/api/auth", func(w http.ResponseWriter, r *http.Request) {
		success := "Log in successful! You can go back to the application :)"

		fmt.Fprintln(w, success)
		log.Println(success)

		token, err := spotifyAuthenticator.Token(r.Context(), state, r)
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
