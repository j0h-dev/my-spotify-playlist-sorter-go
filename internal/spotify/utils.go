package spotify

import (
	"fmt"
	"regexp"
)

func extractPortFromURL(url string) (string, error) {
	re := regexp.MustCompile(`:(\d+)`)
	match := re.FindStringSubmatch(url)

	if len(match) > 1 {
		return match[1], nil
	}
	return "", fmt.Errorf("no port found in the URL")
}
