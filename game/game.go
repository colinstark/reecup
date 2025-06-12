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

func (g *Game) NextPlayerTurn() {
	if len(g.Players) == 0 {
		g.CurrentPlayerTurn = nil
		return
	}

	if g.CurrentPlayerTurn == nil {
		g.CurrentPlayerTurn = &g.Players[0]
		g.StartNewTurn()
		return
	}

	currentIndex := -1
	for i, player := range g.Players {
		if player.ID == g.CurrentPlayerTurn.ID {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		g.CurrentPlayerTurn = &g.Players[0]
		g.StartNewTurn()
		return
	}

	nextIndex := (currentIndex + 1) % len(g.Players)
	g.CurrentPlayerTurn = &g.Players[nextIndex]
	g.StartNewTurn()
}

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

func (g *Game) calculateBoardValue(board Board) int {
	totalValue := 0

	for _, set := range board.Sets {
		for _, stone := range set.Stones {
			if !stone.Joker {
				totalValue += stone.Face
			}
		}
	}

	return totalValue
}

func (g *Game) CheckMeld() bool {
	originalValue := g.calculateBoardValue(g.Board)
	tempValue := g.calculateBoardValue(g.CurrentTurn.TempBoard)

	difference := tempValue - originalValue
	return difference > 30
}

func (g *Game) IsTempBoardPoolEmpty() bool {
	return len(g.CurrentTurn.TempBoard.Pool) == 0
}
