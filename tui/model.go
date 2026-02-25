package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/ahaukis/scoundrel-tui/game"
)

type model struct {
	Game *game.Game
}

func InitialModel() model {
	g := game.NewRandomGame()
	return model{g}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	return tea.NewView("Hello world\n")
}
