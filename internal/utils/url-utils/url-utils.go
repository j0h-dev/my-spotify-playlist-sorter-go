package urlutils

import (
	"fmt"
	"regexp"
	"strings"
)

func ExtractSpotifyID(url string) string {
	parts := strings.Split(url, "/")

	if len(parts) >= 5 {
		idWithQuery := parts[len(parts)-1]
		id := strings.Split(idWithQuery, "?")[0]
		return id
	}

	return url
}

func ExtractPortFromURL(url string) (string, error) {
	re := regexp.MustCompile(`:(\d+)`)
	match := re.FindStringSubmatch(url)

	if len(match) > 1 {
		return match[1], nil
	}
	return "", fmt.Errorf("no port found in the URL")
}
