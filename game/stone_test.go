package game

import (
	"slices"
	"testing"
)

func TestStonesAdhereToDefinitions(t *testing.T) {
	stone := Stone{Face: 1, Color: "yellow", Joker: false}

	colorValid := slices.Contains(COLORS, stone.Color)
	faceValid := slices.Contains(FACES, stone.Face)

	if !stone.Joker {
		if !colorValid || !faceValid {
			t.Error("Stone uses weird parameters")
		}
	}

	stone = Stone{Face: 25, Color: "yellow", Joker: false}

	faceValid = slices.Contains(FACES, stone.Face)
	if faceValid {
		t.Error("Invalid face passes")
	}

}
