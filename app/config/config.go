package config

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Spotify struct {
		// This can be left as is, and it should work.
		// If propblems occur then change the port to something else
		Redirect string `envconfig:"SPOTIFY_REDIRECT" validate:"required,url"`
	}
}

func Load() *Config {
	cfg := &Config{}

	// Read from environment variables
	err := readEnv(cfg)
	if err != nil {
		log.Panic(err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		log.Panicf("Configuration validation failed: %s", err)
	}

	log.Printf("%+v", cfg)

	return cfg
}

func readEnv(cfg *Config) error {
	err := envconfig.Process("", cfg)
	if err != nil {
		return err
	}

	return nil
}
