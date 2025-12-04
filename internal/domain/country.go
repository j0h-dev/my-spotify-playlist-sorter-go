package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/mrz1836/go-countries"
)

func GetDeviceCountry() (string, error) {
	loc := time.Now().Location().String()

	if loc == "Local" {
		return "", fmt.Errorf("device does not have a valid location set")
	}

	capital := strings.Split(loc, "/")[1]

	return countries.GetByCapital(capital).Alpha2, nil
}
