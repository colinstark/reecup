package utils

import (
	"strings"
	"testing"
)

func TestRandomUUID(t *testing.T) {
	id1 := GenerateID()
	id2 := GenerateID()

	if len(id1) != 4 {
		t.Error("Length is not the expected 4")
	}

	if id1 == id2 {
		t.Error("Random ID failed")
	}
}

func TestRandomUserID(t *testing.T) {
	userID := GenerateIDFor("user")

	if !strings.Contains(userID, "user_") {
		t.Error("Random User ID failed")
	}
}
