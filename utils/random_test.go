package utils

import (
	"errors"
	"testing"
)

func TestGenerateRandomStringFailsClosed(t *testing.T) {
	old := readRandom
	readRandom = func([]byte) (int, error) { return 0, errors.New("random source unavailable") }
	t.Cleanup(func() { readRandom = old })

	defer func() {
		if recover() == nil {
			t.Fatal("random source failure must not produce a token")
		}
	}()
	GenerateRandomString(32)
}
