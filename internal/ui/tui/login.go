package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/domain"
	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/spotify"
)

const (
	LoginClientIdFieldFocus = iota
	LoginClientSecretFieldFocus
	LoginRedirectFieldFocus
	LoginButtonsFocus
)

const (
	LoginSubmitButtonFocus = iota
	LoginQuitButtonFocus
)

type LoginModel struct {
	cursor int
	items  []string

	clientIdField       textinput.Model
	clientSecretField   textinput.Model
	clientRedirectField textinput.Model

	errMsg string

	focusIndex       int
	buttonFocusIndex int

	loadingSpinner spinner.Model

	waitingForAuth bool
	authLink       string
	toMenu         bool
}

func NewLoginModel() *LoginModel {
	// initialize text inputs
	clientIdField := textinput.New()
	clientIdField.Placeholder = "client id"
	clientIdField.Prompt = "Client ID: "
	clientIdField.Width = 80
	clientIdField.Focus()

	clientSecretField := textinput.New()
	clientSecretField.Placeholder = "client secret"
	clientSecretField.Prompt = "Client Secret: "
	clientSecretField.Width = 80

	clientRedirectField := textinput.New()
	clientRedirectField.Placeholder = "http://127.0.0.1:8000/api/auth"
	clientRedirectField.Prompt = "Redirect URI: "
	clientRedirectField.Width = 80

	loadingSpinner := spinner.New()
	loadingSpinner.Spinner = spinner.Dot
	loadingSpinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	clientIdField.SetValue("b30517e8fa8444e78583cdd7021effc3")
	clientSecretField.SetValue("abbb803cdb734aaf88d9ba1eff344fe2")
	clientRedirectField.SetValue("http://127.0.0.1:8000/api/auth")

	return &LoginModel{
		clientIdField:       clientIdField,
		clientSecretField:   clientSecretField,
		clientRedirectField: clientRedirectField,

		loadingSpinner: loadingSpinner,

		focusIndex: LoginClientIdFieldFocus,
	}
}

func (m *LoginModel) Init() tea.Cmd {
	return m.loadingSpinner.Tick
}

func (m *LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.loadingSpinner, cmd = m.loadingSpinner.Update(msg)
		return m, cmd

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			if m.focusIndex == LoginButtonsFocus {
				switch m.buttonFocusIndex {
				case LoginQuitButtonFocus:
					return m, tea.Quit

				case LoginSubmitButtonFocus:

					if m.clientIdField.Value() == "" ||
						m.clientSecretField.Value() == "" ||
						m.clientRedirectField.Value() == "" {

						m.errMsg = "All fields are required."
						return m, nil
					}

					// Save credentials
					authLink := spotify.Authenticate(&domain.SpotifyCredentials{
						ClientID:     m.clientIdField.Value(),
						ClientSecret: m.clientSecretField.Value(),
						RedirectURI:  m.clientRedirectField.Value(),
					}, func(err error) {
						if err != nil {
							panic(fmt.Errorf("failed to authenticate: %w", err))
						}

						m.toMenu = true
					})

					m.waitingForAuth = true
					m.authLink = authLink

					return m, nil
				}
			}

		case "tab", "shift+tab", "down", "up":
			// Cycle between fields
			if msg.String() == "down" || msg.String() == "tab" {
				m.focusIndex = (m.focusIndex + 1) % 4
			}
			if msg.String() == "up" || msg.String() == "shift+tab" {
				m.focusIndex = (m.focusIndex + 3) % 4
			}

			m.blurFields()

		case "left", "right":
			if m.focusIndex == LoginButtonsFocus {
				if msg.String() == "right" {
					m.buttonFocusIndex = (m.buttonFocusIndex + 1) % 2
				}
				if msg.String() == "left" {
					m.buttonFocusIndex = (m.buttonFocusIndex + 1) % 2
				}
			}

		}

	}

	// Update text inputs
	switch m.focusIndex {
	case LoginClientIdFieldFocus:
		m.clientIdField, cmd = m.clientIdField.Update(msg)
	case LoginClientSecretFieldFocus:
		m.clientSecretField, cmd = m.clientSecretField.Update(msg)
	case LoginRedirectFieldFocus:
		m.clientRedirectField, cmd = m.clientRedirectField.Update(msg)
	}

	return m, cmd
}

func (m *LoginModel) blurFields() {
	m.clientIdField.Blur()
	m.clientSecretField.Blur()
	m.clientRedirectField.Blur()

	switch m.focusIndex {
	case LoginClientIdFieldFocus:
		m.clientIdField.Focus()
	case LoginClientSecretFieldFocus:
		m.clientSecretField.Focus()
	case LoginRedirectFieldFocus:
		m.clientRedirectField.Focus()
	}
}

func (m *LoginModel) View() string {
	s := fmt.Sprintf("%s\n", APP_TITLE)

	if m.waitingForAuth {
		s += m.loadingSpinner.View()
		s += " Please log in to Spotify by visiting the following URL:\n\n"
		s += lipgloss.NewStyle().Width(100).Render(m.authLink)
		return s
	}

	s += m.clientIdField.View() + "\n"
	s += m.clientSecretField.View() + "\n"
	s += m.clientRedirectField.View() + "\n\n"

	// Submit button render
	if m.focusIndex == LoginButtonsFocus {
		switch m.buttonFocusIndex {
		case LoginSubmitButtonFocus:
			s += "[ Login ] - Quit -\n"
		case LoginQuitButtonFocus:
			s += "- Login - [ Quit ]\n"
		}
	} else {
		s += "- Login - - Quit -\n"
	}
	if m.errMsg != "" {
		s += "\n" + m.errMsg
	}

	s += "\nTab to switch fields, Enter to submit."

	return s
}
