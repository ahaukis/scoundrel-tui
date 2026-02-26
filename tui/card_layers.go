// Lipgloss layers for single cards.

package tui

import (
	"image/color"
	"strings"

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

func lightDark(hasDarkBackground bool, colors [2]color.Color) color.Color {
	return lipgloss.LightDark(hasDarkBackground)(colors[0], colors[1])
}

func cardBorderLayer(selected, hasDarkBackground bool) *lipgloss.Layer {
	var col color.Color
	if selected {
		col = lightDark(hasDarkBackground, colorScheme["selectedBorder"])
	} else {
		col = lightDark(hasDarkBackground, colorScheme["border"])
	}

	bStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(col).
		Width(cardWidth).
		Height(cardHeight).
		Render()

	return lipgloss.NewLayer(bStyle)
}

func emptySlotLayer(selected, hasDarkBackground bool) *lipgloss.Layer {
	var col color.Color
	if selected {
		col = lightDark(hasDarkBackground, colorScheme["selectedEmptyBorder"])
	} else {
		col = lightDark(hasDarkBackground, colorScheme["emptyBorder"])
	}

	bStyle := lipgloss.NewStyle().
		Border(dashedRoundedBorder).
		BorderForeground(col).
		Width(cardWidth).
		Height(cardHeight).
		Render()

	return lipgloss.NewLayer(bStyle)
}

func cardFaceLayer(card *game.Card, selected, hasDarkBackground bool) *lipgloss.Layer {
	var col color.Color
	if card.Suit.IsRed() {
		col = lightDark(hasDarkBackground, colorScheme["redSuit"])
	} else {
		col = lightDark(hasDarkBackground, colorScheme["blackSuit"])
	}

	s := card.String()

	txt1 := lipgloss.NewStyle().Foreground(col).Render(s)
	txt2 := txt1

	txtLayer1 := lipgloss.NewLayer(txt1).X(1).Y(1)
	txtLayer2 := lipgloss.NewLayer(txt2).
		X(cardWidth - lipgloss.Width(s) - 1).
		Y(cardHeight - 2)

	return cardBorderLayer(selected, hasDarkBackground).AddLayers(txtLayer1, txtLayer2)
}

func cardBackLayer(selected, hasDarkBackground bool) *lipgloss.Layer {
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
		Foreground(lightDark(hasDarkBackground, colorScheme["cardBack"])).
		Render(sBuilder.String())

	cardBackLayer := lipgloss.NewLayer(backStyle).X(1).Y(1)

	return cardBorderLayer(selected, hasDarkBackground).AddLayers(cardBackLayer)
}

func playerHandLayer(weapon *game.Card, slain []*game.Card, selected, hasDarkBackground bool) *lipgloss.Layer {
	if weapon == nil {
		return emptySlotLayer(selected, hasDarkBackground)
	}

	playerHand := lipgloss.NewLayer("", cardFaceLayer(weapon, selected, hasDarkBackground))

	for i, s := range slain {
		sLayer := cardFaceLayer(s, selected, hasDarkBackground).Y(2 * (i + 1))
		playerHand.AddLayers(sLayer)
	}

	return playerHand
}
