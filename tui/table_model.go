package tui

import (
	"strconv"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/ahaukis/scoundrel-tui/game"
)

type tableModel struct {
	game              *game.Game
	selectedRoomIdx   int // -1 if no card in room is currently selected
	selectedDungeon   bool
	selectedHand      bool
	weaponEnabled     bool
	hasDarkBackground bool
}

func (m tableModel) Init() tea.Cmd {
	m.game.DealRoom()
	return nil
}

// Unselect the selected room card, if any.
func (m *tableModel) unselectRoom() {
	m.selectedRoomIdx = -1
}

// Check if the any of the room cards is currently selected.
func (m *tableModel) selectionInRoom() bool {
	return m.selectedRoomIdx >= 0
}

// Flip the weaponEnabled bool to the other option
func (m *tableModel) toggleWeapon() {
	m.weaponEnabled = !m.weaponEnabled
}

// Handle a key press that doesn't need to return its own message.
func (m *tableModel) handleKeyPress(msg tea.KeyPressMsg) {
	msgString := msg.String()

	// handle 1,2,...
	if i, err := strconv.Atoi(msgString); err == nil {
		if 1 <= i && i <= game.CardsPerRoom {
			m.selectedRoomIdx = i - 1
			m.game.MakeRoomAction(m.selectedRoomIdx, m.weaponEnabled)
			return
		}
	}

	switch msgString {
	case "left", "h":
		if m.selectedDungeon {
			m.selectedDungeon = false
			m.selectedRoomIdx = game.CardsPerRoom - 1
		} else if m.selectedRoomIdx > 0 {
			m.selectedRoomIdx--
		}

	case "right", "l":
		if m.selectedRoomIdx == game.CardsPerRoom-1 {
			m.unselectRoom()
			m.selectedDungeon = true
		} else if m.selectionInRoom() {
			m.selectedRoomIdx++
		}

	case "down", "j":
		if m.selectionInRoom() {
			m.unselectRoom()
			m.selectedHand = true
		}

	case "up", "k":
		if m.selectedHand {
			m.selectedHand = false
			m.selectedRoomIdx = 0
		}

	case "space", "enter":
		if m.selectionInRoom() {
			m.game.MakeRoomAction(m.selectedRoomIdx, m.weaponEnabled)
		} else if m.selectedDungeon {
			m.game.SkipRoom()
		} else if m.selectedHand {
			m.toggleWeapon()
		}

	case "w":
		m.toggleWeapon()

	case "s":
		m.game.SkipRoom()

	}
}

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyPressMsg, ok := msg.(tea.KeyPressMsg); ok {
		m.handleKeyPress(keyPressMsg)
	}

	if m.game.IsRoomDone() {
		m.game.DealRoom()
	}

	return m, nil
}

func (m tableModel) View() tea.View {
	var discardPile *lipgloss.Layer
	if disc := m.game.LastDiscarded; disc != nil {
		discardPile = cardFaceLayer(disc, false, m.hasDarkBackground)
	} else {
		discardPile = emptySlotLayer(false, m.hasDarkBackground)
	}

	currentRoom := lipgloss.NewLayer("").X(discardPile.GetX() + discardPile.Width() + horizontalSpace)

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
	dungeonPile = dungeonPile.X(currentRoom.GetX() + currentRoom.Width() + horizontalSpace)

	topRow := lipgloss.NewLayer("", discardPile, currentRoom, dungeonPile)

	playerHand := playerHandLayer(m.game.Weapon, m.game.MonstersSlain, m.selectedHand, m.hasDarkBackground).
		X(topRow.GetX() + (topRow.Width()-cardWidth)/2).
		Y(topRow.GetY() + topRow.Height() + verticalSpace)

	footer := footerLayer(m.game.HP, m.weaponEnabled).Y(playerHand.GetY() + playerHand.Height())

	comp := lipgloss.NewCompositor(topRow, playerHand, footer)
	s := comp.Render()
	s += "\n"

	return tea.NewView(s)
}
