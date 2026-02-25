// Custom dashed border for lipgloss.

package tui

import "charm.land/lipgloss/v2"

// A dashed border with rounded corners.
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
