package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullDeckSize(t *testing.T) {
	deck := NewFullDeck()
	assert.Len(t, deck, DeckSize)

}

func TestScoundrelDeckSize(t *testing.T) {
	deck := NewScoundrelDeck()
	assert.Len(t, deck, DeckSize-8)
}

func TestShuffleDeck(t *testing.T) {
	d1 := NewFullDeck()
	d2 := NewFullDeck()
	shuffleDeck(d1)
	shuffleDeck(d2)

	sameOrder := true

	for i := range len(d1) {
		if *d1[i] != *d2[i] {
			sameOrder = false
			break
		}
	}

	assert.False(t, sameOrder, "Shuffled decks are the same!")
}
