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
	game           *game.Game
	palette        *palette.Palette
	gameTable      table.Model
	hpBar          hpbar.Model
	gameInProgress bool
	gameLost       bool
	gameWon        bool
}

func InitialMainModel() mainModel {
	g := game.NewRandomGame()
	p := palette.NewDark()
	return mainModel{
		game:           g,
		palette:        &p,
		gameTable:      table.New(g, &p),
		hpBar:          hpbar.New(&g.HP, &p),
		gameInProgress: true,
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

	var cmd tea.Cmd

	if !m.gameInProgress {
		cmd = m.updateMenu(msg)
	} else {
		cmd = m.updateGame(msg)
		m.gameWon = m.game.Won()
		m.gameLost = m.game.Lost()
		m.gameInProgress = !(m.gameWon || m.gameLost)
	}

	return m, cmd
}

func (m *mainModel) updateGame(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	tModel, tableCmd := m.gameTable.Update(msg)
	m.gameTable = tModel
	cmds = append(cmds, tableCmd)

	hpModel, hpCmd := m.hpBar.Update(msg)
	m.hpBar = hpModel
	cmds = append(cmds, hpCmd)

	m.gameWon = m.game.Won()
	m.gameLost = m.game.Lost()

	if m.game.IsRoomDone() {
		m.game.DealRoom()
	}

	return tea.Batch(cmds...)
}

func (m *mainModel) updateMenu(msg tea.Msg) tea.Cmd {
	switch msg.(type) {
	case tea.KeyPressMsg:
		// any key press starts a new game
		m.gameInProgress = true
		m.gameWon = false
		m.gameLost = false
		*m.game = *game.NewRandomGame()
		return m.gameTable.Init()
	}
	return nil
}

func (m mainModel) View() tea.View {
	var s string
	if m.gameInProgress {
		s = m.viewGame()
	} else {
		s = m.viewMenu()
	}
	view := tea.NewView(s)
	view.AltScreen = true
	view.WindowTitle = "Scoundrel"

	return view
}

func (m *mainModel) viewGame() string {
	hpBarLayer := lipgloss.NewLayer(m.hpBar.View()).X(1).Y(1)
	mainLayer := lipgloss.NewLayer(m.gameTable.View()).
		Y(hpBarLayer.GetY() + hpBarLayer.Height() + 1)

	comp := lipgloss.NewCompositor(hpBarLayer, mainLayer)
	s := comp.Render() + "\n"

	return s
}

func (m *mainModel) viewMenu() string {
	if m.gameWon {
		return "You won!"
	} else {
		return "You lost!"
	}
}
