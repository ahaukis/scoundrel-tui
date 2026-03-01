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
	palette *palette.Palette
}

func New(hp *int, palette *palette.Palette) Model {
	progBar := progress.New(
		progress.WithFillCharacters('|', '·'),
		progress.WithWidth(20),
		progress.WithoutPercentage(),
		progress.WithColorFunc(progBarColorFunc(palette)),
	)
	progBar.EmptyColor = palette.EmptyBorder

	return Model{hp, progBar, palette}
}

func progBarColorFunc(palette *palette.Palette) func(total, current float64) color.Color {
	return func(total, current float64) color.Color {
		if total > 0.5 {
			return palette.HPFull
		} else if total > 0.25 {
			return palette.HPMid
		} else {
			return palette.HPLow
		}
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
	lbl := lipgloss.NewStyle().Foreground(m.palette.Border).Render(
		fmt.Sprintf("HP %2d/%2d ", *m.hp, game.MaxHP))

	return lbl + m.progBar.ViewAs(p)
}
