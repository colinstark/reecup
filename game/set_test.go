package game

import (
	"reflect"
	"testing"
)

func TestMinimumStones(t *testing.T) {
	set := Set{
		Stones: []Stone{
			Stone{Face: 1, Color: "yellow", Joker: false},
			Stone{Face: 3, Color: "yellow", Joker: false}},
	}

	if set.Validate() == true {
		t.Error("Minimum stone failed")
	}
}

func TestSequentialOneTwoThree(t *testing.T) {
	set := Set{
		Stones: []Stone{
			Stone{Face: 1, Color: "yellow", Joker: false},
			Stone{Face: 2, Color: "yellow", Joker: false},
			Stone{Face: 3, Color: "yellow", Joker: false}},
	}

	if set.Validate() == false {
		t.Error("Sequential detection failed")
	}
}

func TestSequentialFailing(t *testing.T) {
	set := Set{
		Stones: []Stone{
			Stone{Face: 2, Color: "yellow", Joker: false},
			Stone{Face: 1, Color: "yellow", Joker: false},
			Stone{Face: 3, Color: "yellow", Joker: false}},
	}

	if set.Validate() == true {
		t.Error("Sequential detection failed")
	}
}

func TestSequentialLoopAround(t *testing.T) {
	set := Set{
		Stones: []Stone{
			Stone{Face: 12, Color: "yellow", Joker: false},
			Stone{Face: 13, Color: "yellow", Joker: false},
			Stone{Face: 1, Color: "yellow", Joker: false}},
	}

	if set.Validate() == false {
		t.Error("Looping not working")
	}
}

func TestSequentialJoker(t *testing.T) {
	set := Set{
		Stones: []Stone{
			Stone{Face: 0, Color: "yellow", Joker: true},
			Stone{Face: 5, Color: "yellow", Joker: false},
			Stone{Face: 6, Color: "yellow", Joker: false},
		},
	}

	if set.Validate() == false {
		t.Error("Joker not working")
	}
}

func TestSequentialJokerInTheMiddle(t *testing.T) {
	set := Set{
		Stones: []Stone{
			Stone{Face: 4, Color: "yellow", Joker: false},
			Stone{Face: 0, Color: "yellow", Joker: true},
			Stone{Face: 6, Color: "yellow", Joker: false},
		},
	}

	if set.Validate() == false {
		t.Error("Joker in the middle not working")
	}
}

func TestRun(t *testing.T) {
	set := Set{
		Stones: []Stone{
			Stone{Face: 9, Color: "red", Joker: false},
			Stone{Face: 9, Color: "yellow", Joker: false},
			Stone{Face: 9, Color: "blue", Joker: false},
		},
	}

	if set.Validate() == false {
		t.Error("Run test failed")
	}
}

func TestRunJoker(t *testing.T) {
	set := Set{
		Stones: []Stone{
			Stone{Face: 9, Color: "red", Joker: false},
			Stone{Face: 9, Color: "red", Joker: true},
			Stone{Face: 9, Color: "blue", Joker: false},
		},
	}

	if set.Validate() == false {
		t.Error("Run test w. Joker failed")
	}
}

func TestRunFailing(t *testing.T) {
	set := Set{
		Stones: []Stone{
			Stone{Face: 9, Color: "red", Joker: false},
			Stone{Face: 10, Color: "yellow", Joker: false},
			Stone{Face: 9, Color: "blue", Joker: false},
		},
	}

	if set.Validate() == true {
		t.Error("Run failing failed")
	}
}

func TestRunColorFailing(t *testing.T) {
	set := Set{
		Stones: []Stone{
			Stone{Face: 9, Color: "red", Joker: false},
			Stone{Face: 9, Color: "yellow", Joker: false},
			Stone{Face: 9, Color: "blue", Joker: false},
			Stone{Face: 9, Color: "black", Joker: false},
			Stone{Face: 9, Color: "yellow", Joker: false},
		},
	}

	if set.Validate() == true {
		t.Error("Run Color failing failed")
	}
}

func TestAddingStone(t *testing.T) {
	set := Set{
		Stones: []Stone{
			Stone{Face: 9, Color: "red", Joker: false},
			Stone{Face: 9, Color: "yellow", Joker: false},
			Stone{Face: 9, Color: "blue", Joker: false},
		},
	}
	targetSet := Set{
		Stones: []Stone{
			Stone{Face: 9, Color: "red", Joker: false},
			Stone{Face: 9, Color: "yellow", Joker: false},
			Stone{Face: 9, Color: "black", Joker: false},
			Stone{Face: 9, Color: "blue", Joker: false},
		},
	}

	newStone := Stone{Face: 9, Color: "black", Joker: false}

	set.AddStone(2, newStone)
	equal := reflect.DeepEqual(set, targetSet)

	if !equal {
		t.Error("Adding stone failed")
	}
}

func TestAddingStoneAtZero(t *testing.T) {
	set := Set{
		Stones: []Stone{
			Stone{Face: 9, Color: "red", Joker: false},
			Stone{Face: 9, Color: "yellow", Joker: false},
			Stone{Face: 9, Color: "blue", Joker: false},
		},
	}
	targetSet := Set{
		Stones: []Stone{
			Stone{Face: 9, Color: "black", Joker: false},
			Stone{Face: 9, Color: "red", Joker: false},
			Stone{Face: 9, Color: "yellow", Joker: false},
			Stone{Face: 9, Color: "blue", Joker: false},
		},
	}

	newStone := Stone{Face: 9, Color: "black", Joker: false}
	set.AddStone(0, newStone)

	equal := reflect.DeepEqual(set, targetSet)

	if !equal {
		t.Error("Adding stone at zero failed")
	}
}

func TestRemovingStone(t *testing.T) {
	set := Set{
		Stones: []Stone{
			Stone{Face: 9, Color: "red", Joker: false},
			Stone{Face: 9, Color: "yellow", Joker: false},
			Stone{Face: 9, Color: "blue", Joker: false},
		},
	}
	targetSet := Set{
		Stones: []Stone{
			Stone{Face: 9, Color: "red", Joker: false},
			Stone{Face: 9, Color: "blue", Joker: false},
		},
	}

	set.RemoveStone(1)
	equal := reflect.DeepEqual(set, targetSet)

	if !equal {
		t.Error("Removing stone failed")
	}
}

func TestRemovingStoneAtZero(t *testing.T) {
	set := Set{
		Stones: []Stone{
			Stone{Face: 9, Color: "red", Joker: false},
			Stone{Face: 9, Color: "yellow", Joker: false},
			Stone{Face: 9, Color: "blue", Joker: false},
		},
	}
	targetSet := Set{
		Stones: []Stone{
			Stone{Face: 9, Color: "yellow", Joker: false},
			Stone{Face: 9, Color: "blue", Joker: false},
		},
	}

	set.RemoveStone(0)
	equal := reflect.DeepEqual(set, targetSet)

	if !equal {
		t.Error("Removing stone failed")
	}
}

func TestMovingStone(t *testing.T) {
	set := Set{
		Stones: []Stone{
			Stone{Face: 9, Color: "red", Joker: false},
			Stone{Face: 9, Color: "yellow", Joker: false},
			Stone{Face: 9, Color: "blue", Joker: false},
			Stone{Face: 9, Color: "black", Joker: false},
		},
	}
	targetSet := Set{
		Stones: []Stone{
			Stone{Face: 9, Color: "red", Joker: false},
			Stone{Face: 9, Color: "blue", Joker: false},
			Stone{Face: 9, Color: "yellow", Joker: false},
			Stone{Face: 9, Color: "black", Joker: false},
		},
	}

	set.MoveStone(2, 1)

	equal := reflect.DeepEqual(set, targetSet)

	if !equal {
		t.Error("Moving stone failed")
	}
}
func TestMovingStoneToZero(t *testing.T) {
	set := Set{
		Stones: []Stone{
			Stone{Face: 9, Color: "red", Joker: false},
			Stone{Face: 9, Color: "yellow", Joker: false},
			Stone{Face: 9, Color: "blue", Joker: false},
			Stone{Face: 9, Color: "black", Joker: false},
		},
	}
	targetSet := Set{
		Stones: []Stone{
			Stone{Face: 9, Color: "yellow", Joker: false},
			Stone{Face: 9, Color: "red", Joker: false},
			Stone{Face: 9, Color: "blue", Joker: false},
			Stone{Face: 9, Color: "black", Joker: false},
		},
	}

	set.MoveStone(1, 0)

	equal := reflect.DeepEqual(set, targetSet)

	if !equal {
		t.Error("Moving stone failed")
	}
}
