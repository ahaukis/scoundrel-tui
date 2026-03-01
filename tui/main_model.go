package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/ahaukis/scoundrel-tui/game"
	hpbar "github.com/ahaukis/scoundrel-tui/tui/hp_bar"
	"github.com/ahaukis/scoundrel-tui/tui/palette"
	"github.com/ahaukis/scoundrel-tui/tui/table"
)

type mainModel struct {
	game      *game.Game
	gameTable table.Model
	hpBar     hpbar.Model
	palette   *palette.Palette
}

func InitialMainModel() mainModel {
	g := game.NewRandomGame()
	p := palette.NewDark()
	return mainModel{
		game:      g,
		gameTable: table.New(g, &p),
		hpBar:     hpbar.New(&g.HP, &p),
		palette:   &p,
	}
}

func (m mainModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, tea.RequestBackgroundColor)
	cmds = append(cmds, m.gameTable.Init())
	cmds = append(cmds, m.hpBar.Init())

	return tea.Batch(cmds...)
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	tModel, tableCmd := m.gameTable.Update(msg)
	m.gameTable = tModel
	cmds = append(cmds, tableCmd)

	hpModel, hpCmd := m.hpBar.Update(msg)
	m.hpBar = hpModel
	cmds = append(cmds, hpCmd)

	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		if msg.IsDark() {
			*m.palette = palette.NewDark()
		} else {
			*m.palette = palette.NewLight()
		}
	case tea.KeyPressMsg:
		if s := msg.String(); s == "ctrl+c" || s == "q" {
			return m, tea.Quit
		}
	}

	if m.game.IsRoomDone() {
		m.game.DealRoom()
	}

	return m, tea.Batch(cmds...)
}

func (m mainModel) View() tea.View {
	hpBarLayer := lipgloss.NewLayer(m.hpBar.View()).X(1).Y(1)
	mainLayer := lipgloss.NewLayer(m.gameTable.View()).
		Y(hpBarLayer.GetY() + hpBarLayer.Height() + 1)

	comp := lipgloss.NewCompositor(hpBarLayer, mainLayer)
	s := comp.Render() + "\n"

	view := tea.NewView(s)
	view.AltScreen = true
	view.WindowTitle = "Scoundrel"

	return view
}
