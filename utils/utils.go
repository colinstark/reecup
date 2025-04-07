package utils

import (
	"fmt"
	"math/rand"
)

func GenerateID() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, 4)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func GenerateIDFor(prefix string) string {
	randomID := GenerateID()
	return fmt.Sprintf("%v_%v", prefix, randomID)
}
