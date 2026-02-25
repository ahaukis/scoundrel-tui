package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
	var discardPile *lipgloss.Layer
	if disc := m.Game.LastDiscarded; disc != nil {
		discardPile = faceLayer(disc, false)
	} else {
		discardPile = emptyLayer()
	}

	currentRoom := lipgloss.NewLayer("").X(discardPile.GetX() + discardPile.Width() + 5)

	var r []*game.Card
	if len(m.Game.Room) > 0 {
		r = m.Game.Room
	} else {
		r = make([]*game.Card, 0, game.CardsPerRoom)
		for range cap(r) {
			r = append(r, nil)
		}
	}

	for i, c := range r {
		var cLayer *lipgloss.Layer
		if c != nil {
			cLayer = faceLayer(c, false)
		} else {
			cLayer = emptyLayer()
		}
		currentRoom.AddLayers(cLayer.X((cardWidth + 1) * i))
	}

	var dungeonPile *lipgloss.Layer
	if len(m.Game.Dungeon) > 0 {
		dungeonPile = backLayer(false)
	} else {
		dungeonPile = emptyLayer()
	}
	dungeonPile = dungeonPile.X(currentRoom.GetX() + currentRoom.Width() + 5)

	topRow := lipgloss.NewLayer("", discardPile, currentRoom, dungeonPile)

	comp := lipgloss.NewCompositor(topRow)
	s := comp.Render()
	s += "\n"

	return tea.NewView(s)
}
