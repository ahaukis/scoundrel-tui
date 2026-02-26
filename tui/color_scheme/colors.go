// Color constants.

package color_scheme

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

var ColorScheme = map[string][2]color.Color{
	"border":              {lipgloss.BrightBlack, lipgloss.White},
	"selectedBorder":      {lipgloss.BrightGreen, lipgloss.BrightGreen},
	"emptyBorder":         {lipgloss.BrightBlack, lipgloss.BrightBlack},
	"selectedEmptyBorder": {lipgloss.BrightGreen, lipgloss.Green},
	"redSuit":             {lipgloss.Red, lipgloss.BrightRed},
	"blackSuit":           {lipgloss.Black, lipgloss.BrightWhite},
	"cardBack":            {lipgloss.BrightBlack, lipgloss.White},
}
