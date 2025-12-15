package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/j0h-dev/my-spotify-playlist-sorter-go/internal/config"
)

const APP_TITLE = `
       _____ _____ _____
      |   __|  _  |   __|
      |__   |   __|__   |
      |_____|__|  |_____|

-- My Spotify Playlist Sorter --
`

type Screen int

const (
	ScreenLogin Screen = iota
	ScreenMenu
	ScreenSort
	ScreenDuplicate
)

type App struct {
	current Screen

	initializedScreens [4]bool
	menu               *MenuModel
	login              *LoginModel
	sort               *SortModel
	duplicate          *DuplicateModel

	width  int
	height int

	initScreens bool
}

func NewApp() App {
	firstScreen := ScreenMenu

	conf := config.ReadConfig()
	if conf.Credentials == nil {
		firstScreen = ScreenLogin
	}

	return App{
		current:   firstScreen,
		menu:      NewMenuModel(),
		login:     NewLoginModel(),
		sort:      NewSortModel(),
		duplicate: NewDuplicateModel(),
	}
}

func (a App) Init() tea.Cmd {
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		a.width, a.height = msg.Width, msg.Height
	}

	switch a.current {

	case ScreenMenu:
		m, cmd := a.menu.Update(msg)
		a.menu = m.(*MenuModel)

		switch a.menu.chosen {
		case "duplicate":
			a.current = ScreenDuplicate
		case "sort":
			a.current = ScreenSort
		}
		a.menu.chosen = ""

		if !a.initializedScreens[ScreenMenu] {
			cmd = tea.Batch(cmd, a.menu.Init())
			a.initializedScreens[ScreenMenu] = true
		}

		return a, cmd

	case ScreenLogin:
		if a.login == nil {
			a.login = NewLoginModel()
			return a, a.login.Init()
		}

		m, cmd := a.login.Update(msg)
		a.login = m.(*LoginModel)

		if a.login.toMenu {
			a.current = ScreenMenu
			a.login.toMenu = false
		}

		if !a.initializedScreens[ScreenLogin] {
			cmd = a.login.Init()
			a.initializedScreens[ScreenLogin] = true
		}

		return a, cmd

	case ScreenSort:
		m, cmd := a.sort.Update(msg)
		a.sort = m.(*SortModel)

		if a.sort.back {
			a.current = ScreenMenu
			a.sort.back = false
		}

		if !a.initializedScreens[ScreenSort] {
			cmd = tea.Batch(cmd, a.sort.Init())
			a.initializedScreens[ScreenSort] = true
		}

		return a, cmd

	case ScreenDuplicate:
		m, cmd := a.duplicate.Update(msg)
		a.duplicate = m.(*DuplicateModel)

		if a.duplicate.back {
			a.current = ScreenMenu
			a.duplicate.back = false
		}

		if !a.initializedScreens[ScreenDuplicate] {
			cmd = tea.Batch(cmd, a.duplicate.Init())
			a.initializedScreens[ScreenDuplicate] = true
		}

		return a, cmd
	}
	return a, nil
}

func (a App) View() string {
	switch a.current {
	case ScreenMenu:
		return a.menu.View()
	case ScreenLogin:
		return a.login.View()
	case ScreenSort:
		return a.sort.View()
	case ScreenDuplicate:
		return a.duplicate.View()
	default:
		return ""
	}
}
