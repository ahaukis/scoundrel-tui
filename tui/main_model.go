package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/ahaukis/scoundrel-tui/game"
)

type mainModel struct {
	game              *game.Game
	gameTable         tableModel
	hasDarkBackground bool
}

func InitialMainModel() mainModel {
	g := game.NewRandomGame()
	table := tableModel{game: g, selectedRoomIdx: 0, weaponEnabled: true}
	return mainModel{game: g, gameTable: table}
}

func (m mainModel) Init() tea.Cmd {
	tableCmd := m.gameTable.Init()
	return tea.Batch(tea.RequestBackgroundColor, tableCmd)
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	tModel, tableCmd := m.gameTable.Update(msg)
	m.gameTable = tModel.(tableModel)
	cmds = append(cmds, tableCmd)

	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		// set background color based on
		m.hasDarkBackground = msg.IsDark()
		m.gameTable.hasDarkBackground = m.hasDarkBackground
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
	view := tea.NewView(m.gameTable.View().Content)
	view.AltScreen = true
	view.WindowTitle = "Scoundrel"

	return view
}
