package tui

import (
	"strconv"

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

func (m *model) unselectRoom() {
	m.selectedRoomIdx = -1
}

func (m *model) isRoomSelected() bool {
	return m.selectedRoomIdx >= 0
}

// Flip the weaponEnabled bool to the other option
func (m *model) toggleWeapon() {
	m.weaponEnabled = !m.weaponEnabled
}

// Handle a key press that doesn't need to return its own message.
func (m *model) handleKeyPress(msg tea.KeyPressMsg) {
	s := msg.String()

	// handle 1,2,...
	if i, err := strconv.Atoi(s); err == nil {
		if 1 <= i && i <= game.CardsPerRoom {
			m.selectedRoomIdx = i - 1
			m.game.MakeRoomAction(m.selectedRoomIdx, m.weaponEnabled)
			return
		}
	}

	switch s {
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
		} else if m.isRoomSelected() {
			m.selectedRoomIdx++
		}

	case "down", "j":
		if m.isRoomSelected() {
			m.unselectRoom()
			m.selectedHand = true
		}

	case "up", "k":
		if m.selectedHand {
			m.selectedHand = false
			m.selectedRoomIdx = 0
		}

	case "space", "enter":
		if m.isRoomSelected() {
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.BackgroundColorMsg:
		// set background color based on
		m.hasDarkBackground = msg.IsDark()
		return m, nil

	case tea.KeyPressMsg:
		if s := msg.String(); s == "ctrl+c" || s == "q" {
			return m, tea.Quit
		}
		m.handleKeyPress(msg)

	}

	if m.game.IsRoomDone() {
		m.game.DealRoom()
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

	playerHand := playerHandLayer(m.game.Weapon, m.game.MonstersSlain, m.selectedHand, m.hasDarkBackground).
		X(topRow.GetX() + (topRow.Width()-cardWidth)/2).
		Y(topRow.GetY() + topRow.Height() + 1)

	footer := footerLayer(m.game.HP, m.weaponEnabled).Y(playerHand.GetY() + playerHand.Height())

	comp := lipgloss.NewCompositor(topRow, playerHand, footer)
	s := comp.Render()
	s += "\n"

	return tea.NewView(s)
}
