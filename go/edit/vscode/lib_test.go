package vscode

import (
	"testing"
)

func TestSeek(t *testing.T) {
	h, err := NewFileHandler("testdata/test.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, err = h.offset(Position{Line: 2, Character: 3})
	if err != nil {
		t.Fatal(err)
	}
	// do it again and same result
	i, err := h.offset(Position{Line: 2, Character: 3})
	if err != nil {
		t.Fatal(err)
	}

	if i != 8 {
		t.Fatalf("expected %d but got %d", 5, i)
	}
}

func TestInsert(t *testing.T) {
	err := Insert("testdata/test.txt", Position{Line: 2, Character: 3}, "abc")
	if err != nil {
		t.Fatal(err)
	}

}
