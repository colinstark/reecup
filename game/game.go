package game

import (
	"math/rand"
)

type Game struct {
	id       string
	Board    Board
	Deck     Deck
	Players  []Player
	GameOver bool
}

func NewGame() Game {
	return Game{
		id:    generateID(),
		Deck:  CreateDeck(),
		Board: Board{},
	}
}

func generateID() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, 4)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
