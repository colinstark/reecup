package game

import (
	"fmt"
	"testing"
)

func TestNewDeck(t *testing.T) {
	deck := CreateDeck()

	if len(deck.stones) != len(COLORS)*len(FACES)+2 {
		t.Error("New Deck does not have all stones")
	}

}

func TestDrawStone(t *testing.T) {
	deck := CreateDeck()
	initialLength := deck.GetCount()

	_, err := deck.Draw()
	if err != nil {
		fmt.Println(err)
		t.Error("Drawing stone failed")
	}

	if deck.GetCount() >= initialLength {
		t.Error("Drawing stone: length mismatch", initialLength, deck.GetCount())
	}

}

func TestShuffleDeck(t *testing.T) {
	deck := Deck{}
	deck.fillDeck()
	initialFirst := deck.stones[0]

	deck.shuffleDeck()

	if initialFirst == deck.stones[0] {
		t.Error("Shuffling didn't randomize order")
	}

}
