package main

import (
	"testing"
)

func TestSimple(t *testing.T) {
	if true {
		t.Errorf("Invalid test")
	}
}
