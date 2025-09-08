package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

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

func Login() (*spotify.Client, error) {
	credentials := readSpotifyCredentials()
	if credentials == nil {
		var spotifyId string
		for spotifyId == "" {
			fmt.Print("Enter your Spotify Client ID: ")
			fmt.Scanln(&spotifyId)
		}

		var spotifySecret string
		for spotifySecret == "" {
			fmt.Print("Enter your Spotify Client Secret: ")
			fmt.Scanln(&spotifySecret)
		}

		redirectLink := "http://localhost:8080/api/auth"
		fmt.Print("Enter redirect link set in for you app (default: http://localhost:8080/api/auth): ")
		fmt.Scanln(&redirectLink)

		credentials = &SpotifyCredentials{
			ClientID:     spotifyId,
			ClientSecret: spotifySecret,
			RedirectURI:  redirectLink,
		}

		saveSpotifyCredentials(credentials)
	}

	spotifyAuthenticator = spotifyauth.New(
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopePlaylistReadPrivate,
			spotifyauth.ScopePlaylistModifyPrivate,
			spotifyauth.ScopePlaylistModifyPublic,
			spotifyauth.ScopePlaylistReadCollaborative,
			spotifyauth.ScopeUserLibraryModify,
			spotifyauth.ScopeUserLibraryRead,
		),
		spotifyauth.WithRedirectURL(credentials.RedirectURI),
		spotifyauth.WithClientID(credentials.ClientID),
		spotifyauth.WithClientSecret(credentials.ClientSecret),
	)

	token := readOAuthToken()
	if token != nil {
		spotifyClient = spotify.New(spotifyAuthenticator.Client(context.Background(), token))
		return spotifyClient, nil
	}

	var wg sync.WaitGroup
	wg.Add(1)

	url := spotifyAuthenticator.AuthURL(STATE)
	fmt.Println("Please log in to Spotify by visiting the following page: \n", url)

	port, err := extractPortFromURL(credentials.RedirectURI)
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

		saveOAuthToken(token)
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
