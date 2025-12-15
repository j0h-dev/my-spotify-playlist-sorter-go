package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	MaxProgressBarWidth = 80
)

// newTextInput creates and configures a new text input field with the given parameters.
func newTextInput(prompt, placeholder string, width int, focused bool) textinput.Model {
	input := textinput.New()
	input.Placeholder = placeholder
	input.Prompt = prompt
	input.Width = width
	if focused {
		input.Focus()
	}
	return input
}

type tickMsg time.Time

// tickCmd creates a command that sends a tick message every second.
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
