package config

import (
	"os"
	"path/filepath"
)

const APP_NAME = "spotify-playlist-sorter"

func getAppDataDir() (string, error) {
	var baseDir string
	var err error

	// os.UserConfigDir returns:
	// - Windows: C:\Users\<User>\AppData\Roaming
	// - Linux: $XDG_CONFIG_HOME or ~/.config
	// - macOS: ~/Library/Application Support
	baseDir, err = os.UserConfigDir()
	if err != nil {
		// fallback to home dir
		home, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return "", homeErr
		}
		baseDir = home
	}

	// Combine baseDir with your app name
	appDir := filepath.Join(baseDir, APP_NAME)

	// Ensure directory exists
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return "", err
	}

	return appDir, nil
}
