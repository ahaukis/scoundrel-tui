package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/ahaukis/scoundrel-tui/game"
)

type model struct {
	game              *game.Game
	selectedRoomIdx   int // -1 if no card in room is currently selected
	selectedDungeon   bool
	selectedHand      bool
	weaponEnabled     bool
	hasDarkBackground bool
}

func InitialModel() model {
	g := game.NewRandomGame()
	return model{game: g, selectedRoomIdx: 0, weaponEnabled: true}
}

func (m model) Init() tea.Cmd {
	m.game.DealRoom()
	return tea.RequestBackgroundColor
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.BackgroundColorMsg:
		m.hasDarkBackground = msg.IsDark()
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "left":
			if m.selectedDungeon {
				m.selectedDungeon = false
				m.selectedRoomIdx = game.CardsPerRoom - 1
			} else if m.selectedRoomIdx > 0 {
				m.selectedRoomIdx--
			}

		case "right":
			if m.selectedRoomIdx == game.CardsPerRoom-1 {
				m.selectedRoomIdx = -1
				m.selectedDungeon = true
			} else if m.selectedRoomIdx >= 0 {
				m.selectedRoomIdx++
			}

		case "space", "enter":
			if i := m.selectedRoomIdx; i >= 0 {
				m.game.MakeRoomAction(i, m.weaponEnabled)
			}

		case "s":
			if !m.game.SkippedLastRoom {
				m.game.SkipRoom()
			}

		}
	}

	if m.game.IsRoomDone() {
		m.game.DealRoom()
		return m, nil
	}

	return m, nil
}

func (m model) View() tea.View {
	var discardPile *lipgloss.Layer
	if disc := m.game.LastDiscarded; disc != nil {
		discardPile = cardFaceLayer(disc, false, m.hasDarkBackground)
	} else {
		discardPile = emptySlotLayer(false, m.hasDarkBackground)
	}

	currentRoom := lipgloss.NewLayer("").X(discardPile.GetX() + discardPile.Width() + 5)

	var r []*game.Card
	if len(m.game.Room) > 0 {
		r = m.game.Room
	} else {
		r = make([]*game.Card, 0, game.CardsPerRoom)
		for range cap(r) {
			r = append(r, nil)
		}
	}

	for i, c := range r {
		var cLayer *lipgloss.Layer
		selected := m.selectedRoomIdx == i
		if c != nil {
			cLayer = cardFaceLayer(c, selected, m.hasDarkBackground)
		} else {
			cLayer = emptySlotLayer(selected, m.hasDarkBackground)
		}
		currentRoom.AddLayers(cLayer.X((cardWidth + 1) * i))
	}

	var dungeonPile *lipgloss.Layer
	if len(m.game.Dungeon) > 0 {
		dungeonPile = cardBackLayer(m.selectedDungeon, m.hasDarkBackground)
	} else {
		dungeonPile = emptySlotLayer(m.selectedDungeon, m.hasDarkBackground)
	}
	dungeonPile = dungeonPile.X(currentRoom.GetX() + currentRoom.Width() + 5)

	topRow := lipgloss.NewLayer("", discardPile, currentRoom, dungeonPile)

	playerHand := playerHandLayer(m.game.Weapon, m.game.MonstersSlain, false, m.hasDarkBackground).
		X(topRow.GetX() + (topRow.Width()-cardWidth)/2).
		Y(topRow.GetY() + topRow.Height() + 1)

	footer := footerLayer(m.game.HP, m.weaponEnabled).Y(playerHand.GetY() + playerHand.Height())

	comp := lipgloss.NewCompositor(topRow, playerHand, footer)
	s := comp.Render()
	s += "\n"

	return tea.NewView(s)
}
