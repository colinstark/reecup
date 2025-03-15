package game

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Deck struct {
	stones []Stone
}

func (d *Deck) shuffleDeck() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.stones), func(i, j int) { d.stones[i], d.stones[j] = d.stones[j], d.stones[i] })
}

func (d *Deck) GetCount() int {
	return len(d.stones)
}

func (d *Deck) fillDeck() {
	// jokers
	d.stones = append(d.stones, Stone{Color: "red", Joker: true})
	d.stones = append(d.stones, Stone{Color: "black", Joker: true})

	for _, color := range COLORS {
		for _, face := range FACES {
			d.stones = append(d.stones, Stone{Face: face, Color: color})
		}
	}
}

func (d *Deck) Draw() (Stone, error) {
	if len(d.stones) > 0 {
		return Stone{}, errors.New("No more stones")
	}
	stone := d.stones[0]
	d.stones = d.stones[1:]
	fmt.Println("drawing", stone)
	return stone, nil
}

func DrawForNewPlayer(name string, deck Deck) []Stone {
	hand := []Stone{}
	for i := 0; i < 14; i++ {
		stone, err := deck.Draw()
		if err != nil {
			// handle error, but unlikely
		}
		hand = append(hand, stone)
	}

	return hand
}

func CreateDeck() Deck {
	deck := Deck{}
	deck.fillDeck()
	deck.shuffleDeck()

	return deck
}
