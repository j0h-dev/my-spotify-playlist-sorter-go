package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/config"
	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/sorter"
	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/spotify"
	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/ui/tui/components"

	spotifyApi "github.com/zmb3/spotify/v2"
)

const (
	SortUrlFieldFocus = iota
	SortCountryFieldFocus
	SortButtonsFocus
)

const (
	SortSubmitButtonFocus = iota
	SortBackButtonFocus
)

type SortModel struct {
	cursor int
	items  []string

	urlInput         textinput.Model
	countryInput     textinput.Model
	focusIndex       int
	buttonFocusIndex int

	back bool

	progressBar *components.ProgressBar
	sorting     bool
	exit        bool
}

// NewSortModel initializes and returns a new SortModel with default values.
func NewSortModel() *SortModel {
	return &SortModel{
		urlInput:     newTextInput("URL: ", "https://open.spotify.com/playlist/...", 80, true),
		countryInput: newTextInput("Country: ", "FI, SV, US...", 20, false),
		focusIndex:   SortUrlFieldFocus,
		progressBar:  components.NewProgressBar(),
	}
}

// Init initializes the SortModel and returns the initial command to start ticking.
func (m *SortModel) Init() tea.Cmd {
	return tickCmd()
}

// Update handles incoming messages and updates the state of the SortModel accordingly.
// It returns the updated model and any commands to be executed.
func (m *SortModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
func (m *SortModel) updateProgressBarWidth(windowWidth int) {
	m.progressBar.SetWidth(windowWidth - 2*2 - 4)
	if m.progressBar.Width() > MaxProgressBarWidth {
		m.progressBar.SetWidth(MaxProgressBarWidth)
	}
}

// handleKeyMsg processes key press events and updates the model state accordingly.
func (m *SortModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
func (m *SortModel) handleEnterKey() (tea.Model, tea.Cmd) {
	if m.focusIndex == SortButtonsFocus {
		if m.buttonFocusIndex == SortSubmitButtonFocus {
			m.startSorting()
			m.sorting = true
		}

		if m.buttonFocusIndex == SortBackButtonFocus {
			m.back = true
		}
	}

	return m, nil
}

// startSorting initiates the playlist sorting process in a separate goroutine.
func (m *SortModel) startSorting() {
	config := config.ReadConfig()
	spotifyClient, err := spotify.Login(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to login to Spotify: %v", err))
	}

	playlistId := spotifyApi.ID(spotify.ExtractSpotifyID(m.urlInput.Value()))
	country := m.urlInput.Value()

	go func() {
		err := sorter.Run(spotifyClient, playlistId, country, m.progressBar)
		if err != nil {
			panic(fmt.Sprintf("Failed to sort playlist: %v", err))
		}
		m.progressBar.SetLabel("Sorting completed!")
		m.exit = true
	}()
}

// blurFields removes focus from all input fields and sets focus to the currently selected field.
func (m *SortModel) blurFields() {
	m.urlInput.Blur()
	m.countryInput.Blur()

	switch m.focusIndex {
	case SortUrlFieldFocus:
		m.urlInput.Focus()
	case SortCountryFieldFocus:
		m.countryInput.Focus()
	}
}

// cycleFocus moves the focus between input fields and buttons based on the navigation direction.
func (m *SortModel) cycleFocus(direction string) {
	switch direction {
	case "down", "tab":
		m.focusIndex = (m.focusIndex + 1) % 3
	case "up", "shift+tab":
		m.focusIndex = (m.focusIndex + 2) % 3
	}
}

// cycleButtonFocus switches the focus between the submit and back buttons.
func (m *SortModel) cycleButtonFocus() {
	if m.focusIndex == SortButtonsFocus {
		m.buttonFocusIndex = (m.buttonFocusIndex + 1) % 2
	}
}

// updateTextInputs updates the state of the currently focused text input field.
func (m *SortModel) updateTextInputs(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.focusIndex {
	case SortUrlFieldFocus:
		m.urlInput, cmd = m.urlInput.Update(msg)
	case SortCountryFieldFocus:
		m.countryInput, cmd = m.countryInput.Update(msg)
	}
	return m, cmd
}

// View renders the current state of the SortModel as a string for display in the terminal.
func (m *SortModel) View() string {
	if m.sorting {
		return m.renderSortingView()
	}
	return m.renderFormView()
}

// renderSortingView generates the view for the sorting progress state.
func (m *SortModel) renderSortingView() string {
	return fmt.Sprintf("%s\n%s\n%s\n", APP_TITLE, m.progressBar.Label(), m.progressBar.View())
}

// renderFormView generates the view for the input form and buttons.
func (m *SortModel) renderFormView() string {
	buttons := m.renderButtons()
	return fmt.Sprintf("%s\n%s\n%s\n\n%s\nTab to switch fields, Enter to submit, Esc to return.",
		APP_TITLE, m.urlInput.View(), m.countryInput.View(), buttons)
}

// renderButtons generates the button view based on the current focus state.
func (m *SortModel) renderButtons() string {
	if m.focusIndex == SortButtonsFocus {
		if m.buttonFocusIndex == SortSubmitButtonFocus {
			return "[ Submit ] - Back -"
		}
		return "- Submit - [ Back ]"
	}
	return "- Submit - - Back -"
}
