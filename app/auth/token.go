package auth

import (
	"encoding/json"
	"os"
	"path"

	"github.com/j0h-dev/my-spotify-playlist-sorter-go/app/appdata"
	"golang.org/x/oauth2"
)

const OAUTH_TOKEN_FILE = "spotify_token.json"
const SPOTIFY_CREDENTIALS_FILE = "spotify_credentials.json"

func readOAuthToken() *oauth2.Token {
	token := &oauth2.Token{}

	appdata, err := appdata.GetAppDataDir()
	if err != nil {
		return nil
	}

	jsonData, err := os.ReadFile(path.Join(appdata, OAUTH_TOKEN_FILE))
	if err != nil {
		return nil
	}

	err = json.Unmarshal(jsonData, token)
	if err != nil {
		return nil
	}

	return token
}

func saveOAuthToken(token *oauth2.Token) error {
	jsonData, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return err
	}

	appdata, err := appdata.GetAppDataDir()
	if err != nil {
		return nil
	}

	return os.WriteFile(path.Join(appdata, OAUTH_TOKEN_FILE), jsonData, 0644)
}

type SpotifyCredentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
}

func readSpotifyCredentials() *SpotifyCredentials {
	credentials := &SpotifyCredentials{}

	appdata, err := appdata.GetAppDataDir()
	if err != nil {
		return nil
	}

	jsonData, err := os.ReadFile(path.Join(appdata, SPOTIFY_CREDENTIALS_FILE))
	if err != nil {
		return nil
	}

	err = json.Unmarshal(jsonData, credentials)
	if err != nil {
		return nil
	}

	return credentials
}

func saveSpotifyCredentials(credentials *SpotifyCredentials) error {
	jsonData, err := json.MarshalIndent(credentials, "", "  ")
	if err != nil {
		return err
	}

	appdata, err := appdata.GetAppDataDir()
	if err != nil {
		return nil
	}

	return os.WriteFile(path.Join(appdata, SPOTIFY_CREDENTIALS_FILE), jsonData, 0644)
}
