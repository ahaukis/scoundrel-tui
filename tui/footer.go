package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/ahaukis/scoundrel-tui/game"
)

func footerLayer(hp int, weaponActive bool) *lipgloss.Layer {
	sBuilder := strings.Builder{}
	sBuilder.WriteString(fmt.Sprintf("[HP] %2d/%2d\n[weapon enabled] %t", hp, game.MaxHP, weaponActive))
	s := sBuilder.String()

	footerStyle := lipgloss.NewStyle().
		Render(s)

	return lipgloss.NewLayer(footerStyle)
}
