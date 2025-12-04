package spotify

import "strings"

func ExtractSpotifyID(url string) string {
	parts := strings.Split(url, "/")

	if len(parts) >= 5 {
		idWithQuery := parts[len(parts)-1]
		id := strings.Split(idWithQuery, "?")[0]
		return id
	}

	return url
}
