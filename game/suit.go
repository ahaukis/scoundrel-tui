package game

type Suit int

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

func (s Suit) IsRed() bool {
	return s == Diamonds || s == Hearts
}

func (s Suit) IsBlack() bool {
	return !s.IsRed()
}

func (s Suit) String() string {
	switch s {
	case Clubs:
		return "♣"
	case Diamonds:
		return "♦"
	case Hearts:
		return "♥"
	case Spades:
		return "♠"
	}
	return "<unknown>"
}
