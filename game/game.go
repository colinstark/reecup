package game

import (
	"reecup/utils"
	"time"
)

type Game struct {
	ID          string    `json:"id"`
	StartedAt   time.Time `json:"startedAt"`
	Board       Board     `json:"board"`
	Deck        Deck
	Players     []Player `json:"players"`
	CurrentTurn Turn     `json:"currentTurn"`
	GameOver    bool     `json:"gameOver"`
	State       string   `json:"state"`
}

type Turn struct {
	Player    Player
	TempBoard Board
	StartedAt time.Time
	IsValid   bool
}

func NewGame() Game {
	return Game{
		ID:    "game_" + utils.GenerateID(),
		Deck:  CreateDeck(),
		Board: Board{},
		State: "waiting",
	}
}
