package game

import (
	"testing"
)

func TestCreateGame(t *testing.T) {
	game := NewGame()
	if game.ID == "" {
		t.Error("No Game ID")
	}

}
