package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type MenuModel struct {
	cursor int
	items  []string

	chosen string // "sort", "duplicate", etc.
}

func NewMenuModel() *MenuModel {
	return &MenuModel{
		cursor: 0,
		items:  []string{"Sort playlist", "Duplicate playlist", "Quit"},
	}
}

func (m *MenuModel) Init() tea.Cmd { return nil }

func (m *MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c", "esc":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}

		case "enter":
			switch m.cursor {
			case 0:
				m.chosen = "sort"
			case 1:
				m.chosen = "duplicate"
			case 2:
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m *MenuModel) View() string {
	s := fmt.Sprintf("%s\n", APP_TITLE)

	for i, item := range m.items {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, item)
	}

	return s
}
