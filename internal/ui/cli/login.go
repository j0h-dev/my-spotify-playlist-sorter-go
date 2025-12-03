package cli

import (
	"fmt"

	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/domain"
	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/spotify"
	"github.com/urfave/cli/v2"
)

const LOGIN_CMD_TITLE = `
     _____ _____ _____
    |   __|  _  |   __|
    |__   |   __|__   |
    |_____|__|  |_____|

        -- Login --
`

type LoginCommand struct{}

func (cmd *LoginCommand) New() *cli.Command {
	return &cli.Command{
		Name:  "login",
		Usage: "Login to Spotify",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Usage:    "Set the client ID",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "secret",
				Usage:    "Set the client secret",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "redirect",
				Usage:    "Set the redirect URI",
				Required: true,
			},
		},
		Action: cmd.Run,
	}
}

func (cmd *LoginCommand) Run(ctx *cli.Context) error {
	println(LOGIN_CMD_TITLE)

	clientID := ctx.String("id")
	clientSecret := ctx.String("secret")
	redirectURI := ctx.String("redirect")

	fmt.Printf("Login details:\n")
	fmt.Printf("Client ID: %s\n", cmd.hideSecret(clientID))
	fmt.Printf("Client Secret: %s\n", cmd.hideSecret(clientSecret))
	fmt.Printf("Redirect URI: %s\n\n", cmd.hideSecret(redirectURI))

	err := spotify.Authenticate(&domain.SpotifyCredentials{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	})

	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	return nil
}

func (cmd *LoginCommand) hideSecret(secret string) string {
	if len(secret) <= 4 {
		return "****"
	}
	return secret[:2] + "****" + secret[len(secret)-2:]
}
