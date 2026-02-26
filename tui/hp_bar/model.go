package hpbar

import (
	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	"github.com/ahaukis/scoundrel-tui/game"
)

type Model struct {
	game    *game.Game
	progBar progress.Model
}

func New(g *game.Game) Model {
	progBar := progress.New(
		progress.WithWidth(20),
		progress.WithoutPercentage(),
	)
	return Model{g, progBar}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var hpCmd tea.Cmd
	m.progBar, hpCmd = m.progBar.Update(msg)
	cmds = append(cmds, hpCmd)

	return m, nil
}

func (m Model) View() string {
	p := float64(m.game.HP) / float64(game.MaxHP)
	s := m.progBar.ViewAs(p)
	return s
}
