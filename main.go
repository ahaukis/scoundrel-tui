package main

import (
	"fmt"

	"github.com/ahaukis/scoundrel-tui/game"
)

func main() {
	c := game.NewCard(4, game.Hearts)
	c2 := game.NewCard(game.Ace, game.Diamonds)
	fmt.Println(c2.RanksAbove(c))

	g := game.NewRandomGame()
	fmt.Println(g)

}
