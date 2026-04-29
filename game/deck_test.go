package game

import "testing"

func TestFullDeckSize(t *testing.T) {
	deck := NewFullDeck()
	size := len(deck)

	if size != DeckSize {
		t.Errorf("Expected %d cards, got %d", DeckSize, size)
	}
}

func TestScoundrelDeckSize(t *testing.T) {
	deck := NewScoundrelDeck()
	size := len(deck)
	scoundrel_size := DeckSize - 8

	if size != scoundrel_size {
		t.Errorf("Expected %d cards, got %d", scoundrel_size, size)
	}
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

	if sameOrder {
		t.Errorf("Shuffled decks are the same!")
	}
}
