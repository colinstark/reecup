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
	GameOver    bool     `json:"gameOver"`
	CurrentTurn Turn     `json:"currentTurn"`
	InProgress  bool     `json:"inProgress"`
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
	}
}
