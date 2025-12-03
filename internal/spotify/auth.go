package spotify

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/config"
	conf "github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/config"
	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/domain"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

const (
	STATE = "main"
)

func NewAuthenticator(creds *domain.SpotifyCredentials) *spotifyauth.Authenticator {
	return spotifyauth.New(
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopePlaylistReadPrivate,
			spotifyauth.ScopePlaylistModifyPrivate,
			spotifyauth.ScopePlaylistModifyPublic,
			spotifyauth.ScopePlaylistReadCollaborative,
			spotifyauth.ScopeUserLibraryModify,
			spotifyauth.ScopeUserLibraryRead,
		),
		spotifyauth.WithRedirectURL(creds.RedirectURI),
		spotifyauth.WithClientID(creds.ClientID),
		spotifyauth.WithClientSecret(creds.ClientSecret),
	)
}

func ClientWithToken(
	auth *spotifyauth.Authenticator,
	token *oauth2.Token,
) *spotify.Client {
	return spotify.New(auth.Client(context.Background(), token))
}

func LoginURL(auth *spotifyauth.Authenticator) string {
	return auth.AuthURL(STATE)
}

func WaitForToken(
	auth *spotifyauth.Authenticator,
	redirectURI string,
) (*oauth2.Token, error) {

	port, err := extractPortFromURL(redirectURI)
	if err != nil {
		return nil, err
	}

	server := &http.Server{Addr: fmt.Sprintf(":%s", port)}

	wg := sync.WaitGroup{}
	wg.Add(1)

	var token *oauth2.Token

	http.HandleFunc("/api/auth", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Login successful! You can return to the app.")

		token, err = auth.Token(r.Context(), STATE, r)
		if err != nil {
			http.Error(w, "Failed to obtain token", http.StatusInternalServerError)
			return
		}

		go func() {
			defer wg.Done()
			server.Close()
		}()
	})

	go func() {
		server.ListenAndServe()
	}()

	wg.Wait()
	return token, nil
}

func Authenticate(
	credentials *domain.SpotifyCredentials,
) error {
	auth := NewAuthenticator(credentials)

	loginUrl := LoginURL(auth)
	fmt.Printf("Please log in to Spotify by visiting the following URL:\n\n%s\n\n", loginUrl)

	token, err := WaitForToken(auth, credentials.RedirectURI)
	if err != nil {
		return err
	}

	conf := config.ReadConfig()
	conf.OAuthToken = token
	conf.Credentials = credentials

	if err := config.SaveConfig(conf); err != nil {
		return err
	}

	return nil
}

func Login(cfg *conf.Config) (*spotify.Client, error) {
	auth := NewAuthenticator(cfg.Credentials)

	if cfg.OAuthToken == nil {
		return nil, fmt.Errorf("no OAuth token found, please login first")
	}

	return ClientWithToken(auth, cfg.OAuthToken), nil
}
