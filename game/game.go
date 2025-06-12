package game

import (
	"reecup/utils"
	"time"
)

type Game struct {
	ID                string    `json:"id"`
	StartedAt         time.Time `json:"startedAt"`
	Board             Board     `json:"board"`
	Deck              Deck
	Players           []Player `json:"players"`
	CurrentPlayerTurn *Player  `json:"currentPlayerTurn"`
	CurrentTurn       Turn     `json:"currentTurn"`
	GameOver          bool     `json:"gameOver"`
	State             string   `json:"state"`
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

// NextPlayerTurn advances to the next player's turn, cycling back to the first player
// if the current player is the last one in the array
func (g *Game) NextPlayerTurn() {
	if len(g.Players) == 0 {
		g.CurrentPlayerTurn = nil
		return
	}

	if g.CurrentPlayerTurn == nil {
		// If no current player, start with the first one
		g.CurrentPlayerTurn = &g.Players[0]
		g.StartNewTurn()
		return
	}

	// Find the current player's index
	currentIndex := -1
	for i, player := range g.Players {
		if player.ID == g.CurrentPlayerTurn.ID {
			currentIndex = i
			break
		}
	}

	// If current player not found, default to first player
	if currentIndex == -1 {
		g.CurrentPlayerTurn = &g.Players[0]
		g.StartNewTurn()
		return
	}

	// Move to next player, cycling back to 0 if at the end
	nextIndex := (currentIndex + 1) % len(g.Players)
	g.CurrentPlayerTurn = &g.Players[nextIndex]
	g.StartNewTurn()
}

// StartNewTurn initializes a new turn for the current player with the current board state
func (g *Game) StartNewTurn() {
	if g.CurrentPlayerTurn == nil {
		return
	}

	g.CurrentTurn = Turn{
		Player:    *g.CurrentPlayerTurn,
		TempBoard: g.Board, // Copy current board state to TempBoard
		StartedAt: time.Now(),
		IsValid:   false,
	}
}
