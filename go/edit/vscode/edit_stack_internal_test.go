package vscode

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCountRunesInLine(t *testing.T) {
	cases := map[string]struct {
		line     string
		expected int
		err      bool
	}{
		"zero count":                   {"", 0, false},
		"ASCII":                        {"0123456789", 10, false},
		"Japanese":                     {"012三四五六七八九", 10, false},
		"ERROR new line":               {"012三四五六七八九\n", 0, true},
		"ERROR new line in the middle": {"012三四五\n六七八九", 0, true},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := countRunesInLine(c.line)
			if err != nil {
				if c.err {
					return // expected error
				}
				t.Fatalf("unexpected error: %s", err)
			}

			if c.err {
				t.Fatalf("Expected error: but succeeded with result = %d", result)
			}
			if c.expected != result {
				t.Errorf("Result = %d is different from expected = %d", result, c.expected)
			}
		})
	}
}

func TestNewPositionAfterAdd(t *testing.T) {
	cases := map[string]struct {
		currentPos Position
		newText    string
		expected   Position
		err        bool
	}{
		"empty":                       {Position{Line: 3, Character: 10}, "" /*********************/, Position{Line: 3, Character: 10}, false},
		"single line":                 {Position{Line: 3, Character: 10}, "0123456789" /***********/, Position{Line: 3, Character: 10 + 10}, false},
		"multi lines":                 {Position{Line: 3, Character: 10}, "0123456789\n012三四" /**/, Position{Line: 4, Character: 5}, false},
		"multi lines end in new-line": {Position{Line: 3, Character: 10}, "0123456789\n" /*********/, Position{Line: 4, Character: 0}, false},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := offsetPosition(c.currentPos, c.newText)
			if err != nil {
				if c.err {
					return // expected error
				}
				t.Fatalf("unexpected error: %s", err)
			}

			if c.err {
				t.Fatalf("Expected error: but succeeded with result = %+v", result)
			}
			if c.expected != result {
				t.Errorf("%s", cmp.Diff(c.expected, result))
			}
		})
	}
}

func TestAddCharByChar(t *testing.T) {
	cases := map[string]struct {
		newText    string
		currentPos Position
		expected   []EditInsert
		err        bool
	}{
		"abcde": {
			"abcde",
			Position{Line: 3, Character: 10},
			[]EditInsert{
				{NewText: "a", Position: Position{Line: 3, Character: 10}},
				{NewText: "b", Position: Position{Line: 3, Character: 11}},
				{NewText: "c", Position: Position{Line: 3, Character: 12}},
				{NewText: "d", Position: Position{Line: 3, Character: 13}},
				{NewText: "e", Position: Position{Line: 3, Character: 14}},
			},
			false},
		"a b c": {
			"a b c",
			Position{Line: 3, Character: 10},
			[]EditInsert{
				{NewText: "a", Position: Position{Line: 3, Character: 10}},
				{NewText: " ", Position: Position{Line: 3, Character: 11}},
				{NewText: "b", Position: Position{Line: 3, Character: 12}},
				{NewText: " ", Position: Position{Line: 3, Character: 13}},
				{NewText: "c", Position: Position{Line: 3, Character: 14}},
			},
			false},
		"ERROR: new line at the end":    {"0123456789\n", Position{Line: 3, Character: 10}, nil, true},
		"ERROR: new line in the middle": {"0123456789\n012三四", Position{Line: 3, Character: 10}, nil, true},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := addCharByChar(c.currentPos, c.newText)
			if err != nil {
				if c.err {
					return // expected error
				}
				t.Fatalf("unexpected error: %s", err)
			}

			if c.err {
				t.Fatalf("Expected error: but succeeded with result = %+v", result)
			}

			diff := cmp.Diff(c.expected, result)
			if len(diff) > 0 {
				t.Errorf("%s", diff)
			}
		})
	}
}

func TestAddWordByWord(t *testing.T) {
	cases := map[string]struct {
		newText    string
		currentPos Position
		expected   []EditInsert
		err        bool
	}{
		"this is a sentence.": {
			/**/ "this is a text.",
			Position{Line: 3, Character: 10},
			[]EditInsert{
				{NewText: "this ", Position: Position{Line: 3, Character: 10}},
				{NewText: "is ", Position: Position{Line: 3, Character: 15}},
				{NewText: "a ", Position: Position{Line: 3, Character: 18}},
				{NewText: "text.", Position: Position{Line: 3, Character: 20}},
			},
			false},
		"sentence with a wrong   spacing": {
			"a b c",
			Position{Line: 3, Character: 10},
			[]EditInsert{
				{NewText: "a ", Position: Position{Line: 3, Character: 10}},
				{NewText: "b ", Position: Position{Line: 3, Character: 12}},
				{NewText: "c", Position: Position{Line: 3, Character: 14}},
			},
			false},
		"ERROR: new line at the end":    {"0123456789\n", Position{Line: 3, Character: 10}, nil, true},
		"ERROR: new line in the middle": {"0123456789\n012三四", Position{Line: 3, Character: 10}, nil, true},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := addWordByWord(c.currentPos, c.newText)
			if err != nil {
				if c.err {
					return // expected error
				}
				t.Fatalf("unexpected error: %s", err)
			}

			if c.err {
				t.Fatalf("Expected error: but succeeded with result = %+v", result)
			}

			diff := cmp.Diff(c.expected, result)
			if len(diff) > 0 {
				t.Errorf("%s", diff)
			}
		})
	}
}

func TestDiffToEdit(t *testing.T) {
	initPos := Position{Line: 3, Character: 10}
	cases := map[string]struct {
		currentPos   Position
		diff         Diff
		expectedEdit Edit
		expectedPos  Position
		err          bool
	}{
		"insert single line": {initPos, Diff{DiffInsert, "a" /********/}, EditInsert{"a" /********/, initPos}, Position{Line: 3, Character: 11}, false},
		"insert multi line":  {initPos, Diff{DiffInsert, "aaaa\nb\nccc"}, EditInsert{"aaaa\nb\nccc", initPos}, Position{Line: 5, Character: 3}, false},

		"equal single line": {initPos, Diff{DiffEqual, "a" /********/}, nil, Position{Line: 3, Character: 11}, false},
		"equal multi line":  {initPos, Diff{DiffEqual, "aaaa\nb\nccc"}, nil, Position{Line: 5, Character: 3}, false},

		"delete single line": {initPos, Diff{DiffDelete, "a" /********/}, EditDelete{Range{initPos, Position{Line: 3, Character: 11}}}, initPos, false},
		"delete multi line":  {initPos, Diff{DiffDelete, "aaaa\nb\nccc"}, EditDelete{Range{initPos, Position{Line: 5, Character: 3}}}, initPos, false},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			resultEdit, resultPos, err := diffToEdit(c.currentPos, c.diff)
			if err != nil {
				if c.err {
					return // expected error
				}
				t.Fatalf("unexpected error: %s", err)
			}

			if c.err {
				t.Fatalf("Expected error: but succeeded with result = %+v, %+v", resultPos, resultEdit)
			}

			diffPos := cmp.Diff(c.expectedPos, resultPos)
			if len(diffPos) > 0 {
				t.Errorf("%s", diffPos)
			}

			diffEdit := cmp.Diff(c.expectedEdit, resultEdit)
			if len(diffEdit) > 0 {
				t.Errorf("%s", diffEdit)
			}
		})
	}
}
