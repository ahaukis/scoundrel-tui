package game

import "math/rand/v2"

func NewFullDeck() []*Card {
	cards := make([]*Card, 0, DeckSize)

	for r := MinRank; r <= MaxRank; r++ {
		for s := MinSuit; s <= MaxSuit; s++ {
			cards = append(cards, NewCard(r, s))
		}
	}

	return cards
}

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

func NewShuffledFullDeck() []*Card {
	d := NewFullDeck()
	rand.Shuffle(len(d), func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})
	return d
}
