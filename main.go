package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/ahaukis/scoundrel-tui/tui"
)

func main() {
	p := tea.NewProgram(tui.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("An error occured:", err)
		os.Exit(1)
	}
}
