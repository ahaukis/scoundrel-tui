package tui

import (
	"image/color"

	"charm.land/lipgloss/v2"
	"github.com/ahaukis/scoundrel-tui/game"
)

const cardHeight = 6
const cardWidth = 8

func border(selected bool) string {
	var col color.Color
	if selected {
		col = lipgloss.Green
	} else {
		col = lipgloss.White
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Foreground(col).
		Width(cardWidth).
		Height(cardHeight).
		Render()
}

func cardLayer(card *game.Card, selected bool) *lipgloss.Layer {
	var col color.Color
	if card.Suit.IsRed() {
		col = lipgloss.Red
	} else {
		col = lipgloss.White
	}

	s := card.String()

	txt1 := lipgloss.NewStyle().
		Foreground(col).
		Render(s)

	txt2 := lipgloss.NewStyle().
		Foreground(col).
		Render(s)

	txtLayer1 := lipgloss.NewLayer(txt1).X(1).Y(1)
	txtLayer2 := lipgloss.NewLayer(txt2).
		X(cardWidth - lipgloss.Width(s) - 1).
		Y(cardHeight - 2)

	return lipgloss.NewLayer(border(selected), txtLayer1, txtLayer2)
}
