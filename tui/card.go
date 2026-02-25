// Lipgloss layers for single cards.

package tui

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/ahaukis/scoundrel-tui/game"
)

const cardHeight = 6
const cardWidth = 8

func borderLayer(selected bool) *lipgloss.Layer {
	var col color.Color
	if selected {
		col = lipgloss.BrightGreen
	} else {
		col = lipgloss.White
	}

	bStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(col).
		Width(cardWidth).
		Height(cardHeight).
		Render()

	return lipgloss.NewLayer(bStyle)
}

func emptyLayer() *lipgloss.Layer {
	bStyle := lipgloss.NewStyle().
		Border(dashedRoundedBorder).
		BorderForeground(lipgloss.BrightBlack).
		Width(cardWidth).
		Height(cardHeight).
		Render()

	return lipgloss.NewLayer(bStyle)
}

func faceLayer(card *game.Card, selected bool) *lipgloss.Layer {
	var col color.Color
	if card.Suit.IsRed() {
		col = lipgloss.BrightRed
	} else {
		col = lipgloss.BrightWhite
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

	return borderLayer(selected).AddLayers(txtLayer1, txtLayer2)
}

func backLayer(selected bool) *lipgloss.Layer {
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

	backStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Black).
		Render(sBuilder.String())

	backLayer := lipgloss.NewLayer(backStyle).X(1).Y(1)

	return borderLayer(selected).AddLayers(backLayer)
}
