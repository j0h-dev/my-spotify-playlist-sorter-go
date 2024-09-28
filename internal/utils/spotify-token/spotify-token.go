package spotifytoken

import (
	"encoding/json"
	"os"

	"golang.org/x/oauth2"
)

const tokenFileName = "spotify_token.json"

func GetSpotifyToken() (*oauth2.Token, error) {
	token := &oauth2.Token{}

	err := readJSONFromFile(tokenFileName, token)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func SaveSpotifyToken(token *oauth2.Token) error {
	return writeJSONToFile(tokenFileName, token)
}

func writeJSONToFile(fileName string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, jsonData, 0644)
}

func readJSONFromFile(fileName string, data interface{}) error {
	jsonData, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonData, data)
}
