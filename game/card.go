package game

import "fmt"

type Card struct {
	rank Rank
	suit Suit
}

// Creates a new card. Guards againts invalid rank/suit values.
func NewCard(rank Rank, suit Suit) *Card {
	if rank < MinRank || rank > MaxRank {
		return &Card{}
	}
	if suit < MinSuit || suit > MaxSuit {
		return &Card{}
	}
	return &Card{rank, suit}
}

// Get the card's rank as an integer.
func (c *Card) IntRank() int {
	return int(c.rank)
}

func (c *Card) RanksEqual(other *Card) bool {
	return c.rank == other.rank
}

func (c *Card) RanksAbove(other *Card) bool {
	return c.rank > other.rank
}

func (c *Card) RanksBelow(other *Card) bool {
	return c.rank < other.rank
}

func (c *Card) IsWeapon() bool {
	return c.suit == Diamonds
}

func (c *Card) IsHealthPotion() bool {
	return c.suit == Hearts
}

func (c *Card) IsMonster() bool {
	return c.suit.IsBlack()
}

func (c Card) String() string {
	return fmt.Sprintf("%v%v", c.rank, c.suit)
}
