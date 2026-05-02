// Color constants.

package palette

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

type Palette struct {
	Border              color.Color
	SelectedBorder      color.Color
	EmptyBorder         color.Color
	SelectedEmptyBorder color.Color
	RedSuit             color.Color
	BlackSuit           color.Color
	CardBack            color.Color
	HPFull              color.Color
	HPMid               color.Color
	HPLow               color.Color
}

func NewLight() Palette {
	return Palette{
		Border:              lipgloss.BrightBlack,
		SelectedBorder:      lipgloss.BrightGreen,
		EmptyBorder:         lipgloss.BrightBlack,
		SelectedEmptyBorder: lipgloss.BrightGreen,
		RedSuit:             lipgloss.Red,
		BlackSuit:           lipgloss.Black,
		CardBack:            lipgloss.BrightBlack,
		HPFull:              lipgloss.Green,
		HPMid:               lipgloss.Yellow,
		HPLow:               lipgloss.Red,
	}
}

func NewDark() Palette {
	return Palette{
		Border:              lipgloss.White,
		SelectedBorder:      lipgloss.BrightGreen,
		EmptyBorder:         lipgloss.BrightBlack,
		SelectedEmptyBorder: lipgloss.Green,
		RedSuit:             lipgloss.BrightRed,
		BlackSuit:           lipgloss.BrightWhite,
		CardBack:            lipgloss.White,
		HPFull:              lipgloss.BrightGreen,
		HPMid:               lipgloss.BrightYellow,
		HPLow:               lipgloss.BrightRed,
	}
}
