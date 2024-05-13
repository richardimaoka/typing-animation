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
			err := c.pos.Validate()
			if err != nil && c.success {
				t.Fatalf("unexpected error %s", err)
			}
			if err == nil && !c.success {
				t.Fatalf("error expected but not encounterd")
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
		err      bool
	}{
		// line is zero-based, so the minimum is up to first line, copying no line is impossible
		"Up to the 1st line with line zero": {
			original,
			0,
			`Hello this is a test file.` + "\n",
			false,
		},
		"Up to the 2nd line": {
			original,
			1,
			`Hello this is a test file.
There are multiple lines` + "\n",
			false,
		},
		"Up to the 3rd line": {
			original,
			2,
			`Hello this is a test file.
There are multiple lines
in this text file.` + "\n",
			false,
		},
		"Up to the 4th line": {
			original,
			3,
			`Hello this is a test file.
There are multiple lines
in this text file.
And この文章のいくつかのpartは` + "\n",
			false,
		},
		"Up to the 5th line": {
			original,
			4,
			`Hello this is a test file.
There are multiple lines
in this text file.
And この文章のいくつかのpartは
英語とJapaneseを混ぜて` + "\n",
			false,
		},
		"Up to the 6th(last) line, without trailing new-line": {
			original,
			5,
			`Hello this is a test file.
There are multiple lines
in this text file.
And この文章のいくつかのpartは
英語とJapaneseを混ぜて
writtenされています。`,
			false,
		},
		//ERROR when trying to read more lines than exist
		"error upon reading non-existent 7th line": {
			original,
			6,
			"",
			true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var builder strings.Builder
			bufReader := bufio.NewReader(strings.NewReader(c.original))

			err := copyUpTo(bufReader, &builder, c.toLine)
			if err != nil {
				if c.err {
					return // expected error
				}
				t.Fatalf("unexpected error: %s", err)
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
		"at the beginning": {
			//              1         2
			//    01234567890123456789012345
			/**/ `Hello this is a test file.`,
			Position{Line: 0, Character: 0},
			"Good morning. ",
			`Good morning. Hello this is a test file.`,
			nil,
		},
		"in the middle, English": {
			//              1         2
			//    01234567890123456789012345
			/**/ `Hello this is a test file.`,
			Position{Line: 0, Character: 15},
			"n amazing",
			`Hello this is a` + "n amazing " + `test file.`,
			nil,
		},
		"in the middle, Japanese": {
			//                   1             2
			//    01234 5 6 7 8 90 1 2 34567 89012345
			/**/ `And この文章のいくつかのpartは`,
			Position{Line: 0, Character: 9},
			"中の",
			`And この文章の中のいくつかのpartは`,
			nil,
		},
		"at the beginning, end in newline": {
			//              1         2
			//    01234567890123456789012345
			/**/ `Hello this is a test file.\n`,
			Position{Line: 0, Character: 0},
			"Good morning. ",
			`Good morning. Hello this is a test file.\n`,
			nil,
		},
		"in the middle, English, end in newline": {
			//              1         2
			//    01234567890123456789012345
			/**/ `Hello this is a test file.\n`,
			Position{Line: 0, Character: 15},
			"n amazing",
			`Hello this is a` + "n amazing " + `test file.\n`,
			nil,
		},
		"in the middle, Japanese, end in newline": {
			//                   1             2
			//    01234 5 6 7 8 90 1 2 34567 89012345
			/**/ `And この文章のいくつかのpartは\n`,
			Position{Line: 0, Character: 9},
			"中の",
			`And この文章の中のいくつかのpartは\n`,
			nil,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := insertInLine(c.pos.Character, c.newText, []byte(c.original))
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

func TestProcessLine(t *testing.T) {

	cases := map[string]struct {
		original string
		pos      Position
		newText  string
		expected string
		err      bool
	}{
		"Up to the 1st line with line zero": {
			`Hello this is a test file.`,
			Position{-1, 0},
			"aaa",
			"",
			true,
		},
		"at the beginning": {
			//              1         2
			//    01234567890123456789012345
			/**/ `Hello this is a test file.`,
			Position{Line: 0, Character: 0},
			"Good morning. ",
			`Good morning. Hello this is a test file.`,
			false,
		},
		"in the middle, English": {
			//              1         2
			//    01234567890123456789012345
			/**/ `Hello this is a test file.`,
			Position{Line: 0, Character: 15},
			"n amazing",
			`Hello this is a` + "n amazing " + `test file.`,
			false,
		},
		"in the middle, Japanese": {
			//                   1             2
			//    01234 5 6 7 8 90 1 2 34567 89012345
			/**/ `And この文章のいくつかのpartは`,
			Position{Line: 0, Character: 9},
			"中の",
			`And この文章の中のいくつかのpartは`,
			false,
		},
		"at the beginning, end in newline": {
			//              1         2
			//    01234567890123456789012345
			/**/ `Hello this is a test file.\n`,
			Position{Line: 0, Character: 0},
			"Good morning. ",
			`Good morning. Hello this is a test file.\n`,
			false,
		},
		"in the middle, English, end in newline": {
			//              1         2
			//    01234567890123456789012345
			/**/ `Hello this is a test file.\n`,
			Position{Line: 0, Character: 15},
			"n amazing",
			`Hello this is a` + "n amazing " + `test file.\n`,
			false,
		},
		"in the middle, Japanese, end in newline": {
			//                   1             2
			//    01234 5 6 7 8 90 1 2 34567 89012345
			/**/ `And この文章のいくつかのpartは\n`,
			Position{Line: 0, Character: 9},
			"中の",
			`And この文章の中のいくつかのpartは\n`,
			false,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var builder strings.Builder
			bufReader := bufio.NewReader(strings.NewReader(c.original))

			err := processLine(bufReader, &builder, c.pos, c.newText)
			if err != nil {
				if c.err {
					return // expected error
				}
				t.Fatalf("unexpected error: %s", err)
			}

			result := builder.String()
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

// func TestInsert(t *testing.T) {
// 	err := Insert("testdata/test.txt", Position{Line: 2, Character: 3}, "abc")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	err = Insert("testdata/test.txt", Position{Line: 2, Character: 3}, "abc")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	err = Insert("testdata/test.txt", Position{Line: 2, Character: 3}, "abc")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestDelete(t *testing.T) {
// 	r := Range{Position{Line: 2, Character: 2}, Position{Line: 2, Character: 3}}
// 	err := Delete("testdata/test_delete.txt", r)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
