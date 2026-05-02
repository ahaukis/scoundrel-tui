package game

import "fmt"

type Rank int

const (
	Jack Rank = iota + 11
	Queen
	King
	Ace
)

func (r Rank) String() string {
	if r < Jack {
		return fmt.Sprint(int(r))
	}
	switch r {
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	case Ace:
		return "A"
	}
	return "<unknown>"
}
