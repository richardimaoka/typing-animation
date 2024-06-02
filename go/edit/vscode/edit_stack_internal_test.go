package vscode

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOffsetPosition(t *testing.T) {
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

		"delete single line": {initPos, Diff{DiffDelete, "a" /********/}, EditDelete{"a" /********/, Range{initPos, Position{Line: 3, Character: 11}}}, initPos, false},
		"delete multi line":  {initPos, Diff{DiffDelete, "aaaa\nb\nccc"}, EditDelete{"aaaa\nb\nccc", Range{initPos, Position{Line: 5, Character: 3}}}, initPos, false},
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
