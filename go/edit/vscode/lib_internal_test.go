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
		"at the beginning":                             {"0123456789" /********/, Position{Line: 0, Character: 0}, "Insert ", "Insert 0123456789", false},
		"at the beginning, end in newline":             {"0123456789\n" /******/, Position{Line: 0, Character: 0}, "Insert ", "Insert 0123456789\n", false},
		"in the middle, 1":                             {"0123456789" /********/, Position{Line: 0, Character: 1}, " insert ", "0 insert 123456789", false},
		"in the middle, 2":                             {"0123456789" /********/, Position{Line: 0, Character: 2}, " insert ", "01 insert 23456789", false},
		"in the middle, 3":                             {"0123456789" /********/, Position{Line: 0, Character: 3}, " insert ", "012 insert 3456789", false},
		"in the middle, Japanese":                      {"012三四五六七89" /****/, Position{Line: 0, Character: 3}, " 中間 ", "012 中間 三四五六七89", false},
		"in the middle, English, end in newline":       {"0123456789\n" /******/, Position{Line: 0, Character: 3}, " insert ", "012 insert 3456789\n", false},
		"in the middle, Japanese, end in newline":      {"012三四五六七89\n" /**/, Position{Line: 0, Character: 7}, " 中間 ", "012三四五六 中間 七89\n", false},
		"close to the end, Japanese":                   {"012三四五六七89" /****/, Position{Line: 0, Character: 9}, " 中間 ", "012三四五六七8 中間 9", false},
		"at the end, Japanese":                         {"012三四五六七89" /****/, Position{Line: 0, Character: 10}, " 最後", "012三四五六七89 最後", false},
		"ERROR: at the end, Japanese, after end 1":     {"012三四五六七89" /****/, Position{Line: 0, Character: 11}, " 最後より後", "", true},
		"ERROR: at the end, Japanese, after end 2":     {"012三四五六七89" /****/, Position{Line: 0, Character: 12}, " 最後より後", "", true},
		"at the end, Japanese, end in newline":         {"012三四五六七89\n" /**/, Position{Line: 0, Character: 10}, " 最後", "012三四五六七89 最後\n", false},
		"ERROR: at the end, Japanese, after newline 1": {"012三四五六七89\n" /**/, Position{Line: 0, Character: 11}, " 最後より後", "", true},
		"ERROR: at the end, Japanese, after newline 2": {"012三四五六七89\n" /**/, Position{Line: 0, Character: 12}, " 最後より後", "", true},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := insertInLine(c.pos.Character, c.newText, []byte(c.original))
			if err != nil {
				if c.err {
					return // expected error
				}
				t.Fatalf("error: %s", err)
			}

			if c.err {
				t.Fatalf("Expected error: but succeeded with result = %s", result)
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
		"at the beginning":                             {"0123456789" /********/, Position{Line: 0, Character: 0}, "Insert ", "Insert 0123456789", false},
		"at the beginning, end in newline":             {"0123456789\n" /******/, Position{Line: 0, Character: 0}, "Insert ", "Insert 0123456789\n", false},
		"in the middle, 1":                             {"0123456789" /********/, Position{Line: 0, Character: 1}, " insert ", "0 insert 123456789", false},
		"in the middle, 2":                             {"0123456789" /********/, Position{Line: 0, Character: 2}, " insert ", "01 insert 23456789", false},
		"in the middle, 3":                             {"0123456789" /********/, Position{Line: 0, Character: 3}, " insert ", "012 insert 3456789", false},
		"in the middle, Japanese":                      {"012三四五六七89" /****/, Position{Line: 0, Character: 3}, " 中間 ", "012 中間 三四五六七89", false},
		"in the middle, English, end in newline":       {"0123456789\n" /******/, Position{Line: 0, Character: 3}, " insert ", "012 insert 3456789\n", false},
		"in the middle, Japanese, end in newline":      {"012三四五六七89\n" /**/, Position{Line: 0, Character: 7}, " 中間 ", "012三四五六 中間 七89\n", false},
		"close to the end, Japanese":                   {"012三四五六七89" /****/, Position{Line: 0, Character: 9}, " 中間 ", "012三四五六七8 中間 9", false},
		"at the end, Japanese":                         {"012三四五六七89" /****/, Position{Line: 0, Character: 10}, " 最後", "012三四五六七89 最後", false},
		"ERROR: at the end, Japanese, after end 1":     {"012三四五六七89" /****/, Position{Line: 0, Character: 11}, " 最後より後", "", true},
		"ERROR: at the end, Japanese, after end 2":     {"012三四五六七89" /****/, Position{Line: 0, Character: 12}, " 最後より後", "", true},
		"at the end, Japanese, end in newline":         {"012三四五六七89\n" /**/, Position{Line: 0, Character: 10}, " 最後", "012三四五六七89 最後\n", false},
		"ERROR: at the end, Japanese, after newline 1": {"012三四五六七89\n" /**/, Position{Line: 0, Character: 11}, " 最後より後", "", true},
		"ERROR: at the end, Japanese, after newline 2": {"012三四五六七89\n" /**/, Position{Line: 0, Character: 12}, " 最後より後", "", true},
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
			if c.err {
				t.Fatalf("Expected error: but succeeded with result = %s", result)
			}
			if c.expected != result {
				t.Errorf("%s", cmp.Diff(c.expected, result))
			}
		})
	}
}

func TestReadUpToPrevChar(t *testing.T) {
	cases := map[string]struct {
		original string
		charAt   int
		expected string
		err      bool
	}{
		"0":           {"0123456789\n", 0, "", false},
		"1":           {"0123456789\n", 1, "0", false},
		"3":           {"0123456789\n", 3, "012", false},
		"9":           {"0123456789\n", 9, "012345678", false},
		"10":          {"0123456789\n", 10, "0123456789", false},
		"ERROR: 11":   {"0123456789\n", 11, "0123456789", true}, // error expected, reading the newline char is not allowed by this method
		"4 Japanese":  {"012三四五六七八九\n", 4, "012三", false},
		"5 Japanese":  {"012三四五六七八九\n", 5, "012三四", false},
		"6 Japanese":  {"012三四五六七八九\n", 6, "012三四五", false},
		"10 Japanese": {"012三四五六七八九\n", 10, "012三四五六七八九", false},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			line := []byte(c.original)
			result, err := readUptoPrevChar(line, c.charAt)
			if err != nil {
				if c.err {
					return // expected error
				}
				t.Fatalf("unexpected error: %s", err)
			}

			if c.err {
				t.Fatalf("Expected error: but succeeded with result = %s", result)
			}
			if c.expected != result {
				t.Errorf("%s", cmp.Diff(c.expected, result))
			}
		})
	}
}

func TestReadLineWithSkip(t *testing.T) {
	// 	original := `0123456789
	// 0123456789
	// 012三四五六七89
	// `
	cases := map[string]struct {
		original  string
		skipStart int
		skipEnd   int
		expected  string
		err       bool
	}{
		"0-0":                  {"0123456789\n" /********/, 0, 0, "0123456789\n", false},
		"2-3 Japanese":         {"012三四五六七八九\n" /**/, 2, 3, "01三四五六七八九\n", false},
		"2-4 Japanese":         {"012三四五六七八九\n" /**/, 2, 4, "01四五六七八九\n", false},
		"7-9 Japanese":         {"012三四五六七八九\n" /**/, 7, 9, "012三四五六九\n", false},
		"7-10 Japanese":        {"012三四五六七八九\n" /**/, 7, 10, "012三四五六\n", false},
		"ERROR: 7-11 Japanese": {"012三四五六七八九\n" /**/, 7, 11, "012三四五六\n", true}, // error expected, skipping the newline char is not allowed by this method
		"ERROR: 7-12 Japanese": {"012三四五六七八九\n" /**/, 7, 12, "012三四五六\n", true},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			line := []byte(c.original)
			result, err := readLineWithSkip(line, c.skipStart, c.skipEnd)
			if err != nil {
				if c.err {
					return // expected error
				}
				t.Fatalf("unexpected error: %s", err)
			}

			if c.err {
				t.Fatalf("Expected error: but succeeded with result = %s", result)
			}
			if c.expected != result {
				t.Errorf("%s", cmp.Diff(c.expected, result))
			}
		})
	}
}

func TestProcessLinesOnRange(t *testing.T) {
	// 	original := `0123456789
	// 0123456789
	// 012三四五六七89
	// `
	cases := map[string]struct {
		original string
		delRange Range
		expected string
		err      bool
	}{
		// "no deletion 1": {"0123456789", Range{Position{Line: 0, Character: 3}, Position{Line: 0, Character: 0}}, "0123456789", false},
		// "at the beginning 1": {"0123456789", Range{Position{Line: 0, Character: 0}, Position{Line: 0, Character: 1}}, "123456789", false},
		// "at the beginning 2": {"0123456789", Range{Position{Line: 0, Character: 0}, Position{Line: 0, Character: 2}}, "23456789", false},
		// "at the beginning 3": {"0123456789", Range{Position{Line: 0, Character: 0}, Position{Line: 0, Character: 3}}, "3456789", false},
		"no deletion 2": {"0123456789", Range{Position{Line: 0, Character: 3}, Position{Line: 0, Character: 3}}, "0123456789", false},

		// "at the beginning, end in newline":             {"0123456789\n" /******/, Position{Line: 0, Character: 0}, "Insert ", "Insert 0123456789\n", false},
		// "in the middle, 1":                             {"0123456789" /********/, Position{Line: 0, Character: 1}, " insert ", "0 insert 123456789", false},
		// "in the middle, 2":                             {"0123456789" /********/, Position{Line: 0, Character: 2}, " insert ", "01 insert 23456789", false},
		// "in the middle, 3":                             {"0123456789" /********/, Position{Line: 0, Character: 3}, " insert ", "012 insert 3456789", false},
		// "in the middle, Japanese":                      {"012三四五六七89" /****/, Position{Line: 0, Character: 3}, " 中間 ", "012 中間 三四五六七89", false},
		// "in the middle, English, end in newline":       {"0123456789\n" /******/, Position{Line: 0, Character: 3}, " insert ", "012 insert 3456789\n", false},
		// "in the middle, Japanese, end in newline":      {"012三四五六七89\n" /**/, Position{Line: 0, Character: 7}, " 中間 ", "012三四五六 中間 七89\n", false},
		// "close to the end, Japanese":                   {"012三四五六七89" /****/, Position{Line: 0, Character: 9}, " 中間 ", "012三四五六七8 中間 9", false},
		// "at the end, Japanese":                         {"012三四五六七89" /****/, Position{Line: 0, Character: 10}, " 最後", "012三四五六七89 最後", false},
		// "ERROR: at the end, Japanese, after end 1":     {"012三四五六七89" /****/, Position{Line: 0, Character: 11}, " 最後より後", "", true},
		// "ERROR: at the end, Japanese, after end 2":     {"012三四五六七89" /****/, Position{Line: 0, Character: 12}, " 最後より後", "", true},
		// "at the end, Japanese, end in newline":         {"012三四五六七89\n" /**/, Position{Line: 0, Character: 10}, " 最後", "012三四五六七89 最後\n", false},
		// "ERROR: at the end, Japanese, after newline 1": {"012三四五六七89\n" /**/, Position{Line: 0, Character: 11}, " 最後より後", "", true},
		// "ERROR: at the end, Japanese, after newline 2": {"012三四五六七89\n" /**/, Position{Line: 0, Character: 12}, " 最後より後", "", true},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var builder strings.Builder
			bufReader := bufio.NewReader(strings.NewReader(c.original))

			err := processLinesOnRange(bufReader, &builder, c.delRange)
			if err != nil {
				if c.err {
					return // expected error
				}
				t.Fatalf("unexpected error: %s", err)
			}

			result := builder.String()
			if c.err {
				t.Fatalf("Expected error: but succeeded with result = %s", result)
			}
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
