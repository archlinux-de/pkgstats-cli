package main

import (
	"os"
	"testing"
)

var Mocks = make(map[string]func())

func TestMain(m *testing.M) {
	mockName := os.Getenv("TEST_MOCK")
	if mockName != "" {
		mock, ok := Mocks[mockName]
		if ok {
			mock()
		}
	}

	os.Exit(m.Run())
}
