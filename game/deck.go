package game

import "math/rand/v2"

// Create a new full deck of 52 cards.
func NewFullDeck() []*Card {
	cards := make([]*Card, 0, DeckSize)

	for r := MinRank; r <= MaxRank; r++ {
		for s := MinSuit; s <= MaxSuit; s++ {
			cards = append(cards, NewCard(r, s))
		}
	}

	return cards
}

// Create a new deck with red suits of ranks J,Q,K,A removed, suitable for playing Scoundrel.
func NewScoundrelDeck() []*Card {
	cards := make([]*Card, 0, DeckSize)

	for r := MinRank; r <= MaxRank; r++ {
		for s := MinSuit; s <= MaxSuit; s++ {
			if r < Jack || !s.IsRed() {
				cards = append(cards, NewCard(r, s))
			}
		}
	}

	return cards
}

// Create a new shuffled deck with red suits of ranks J,Q,K,A removed, suitable for playing Scoundrel.
func NewShuffledScoundrelDeck() []*Card {
	d := NewScoundrelDeck()
	shuffleDeck(d)
	return d
}

func shuffleDeck(d []*Card) {
	rand.Shuffle(len(d), func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})
}
