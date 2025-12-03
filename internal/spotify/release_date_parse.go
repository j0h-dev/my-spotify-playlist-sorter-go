package spotify

import "time"

func ParseSpotifyReleaseDate(date string, precision string) (time.Time, error) {
	var strDate string

	if precision == "year" {
		strDate = date + "-01-01"
	}

	if precision == "month" {
		strDate = date + "-01-01"
	}

	if precision == "day" {
		strDate = date
	}

	return time.Parse("2006-01-02", strDate)
}
