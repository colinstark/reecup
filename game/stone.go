package game

var COLORS = []string{"blue", "red", "yellow", "black"}
var FACES = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

type Stone struct {
	Face  int    `json:"face"`
	Color string `json:"color"`
	Joker bool   `json:"joker"`
}
