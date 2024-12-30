package game

import (
	"fmt"
	"math/rand"
	"time"
)

type Stone struct {
	face  int
	color string
	joker bool
}

type Deck struct {
	stones []Stone
}

func (d *Deck) shuffleDeck() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.stones), func(i, j int) { d.stones[i], d.stones[j] = d.stones[j], d.stones[i] })
}

func (d *Deck) fillDeck() {
	colors := [...]string{"blue", "red", "yellow", "black"}
	faces := [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

	// jokers
	d.stones = append(d.stones, Stone{color: "red", joker: true})
	d.stones = append(d.stones, Stone{color: "black", joker: true})

	for _, color := range colors {
		for _, face := range faces {
			d.stones = append(d.stones, Stone{face: face, color: color})
		}
	}

}

func (d *Deck) Draw() (Stone, error) {
	// if len(d.stones) > 0 {}
	stone := d.stones[0]
	d.stones = d.stones[1:]
	fmt.Println("drawing", stone)
	return stone, nil
}

func CreateDeck() Deck {
	deck := Deck{}
	deck.fillDeck()
	deck.shuffleDeck()

	return deck
}
