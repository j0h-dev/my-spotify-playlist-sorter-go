package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/ui/tui"
)

func main() {
	p := tea.NewProgram(tui.NewApp())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
