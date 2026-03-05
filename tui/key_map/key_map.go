package keymap

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	"github.com/ahaukis/scoundrel-tui/game"
)

type KeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding

	Interact     key.Binding
	SelectInRoom key.Binding
	SkipRoom     key.Binding
	ToggleWeapon key.Binding

	Help key.Binding
	Quit key.Binding
}

func New() KeyMap {
	roomIdcs := make([]string, game.CardsPerRoom)
	for i := range len(roomIdcs) {
		roomIdcs[i] = fmt.Sprint(i + 1)
	}

	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "move right"),
		),

		Interact: key.NewBinding(
			key.WithKeys("space", "enter"),
			key.WithHelp("⏎/␣", "interact"),
		),
		SelectInRoom: key.NewBinding(
			key.WithKeys(roomIdcs...),
			key.WithHelp(fmt.Sprintf("%s-%s", roomIdcs[0], roomIdcs[len(roomIdcs)-1]), "select room card"),
		),
		SkipRoom: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "skip room"),
		),
		ToggleWeapon: key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "toggle weapon"),
		),

		Help: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},                          // first col
		{k.Interact, k.SelectInRoom, k.SkipRoom, k.ToggleWeapon}, // second col ...
		{k.Help, k.Quit},
	}
}
