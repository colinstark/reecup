package main

import (
	"fmt"
	"reecup/game"
)

func main() {
	fmt.Println("Welcome to ReeCup")
	new_game := game.NewGame()

	new_game.Players = append(new_game.Players, game.NewPlayer("Pauly", new_game.Deck))
	new_game.Players = append(new_game.Players, game.NewPlayer("Christopher", new_game.Deck))

	fmt.Println("created game", new_game)

}
