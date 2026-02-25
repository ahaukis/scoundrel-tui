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
		discardPile = cardFaceLayer(disc, false)
	} else {
		discardPile = emptySlotLayer()
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
			cLayer = cardFaceLayer(c, false)
		} else {
			cLayer = emptySlotLayer()
		}
		currentRoom.AddLayers(cLayer.X((cardWidth + 1) * i))
	}

	var dungeonPile *lipgloss.Layer
	if len(m.Game.Dungeon) > 0 {
		dungeonPile = cardBackLayer(false)
	} else {
		dungeonPile = emptySlotLayer()
	}
	dungeonPile = dungeonPile.X(currentRoom.GetX() + currentRoom.Width() + 5)

	topRow := lipgloss.NewLayer("", discardPile, currentRoom, dungeonPile)

	playerHand := playerHandLayer(m.Game.Weapon, m.Game.MonstersSlain, false).
		X(topRow.GetX() + (topRow.Width()-cardWidth)/2).
		Y(topRow.GetY() + topRow.Height() + 1)

	comp := lipgloss.NewCompositor(topRow, playerHand)
	s := comp.Render()
	s += "\n"

	return tea.NewView(s)
}
