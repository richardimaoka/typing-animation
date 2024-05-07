package vscode

import (
	"testing"
)

func TestValidatePosition(t *testing.T) {
	cases := map[string]struct {
		success bool
		pos     Position
	}{
		"negative char": {false, Position{Character: -1, Line: 10}},
		"zero char":     {true, Position{Character: 0, Line: 10}},
		"negative line": {false, Position{Character: 1, Line: -10}},
		"zero line":     {true, Position{Character: 1, Line: 0}},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result := c.pos.Validate()
			if c.success != result {
				t.Errorf("%t expected but result = %t", c.success, result)
			}
		})
	}
}

func TestLessThanOrEqualToPosition(t *testing.T) {
	cases := map[string]struct {
		expected bool
		p1       Position
		p2       Position
	}{
		"p1 line above p2 line":           {true, Position{Character: 5, Line: 10}, Position{Character: 1, Line: 12}},
		"p1 line below p2 line":           {false, Position{Character: 1, Line: 10}, Position{Character: 5, Line: 8}},
		"p1 left to p2 on the same line":  {true, Position{Character: 0, Line: 10}, Position{Character: 1, Line: 10}},
		"p1 right to p2 on the same line": {false, Position{Character: 5, Line: 10}, Position{Character: 1, Line: 10}},
		"p1 and p2 same position":         {true, Position{Character: 0, Line: 10}, Position{Character: 0, Line: 10}},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result := c.p1.LessThanOrEqualTo(c.p2)
			if c.expected != result {
				t.Errorf("%t expected but result = %t", c.expected, result)
			}
		})
	}
}

func TestValidateRange(t *testing.T) {
	// cases := map[string]struct {
	// 	success bool
	// 	r       Range
	// }{
	// 	"negative pos": {true, Range{Position{Character: -1, Line: 10}, Position{Character: 1, Line: 10}}},
	// }

	// for name, c := range cases {
	// 	t.Run(name, func(t *testing.T) {
	// 		result := c.pos.Validate()
	// 		if c.success != result {
	// 			t.Errorf("%t expected but result = %t", c.success, result)
	// 		}
	// 	})
	// }
}

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
	err = Insert("testdata/test.txt", Position{Line: 2, Character: 3}, "abc")
	if err != nil {
		t.Fatal(err)
	}
	err = Insert("testdata/test.txt", Position{Line: 2, Character: 3}, "abc")
	if err != nil {
		t.Fatal(err)
	}
}

// func TestDelete(t *testing.T) {
// 	r := Range{Position{Line: 2, Character: 2}, Position{Line: 2, Character: 3}}
// 	err := Delete("testdata/test_delete.txt", r)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
