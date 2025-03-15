package utils

import "testing"

func TestRandomUUID(t *testing.T) {
	id1 := GenerateID()
	id2 := GenerateID()

	if id1 == id2 {
		t.Error("Random ID failed")
	}
}
