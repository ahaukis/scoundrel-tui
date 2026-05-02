package table

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/ahaukis/scoundrel-tui/internal/game"
	keymap "github.com/ahaukis/scoundrel-tui/internal/tui/keymap"
	"github.com/ahaukis/scoundrel-tui/internal/tui/palette"
)

const cardHeight = 6
const cardWidth = 8

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

type Model struct {
	game            *game.Game
	selectedRoomIdx int // -1 if no card in room is currently selected
	selectedDungeon bool
	selectedHand    bool
	weaponEnabled   bool
	palette         *palette.Palette
}

func New(g *game.Game, p *palette.Palette) Model {
	return Model{
		game:            g,
		selectedRoomIdx: 0,
		weaponEnabled:   true,
		palette:         p,
	}
}

func (m Model) Init() tea.Cmd {
	m.game.DealRoom()
	return nil
}

// Unselect the selected room card, if any.
func (m *Model) unselectRoom() {
	m.selectedRoomIdx = -1
}

// Check if the any of the room cards is currently selected.
func (m *Model) selectionInRoom() bool {
	return m.selectedRoomIdx >= 0
}

// Flip the weaponEnabled bool to the other option
func (m *Model) toggleWeapon() {
	m.weaponEnabled = !m.weaponEnabled
}

// Handle a key press that doesn't need to return its own message.
func (m *Model) handleKeyPress(msg tea.KeyPressMsg, keys keymap.KeyMap) {
	switch {
	case key.Matches(msg, keys.SelectInRoom):
		i, err := strconv.Atoi(msg.String())
		if err != nil {
			panic(fmt.Errorf("invalid index: %w", err))
		}
		if i < 1 || i > game.CardsPerRoom {
			panic(fmt.Errorf("card index out of range: %d", i))
		}
		m.selectedRoomIdx = i - 1
		m.game.MakeRoomAction(m.selectedRoomIdx, m.weaponEnabled)

	case key.Matches(msg, keys.Left):
		if m.selectedDungeon {
			m.selectedDungeon = false
			m.selectedRoomIdx = game.CardsPerRoom - 1
		} else if m.selectedRoomIdx > 0 {
			m.selectedRoomIdx--
		}
	case key.Matches(msg, keys.Right):
		if m.selectedRoomIdx == game.CardsPerRoom-1 {
			m.unselectRoom()
			m.selectedDungeon = true
		} else if m.selectionInRoom() {
			m.selectedRoomIdx++
		}
	case key.Matches(msg, keys.Down):
		if m.selectionInRoom() || m.selectedDungeon {
			m.unselectRoom()
			m.selectedDungeon = false
			m.selectedHand = true
		}
	case key.Matches(msg, keys.Up):
		if m.selectedHand {
			m.selectedHand = false
			m.selectedRoomIdx = 0
		}

	case key.Matches(msg, keys.Interact):
		if m.selectionInRoom() {
			m.game.MakeRoomAction(m.selectedRoomIdx, m.weaponEnabled)
		} else if m.selectedDungeon {
			m.game.SkipRoom()
		} else if m.selectedHand {
			m.toggleWeapon()
		}
	case key.Matches(msg, keys.ToggleWeapon):
		m.toggleWeapon()
	case key.Matches(msg, keys.SkipRoom):
		m.game.SkipRoom()
	}
}

func (m Model) Update(msg tea.Msg, keys keymap.KeyMap) (Model, tea.Cmd) {
	if keyPressMsg, ok := msg.(tea.KeyPressMsg); ok {
		m.handleKeyPress(keyPressMsg, keys)
	}

	if m.game.IsRoomDone() {
		m.game.DealRoom()
	}

	return m, nil
}

func (m Model) View() string {
	var discardPile *lipgloss.Layer
	if disc := m.game.LastDiscarded; disc != nil {
		discardPile = m.cardFaceLayer(disc, false, true)
	} else {
		discardPile = m.emptySlotLayer(false)
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
			cLayer = m.cardFaceLayer(c, selected, true)
		} else {
			cLayer = m.emptySlotLayer(selected)
		}
		currentRoom.AddLayers(cLayer.X((cardWidth + 1) * i))
	}

	dungeonPile := m.dungeonPileLayer().X(currentRoom.GetX() + currentRoom.Width() + 5)

	topRow := lipgloss.NewLayer("", discardPile, currentRoom, dungeonPile)

	playerHand := m.playerHandLayer(m.game.Weapon, m.game.MonstersSlain, m.selectedHand, m.weaponEnabled).
		X(topRow.GetX() + (topRow.Width()-cardWidth)/2).
		Y(currentRoom.GetY() + currentRoom.Height() + 1)

	comp := lipgloss.NewCompositor(topRow, playerHand)
	s := comp.Render()
	s += "\n"

	return s
}

func (m *Model) cardBorderStyle(selected, active bool) lipgloss.Style {
	var col color.Color
	switch {
	case selected && active:
		col = m.palette.SelectedBorder
	case selected && !active:
		col = m.palette.SelectedEmptyBorder
	case !selected && active:
		col = m.palette.Border
	default:
		col = m.palette.EmptyBorder
	}

	bStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(col).
		Width(cardWidth).
		Height(cardHeight)

	return bStyle
}

func (m *Model) emptySlotLayer(selected bool) *lipgloss.Layer {
	var col color.Color
	if selected {
		col = m.palette.SelectedEmptyBorder
	} else {
		col = m.palette.EmptyBorder
	}

	bStyle := lipgloss.NewStyle().
		Border(dashedRoundedBorder).
		BorderForeground(col).
		Width(cardWidth).
		Height(cardHeight).
		Render()

	return lipgloss.NewLayer(bStyle)
}

func (m *Model) cardFaceLayer(card *game.Card, selected, active bool) *lipgloss.Layer {
	var col color.Color
	if card.Suit.IsRed() {
		col = m.palette.RedSuit
	} else {
		col = m.palette.BlackSuit
	}

	s := card.String()

	txt1 := lipgloss.NewStyle().Foreground(col).Render(s)
	txt2 := txt1

	txtLayer1 := lipgloss.NewLayer(txt1).X(1).Y(1)
	txtLayer2 := lipgloss.NewLayer(txt2).
		X(cardWidth - lipgloss.Width(s) - 1).
		Y(cardHeight - 2)

	bLayer := m.cardBorderStyle(selected, active).Render()

	return lipgloss.NewLayer(bLayer, txtLayer1, txtLayer2)
}

func (m *Model) cardBackLayer(selected, active bool) *lipgloss.Layer {
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

	backStyle := m.cardBorderStyle(selected, active).
		Foreground(m.palette.CardBack).
		Render(sBuilder.String())

	cardBackLayer := lipgloss.NewLayer(backStyle)

	return cardBackLayer
}

func (m *Model) playerHandLayer(weapon *game.Card, slain []*game.Card, selected, active bool) *lipgloss.Layer {
	if weapon == nil {
		return m.emptySlotLayer(selected)
	}

	playerHand := lipgloss.NewLayer("", m.cardFaceLayer(weapon, selected, active))

	for i, s := range slain {
		sLayer := m.cardFaceLayer(s, selected, active).Y(2 * (i + 1))
		playerHand.AddLayers(sLayer)
	}

	return playerHand
}

func (m *Model) dungeonPileLayer() *lipgloss.Layer {
	lenDungeon := len(m.game.Dungeon)

	var dungeonPile *lipgloss.Layer
	if lenDungeon > 0 {
		dungeonPile = m.cardBackLayer(m.selectedDungeon, m.game.CanSkipRoom())
	} else {
		dungeonPile = m.emptySlotLayer(m.selectedDungeon)
	}

	txt := lipgloss.NewStyle().
		Foreground(m.palette.Border).
		AlignHorizontal(lipgloss.Right).
		Width(dungeonPile.Width()).
		Render(fmt.Sprintf("%2d left", lenDungeon))

	dungeonPile = lipgloss.NewLayer(dungeonPile.GetContent() + "\n" + txt)

	return dungeonPile
}
