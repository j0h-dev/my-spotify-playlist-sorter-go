package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/config"
	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/duplicator"
	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/spotify"
	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/ui/tui/components"

	spotifyApi "github.com/zmb3/spotify/v2"
)

const (
	DuplicateUrlFieldFocus = iota
	DuplicateNameFieldFocus
	DuplicateDescriptionFieldFocus
	DuplicatePublicToggleFocus
	DuplicateCollaborationToggleFocus
	DuplicateCountryFieldFocus
	DuplicateButtonsFocus
)

const (
	DuplicateSubmitButtonFocus = iota
	DuplicateBackButtonFocus
)

type DuplicateModel struct {
	cursor int
	items  []string

	urlInput            textinput.Model
	nameInput           textinput.Model
	descriptionInput    textinput.Model
	publicToggle        components.Toggle
	collaborationToggle components.Toggle
	countryInput        textinput.Model

	focusIndex       int
	buttonFocusIndex int

	back bool

	progressBar *components.ProgressBar
	duplicating bool
	exit        bool
}

// NewDuplicateModel initializes and returns a new DuplicateModel with default values.
func NewDuplicateModel() *DuplicateModel {
	return &DuplicateModel{
		urlInput:            newTextInput("URL: ", "https://open.spotify.com/playlist/...", 80, true),
		nameInput:           newTextInput("New Playlist Name: ", "Copy of <playlist name>", 50, false),
		descriptionInput:    newTextInput("New Playlist Description: ", "Copy of playlist <playlist name>", 100, false),
		publicToggle:        components.NewToggle("Public", false),
		collaborationToggle: components.NewToggle("Collaborative", false),
		countryInput:        newTextInput("Country: ", "FI, SV, US...", 20, false),
		focusIndex:          DuplicateUrlFieldFocus,
		progressBar:         components.NewProgressBar(),
	}
}

// Init initializes the DuplicateModel and returns the initial command to start ticking.
func (m *DuplicateModel) Init() tea.Cmd {
	return tickCmd()
}

// Update handles incoming messages and updates the state of the DuplicateModel accordingly.
// It returns the updated model and any commands to be executed.
func (m *DuplicateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.exit {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.updateProgressBarWidth(msg.Width)
		return m, nil

	case tickMsg:
		return m, tickCmd()

	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	}

	return m, nil
}

// updateProgressBarWidth adjusts the width of the progress bar based on the window size.
func (m *DuplicateModel) updateProgressBarWidth(windowWidth int) {
	m.progressBar.SetWidth(windowWidth - 2*2 - 4)
	if m.progressBar.Width() > MaxProgressBarWidth {
		m.progressBar.SetWidth(MaxProgressBarWidth)
	}
}

// handleKeyMsg processes key press events and updates the model state accordingly.
func (m *DuplicateModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {

	case "ctrl+c":
		return m, tea.Quit

	case "esc":
		m.back = true
		return m, nil

	case "enter":
		return m.handleEnterKey()

	case "tab", "shift+tab", "down", "up":
		m.cycleFocus(msg.String())
		m.blurFields()

	case "left", "right":
		m.cycleButtonFocus()
	}

	return m.updateTextInputs(msg)
}

// handleEnterKey handles the "enter" key press, triggering actions based on the current focus.
func (m *DuplicateModel) handleEnterKey() (tea.Model, tea.Cmd) {
	if m.focusIndex == DuplicateButtonsFocus {
		if m.buttonFocusIndex == DuplicateSubmitButtonFocus {
			m.startDuplicating()
			m.duplicating = true
		}

		if m.buttonFocusIndex == DuplicateBackButtonFocus {
			m.back = true
		}
	}

	return m, nil
}

// startDuplicating initiates the playlist duplicating process in a separate goroutine.
func (m *DuplicateModel) startDuplicating() {
	config := config.ReadConfig()
	spotifyClient, err := spotify.Login(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to login to Spotify: %v", err))
	}

	playlistId := spotifyApi.ID(spotify.ExtractSpotifyID(m.urlInput.Value()))
	country := m.urlInput.Value()

	playlist, err := spotify.GetPlaylist(spotifyClient, playlistId, country, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to get playlist: %v", err))
	}

	playlistName := m.nameInput.Value()
	if playlistName == "" {
		playlistName = fmt.Sprintf("Copy of %s", playlist.Name)
	}

	description := m.descriptionInput.Value()
	if description == "" {
		description = fmt.Sprintf("Copy of %s", playlist.Name)
	}

	go func() {
		err := duplicator.Run(
			spotifyClient,
			playlistId,
			country,
			duplicator.PlaylistInfo{
				Name:          playlistName,
				Description:   description,
				Public:        m.publicToggle.Value(),
				Collaborative: m.collaborationToggle.Value(),
			},
			m.progressBar,
		)

		if err != nil {
			panic(fmt.Sprintf("Failed to sort playlist: %v", err))
		}

		m.progressBar.SetLabel("Duplication completed!")
		m.exit = true
	}()
}

// blurFields removes focus from all input fields and sets focus to the currently selected field.
func (m *DuplicateModel) blurFields() {
	m.urlInput.Blur()
	m.nameInput.Blur()
	m.descriptionInput.Blur()
	m.publicToggle.Blur()
	m.collaborationToggle.Blur()
	m.countryInput.Blur()

	switch m.focusIndex {
	case DuplicateUrlFieldFocus:
		m.urlInput.Focus()
	case DuplicateNameFieldFocus:
		m.nameInput.Focus()
	case DuplicateDescriptionFieldFocus:
		m.descriptionInput.Focus()
	case DuplicateCollaborationToggleFocus:
		m.collaborationToggle.Focus()
	case DuplicatePublicToggleFocus:
		m.publicToggle.Focus()
	case DuplicateCountryFieldFocus:
		m.countryInput.Focus()
	}
}

// cycleFocus moves the focus between input fields and buttons based on the navigation direction.
func (m *DuplicateModel) cycleFocus(direction string) {
	switch direction {
	case "down", "tab":
		m.focusIndex = (m.focusIndex + 1) % 7
	case "up", "shift+tab":
		m.focusIndex = (m.focusIndex + 6) % 7
	}
}

// cycleButtonFocus switches the focus between the submit and back buttons.
func (m *DuplicateModel) cycleButtonFocus() {
	if m.focusIndex == DuplicateButtonsFocus {
		m.buttonFocusIndex = (m.buttonFocusIndex + 1) % 2
	}
}

// updateTextInputs updates the state of the currently focused text input field.
func (m *DuplicateModel) updateTextInputs(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.focusIndex {
	case DuplicateUrlFieldFocus:
		m.urlInput, cmd = m.urlInput.Update(msg)
	case DuplicateNameFieldFocus:
		m.nameInput, cmd = m.nameInput.Update(msg)
	case DuplicateDescriptionFieldFocus:
		m.descriptionInput, cmd = m.descriptionInput.Update(msg)
	case DuplicatePublicToggleFocus:
		m.publicToggle, cmd = m.publicToggle.Update(msg)
	case DuplicateCollaborationToggleFocus:
		m.collaborationToggle, cmd = m.collaborationToggle.Update(msg)
	case DuplicateCountryFieldFocus:
		m.countryInput, cmd = m.countryInput.Update(msg)
	}
	return m, cmd
}

// View renders the current state of the DuplicateModel as a string for display in the terminal.
func (m *DuplicateModel) View() string {
	if m.duplicating {
		return m.renderDuplicatingView()
	}
	return m.renderFormView()
}

// renderSortingView generates the view for the sorting progress state.
func (m *DuplicateModel) renderDuplicatingView() string {
	return fmt.Sprintf("%s\n%s\n%s\n", APP_TITLE, m.progressBar.Label(), m.progressBar.View())
}

// renderFormView generates the view for the input form and buttons.
func (m *DuplicateModel) renderFormView() string {
	return fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n%s\n%s\n\n%s\nTab to switch fields, Enter to submit, Esc to return.",
		APP_TITLE,
		m.urlInput.View(),
		m.nameInput.View(),
		m.descriptionInput.View(),
		m.publicToggle.View(),
		m.collaborationToggle.View(),
		m.countryInput.View(),
		m.renderButtons(),
	)
}

// renderButtons generates the button view based on the current focus state.
func (m *DuplicateModel) renderButtons() string {
	if m.focusIndex == DuplicateButtonsFocus {
		if m.buttonFocusIndex == DuplicateSubmitButtonFocus {
			return "[ Submit ] - Back -"
		}
		return "- Submit - [ Back ]"
	}
	return "- Submit - - Back -"
}
