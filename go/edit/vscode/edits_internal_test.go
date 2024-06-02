package vscode

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAddCharByChar(t *testing.T) {
	cases := map[string]struct {
		newText    string
		currentPos Position
		expected   []Edit
		err        bool
	}{
		"abcde": {
			"abcde",
			Position{Line: 3, Character: 10},
			[]Edit{
				EditInsert{NewText: "a", Position: Position{Line: 3, Character: 10}},
				EditInsert{NewText: "b", Position: Position{Line: 3, Character: 11}},
				EditInsert{NewText: "c", Position: Position{Line: 3, Character: 12}},
				EditInsert{NewText: "d", Position: Position{Line: 3, Character: 13}},
				EditInsert{NewText: "e", Position: Position{Line: 3, Character: 14}},
			},
			false},
		"a b c": {
			"a b c",
			Position{Line: 3, Character: 10},
			[]Edit{
				EditInsert{NewText: "a", Position: Position{Line: 3, Character: 10}},
				EditInsert{NewText: " ", Position: Position{Line: 3, Character: 11}},
				EditInsert{NewText: "b", Position: Position{Line: 3, Character: 12}},
				EditInsert{NewText: " ", Position: Position{Line: 3, Character: 13}},
				EditInsert{NewText: "c", Position: Position{Line: 3, Character: 14}},
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
		expected   []Edit
		err        bool
	}{
		"this is a sentence.": {
			/**/ "this is a text.",
			Position{Line: 3, Character: 10},
			[]Edit{
				EditInsert{NewText: "this ", Position: Position{Line: 3, Character: 10}},
				EditInsert{NewText: "is ", Position: Position{Line: 3, Character: 15}},
				EditInsert{NewText: "a ", Position: Position{Line: 3, Character: 18}},
				EditInsert{NewText: "text.", Position: Position{Line: 3, Character: 20}},
			},
			false},
		"sentence with a wrong   spacing": {
			"a b c",
			Position{Line: 3, Character: 10},
			[]Edit{
				EditInsert{NewText: "a ", Position: Position{Line: 3, Character: 10}},
				EditInsert{NewText: "b ", Position: Position{Line: 3, Character: 12}},
				EditInsert{NewText: "c", Position: Position{Line: 3, Character: 14}},
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
