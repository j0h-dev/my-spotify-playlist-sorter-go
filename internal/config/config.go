package config

import (
	"encoding/json"
	"os"
	"path"

	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/domain"
	"golang.org/x/oauth2"
)

const CONFIG_FILE = "config.json"

type Config struct {
	OAuthToken  *oauth2.Token              `json:"oauth_token"`
	Credentials *domain.SpotifyCredentials `json:"credentials"`
}

func ReadConfig() *Config {
	config := &Config{}

	appdata, err := getAppDataDir()
	if err != nil {
		return config
	}

	jsonData, err := os.ReadFile(path.Join(appdata, CONFIG_FILE))
	if err != nil {
		return config
	}

	err = json.Unmarshal(jsonData, config)
	if err != nil {
		return config
	}

	return config
}

func SaveConfig(config *Config) error {
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	appdata, err := getAppDataDir()
	if err != nil {
		return nil
	}

	return os.WriteFile(path.Join(appdata, CONFIG_FILE), jsonData, 0644)
}
