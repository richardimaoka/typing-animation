package vscode

import (
	"bufio"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

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

			err := copyUpToLine(bufReader, &builder, c.toLine)
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
		err      bool
	}{
		"at the beginning":                        {"0123456789" /********/, Position{Line: 0, Character: 0}, "Insert ", "Insert 0123456789", false},
		"in the middle, 1":                        {"0123456789" /********/, Position{Line: 0, Character: 1}, " insert ", "0 insert 123456789", false},
		"in the middle, 2":                        {"0123456789" /********/, Position{Line: 0, Character: 2}, " insert ", "01 insert 23456789", false},
		"in the middle, 3":                        {"0123456789" /********/, Position{Line: 0, Character: 3}, " insert ", "012 insert 3456789", false},
		"in the middle, Japanese":                 {"012三四五六七89" /****/, Position{Line: 0, Character: 3}, " 中間 ", "012 中間 三四五六七89", false},
		"at the beginning, end in newline":        {"0123456789\n" /******/, Position{Line: 0, Character: 0}, "Insert ", "Insert 0123456789\n", false},
		"in the middle, English, end in newline":  {"0123456789\n" /******/, Position{Line: 0, Character: 3}, " insert ", "012 insert 3456789\n", false},
		"in the middle, Japanese, end in newline": {"012三四五六七89\n" /**/, Position{Line: 0, Character: 7}, " 中間 ", "012三四五六 中間 七89\n", false},
		"at the end, Japanese, end in newline":    {"012三四五六七89\n" /**/, Position{Line: 0, Character: 10}, " 最後", "012三四五六七89 最後\n", false},
		// error cases
		// "at the end, Japanese, after newline": {"012三四五六七89\n" /**/, Position{Line: 0, Character: 11}, " 最後より後", "", true},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := insertInLine(c.pos.Character, c.newText, []byte(c.original))
			if err != nil {
				if c.err {
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

func TestCopyUntilEOF(t *testing.T) {
	cases := map[string]struct {
		original string
		err      bool
	}{
		"empty": {
			``,
			false,
		},
		"only newline": {
			"\n",
			false,
		},
		"one line": {
			`Hello this is a test file.`,
			false,
		},
		"one line, endling in newline": {
			`Hello this is a test file.` + "\n",
			false,
		},
		"two lines": {
			`Hello this is a test file.
Good morning`,
			false,
		},
		"two lines ending in newline": {
			`Hello this is a test file.
Good morning` + "\n",
			false,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var builder strings.Builder
			bufReader := bufio.NewReader(strings.NewReader(c.original))

			err := copyUntilEOF(bufReader, &builder)
			if err != nil {
				if c.err {
					return // expected error
				}
				t.Fatalf("unexpected error: %s", err)
			}

			result := builder.String()
			if c.original != result {
				t.Errorf("%s", cmp.Diff(c.original, result))
			}
		})
	}
}
