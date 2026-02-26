package tui

import (
	"image/color"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/ahaukis/scoundrel-tui/game"
)

// A custom dashed border with rounded corners.
var dashedRoundedBorder = lipgloss.Border{
	Top:          "╌",
	Bottom:       "╌",
	Left:         "┆",
	Right:        "┆",
	TopLeft:      "╭",
	TopRight:     "╮",
	BottomLeft:   "╰",
	BottomRight:  "╯",
	MiddleLeft:   "├",
	MiddleRight:  "┤",
	Middle:       "┼",
	MiddleTop:    "┬",
	MiddleBottom: "┴",
}

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
		discardPile = m.cardFaceLayer(disc, false)
	} else {
		discardPile = m.emptySlotLayer(false)
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
			cLayer = m.cardFaceLayer(c, selected)
		} else {
			cLayer = m.emptySlotLayer(selected)
		}
		currentRoom.AddLayers(cLayer.X((cardWidth + 1) * i))
	}

	var dungeonPile *lipgloss.Layer
	if len(m.game.Dungeon) > 0 {
		dungeonPile = m.cardBackLayer(m.selectedDungeon)
	} else {
		dungeonPile = m.emptySlotLayer(m.selectedDungeon)
	}
	dungeonPile = dungeonPile.X(currentRoom.GetX() + currentRoom.Width() + horizontalSpace)

	topRow := lipgloss.NewLayer("", discardPile, currentRoom, dungeonPile)

	playerHand := m.playerHandLayer(m.game.Weapon, m.game.MonstersSlain, m.selectedHand).
		X(topRow.GetX() + (topRow.Width()-cardWidth)/2).
		Y(topRow.GetY() + topRow.Height() + verticalSpace)

	footer := footerLayer(m.game.HP, m.weaponEnabled).Y(playerHand.GetY() + playerHand.Height())

	comp := lipgloss.NewCompositor(topRow, playerHand, footer)
	s := comp.Render()
	s += "\n"

	return tea.NewView(s)
}

func (m *tableModel) lightDark(colors [2]color.Color) color.Color {
	return lipgloss.LightDark(m.hasDarkBackground)(colors[0], colors[1])
}

func (m *tableModel) cardBorderStyle(selected bool) lipgloss.Style {
	var col color.Color
	if selected {
		col = m.lightDark(colorScheme["selectedBorder"])
	} else {
		col = m.lightDark(colorScheme["border"])
	}
	bStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(col).
		Width(cardWidth).
		Height(cardHeight)

	return bStyle
}

func (m *tableModel) emptySlotLayer(selected bool) *lipgloss.Layer {
	var col color.Color
	if selected {
		col = m.lightDark(colorScheme["selectedEmptyBorder"])
	} else {
		col = m.lightDark(colorScheme["emptyBorder"])
	}

	bStyle := lipgloss.NewStyle().
		Border(dashedRoundedBorder).
		BorderForeground(col).
		Width(cardWidth).
		Height(cardHeight).
		Render()

	return lipgloss.NewLayer(bStyle)
}

func (m *tableModel) cardFaceLayer(card *game.Card, selected bool) *lipgloss.Layer {
	var col color.Color
	if card.Suit.IsRed() {
		col = m.lightDark(colorScheme["redSuit"])
	} else {
		col = m.lightDark(colorScheme["blackSuit"])
	}

	s := card.String()

	txt1 := lipgloss.NewStyle().Foreground(col).Render(s)
	txt2 := txt1

	txtLayer1 := lipgloss.NewLayer(txt1).X(1).Y(1)
	txtLayer2 := lipgloss.NewLayer(txt2).
		X(cardWidth - lipgloss.Width(s) - 1).
		Y(cardHeight - 2)

	bLayer := m.cardBorderStyle(selected).Render()

	return lipgloss.NewLayer(bLayer, txtLayer1, txtLayer2)
}

func (m *tableModel) cardBackLayer(selected bool) *lipgloss.Layer {
	sBuilder := strings.Builder{}
	rows := cardHeight - 2
	columns := cardWidth - 2

	for i := range rows {
		for range columns {
			sBuilder.WriteString("#")
		}
		if i < rows-1 {
			sBuilder.WriteString("\n")
		}
	}

	backStyle := m.cardBorderStyle(selected).
		Foreground(m.lightDark(colorScheme["cardBack"])).
		Render(sBuilder.String())

	cardBackLayer := lipgloss.NewLayer(backStyle)

	return cardBackLayer
}

func (m *tableModel) playerHandLayer(weapon *game.Card, slain []*game.Card, selected bool) *lipgloss.Layer {
	if weapon == nil {
		return m.emptySlotLayer(selected)
	}

	playerHand := lipgloss.NewLayer("", m.cardFaceLayer(weapon, selected))

	for i, s := range slain {
		sLayer := m.cardFaceLayer(s, selected).Y(2 * (i + 1))
		playerHand.AddLayers(sLayer)
	}

	return playerHand
}
