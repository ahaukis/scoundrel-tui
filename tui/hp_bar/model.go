package hpbar

import (
	"fmt"
	"image/color"

	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/ahaukis/scoundrel-tui/game"
	"github.com/ahaukis/scoundrel-tui/tui/palette"
)

type Model struct {
	hp      *int
	progBar progress.Model
}

func New(hp *int) Model {
	progBar := progress.New(
		progress.WithFillCharacters('|', '·'),
		progress.WithWidth(20),
		progress.WithoutPercentage(),
		progress.WithColorFunc(progBarColor),
	)
	progBar.EmptyColor = palette.Colors["emptyBorder"][1]

	return Model{hp, progBar}
}

func progBarColor(total, current float64) color.Color {
	if total > 0.5 {
		return lipgloss.BrightGreen
	} else if total > 0.25 {
		return lipgloss.BrightYellow
	} else {
		return lipgloss.BrightRed
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var hpCmd tea.Cmd
	m.progBar, hpCmd = m.progBar.Update(msg)
	cmds = append(cmds, hpCmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	p := float64(*m.hp) / float64(game.MaxHP)
	s := fmt.Sprintf("HP %2d/%2d ", *m.hp, game.MaxHP) + m.progBar.ViewAs(p)

	return s
}
