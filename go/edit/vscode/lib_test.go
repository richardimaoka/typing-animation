package vscode

import (
	"bufio"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
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

func TestCopyUpTo(t *testing.T) {
	original := `Hello this is a test file.
There are multiple lines
in this text file.
And この文章のいくつかのpartは
英語とJapaneseを混ぜて
writtenされています。`

	cases := map[string]struct {
		original string
		toLine   int
		expected string
	}{
		"0": {
			original,
			0,
			`Hello this is a test file.` + "\n",
		},
		"1": {
			original,
			1,
			`Hello this is a test file.
There are multiple lines` + "\n",
		},
		"2": {
			original,
			2,
			`Hello this is a test file.
There are multiple lines
in this text file.` + "\n",
		},
		"3": {
			original,
			3,
			`Hello this is a test file.
There are multiple lines
in this text file.
And この文章のいくつかのpartは` + "\n",
		},
		"4": {
			original,
			4,
			`Hello this is a test file.
There are multiple lines
in this text file.
And この文章のいくつかのpartは
英語とJapaneseを混ぜて` + "\n",
		},
		"5": {
			original,
			5,
			`Hello this is a test file.
There are multiple lines
in this text file.
And この文章のいくつかのpartは
英語とJapaneseを混ぜて
writtenされています。`,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var builder strings.Builder
			bufReader := bufio.NewReader(strings.NewReader(c.original))
			err := copyUpTo(bufReader, &builder, c.toLine)
			if err != nil {
				t.Errorf("error: %s", err)
			}

			result := builder.String()
			if c.expected != result {
				t.Errorf("%s", cmp.Diff(c.expected, result))
			}
		})
	}
}

func TestIinsertInLine(t *testing.T) {
	cases := map[string]struct {
		original string
		pos      Position
		newText  string
		expected string
		err      error
	}{
		"0": {
			//              1         2
			//    01234567890123456789012345
			/**/ `Hello this is a test file.`,
			Position{Line: 0, Character: 0},
			"Good morning. ",
			`Good morning. Hello this is a test file.`,
			nil,
		},
		"1": {
			//              1         2
			//    01234567890123456789012345
			/**/ `Hello this is a test file.`,
			Position{Line: 0, Character: 15},
			"n amazing",
			`Hello this is a` + "n amazing " + `test file.`,
			nil,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := insertInLine(c.pos, c.newText, []byte(c.original))
			if err != nil {
				if c.err == nil {
					t.Fatalf("error: %s", err)
				}
				return // expected error
			}

			if c.expected != result {
				t.Errorf("%s", cmp.Diff(c.expected, result))
			}
		})
	}
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
