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

func ListenForToken(
	auth *spotifyauth.Authenticator,
	redirectURI string,
	callback func(error, *oauth2.Token),
) {

	port, err := extractPortFromURL(redirectURI)
	if err != nil {
		callback(err, nil)
		return
	}

	server := &http.Server{Addr: fmt.Sprintf(":%s", port)}

	wg := sync.WaitGroup{}
	wg.Add(1)

	http.HandleFunc("/api/auth", func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.Token(r.Context(), STATE, r)
		if err != nil {
			http.Error(w, "Failed to obtain token", http.StatusInternalServerError)
			callback(err, nil)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Authentication successful! You can close this window.")

		go func() {
			defer wg.Done()
			server.Close()
			callback(nil, token)
		}()
	})

	go func() {
		server.ListenAndServe()
	}()
}

func Authenticate(
	credentials *domain.SpotifyCredentials,
	callback func(error),
) (loginUrl string) {
	auth := NewAuthenticator(credentials)
	loginUrl = LoginURL(auth)

	ListenForToken(auth, credentials.RedirectURI, func(err error, token *oauth2.Token) {
		if err != nil {
			callback(err)
			return
		}

		conf := config.ReadConfig()
		conf.OAuthToken = token
		conf.Credentials = credentials

		if err := config.SaveConfig(conf); err != nil {
			callback(err)
			return
		}

		callback(nil)
	})

	return loginUrl
}

func Login(cfg *conf.Config) (*spotify.Client, error) {
	auth := NewAuthenticator(cfg.Credentials)

	if cfg.OAuthToken == nil {
		return nil, fmt.Errorf("no OAuth token found, please login first")
	}

	return ClientWithToken(auth, cfg.OAuthToken), nil
}
