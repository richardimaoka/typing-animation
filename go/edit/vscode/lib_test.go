package vscode

import (
	"testing"
)

func TestSeek(t *testing.T) {
	h, err := NewFileHandler("testdata/test.txt")
	if err != nil {
		t.Fatal(err)
	}
	i, err := h.byteIndexAt(Position{Line: 2, Character: 3})
	if err != nil {
		t.Fatal(err)
	}

	if i != 8 {
		t.Fatalf("expected %d but got %d", 5, i)
	}
}
