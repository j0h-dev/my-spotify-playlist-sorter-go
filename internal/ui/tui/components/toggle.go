package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Toggle struct {
	label string
	value bool

	focus bool
}

func NewToggle(label string, value bool) Toggle {
	return Toggle{
		label: label,
		value: value,
	}
}

func (t Toggle) Init() tea.Cmd {
	return nil
}

func (t Toggle) Update(msg tea.Msg) (Toggle, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			t.Toggle()
		}
	}

	return t, nil
}

func (t *Toggle) Toggle() {
	t.value = !t.value
}

func (t Toggle) View() string {
	s := ""

	if t.value {
		s = fmt.Sprintf("[x] %s", t.label)
	} else {
		s = fmt.Sprintf("[ ] %s", t.label)
	}

	if t.focus {
		return fmt.Sprintf("\033[7m%s\033[0m", s)
	}

	return s
}

func (t *Toggle) Focus() {
	t.focus = true
}

func (t *Toggle) Blur() {
	t.focus = false
}

func (t *Toggle) Value() bool {
	return t.value
}
