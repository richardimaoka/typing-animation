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
		"new line at the end": {
			"abc\n",
			Position{Line: 3, Character: 10},
			[]Edit{
				EditInsert{NewText: "\n", Position: Position{Line: 3, Character: 10}},
				EditInsert{NewText: "a", Position: Position{Line: 3, Character: 10}},
				EditInsert{NewText: "b", Position: Position{Line: 3, Character: 11}},
				EditInsert{NewText: "c", Position: Position{Line: 3, Character: 12}},
			},
			false},
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

func TestDeleteCharByChar(t *testing.T) {
	cases := map[string]struct {
		deleteText string
		startPos   Position
		expected   []Edit
		err        bool
	}{
		"abcde": {
			"abcde",
			Position{Line: 3, Character: 10},
			[]Edit{
				EditDelete{DeleteText: "a", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 11}}},
				EditDelete{DeleteText: "b", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 11}}},
				EditDelete{DeleteText: "c", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 11}}},
				EditDelete{DeleteText: "d", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 11}}},
				EditDelete{DeleteText: "e", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 11}}},
			},
			false},
		"a b c": {
			"a b c",
			Position{Line: 3, Character: 10},
			[]Edit{
				EditDelete{DeleteText: "a", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 11}}},
				EditDelete{DeleteText: " ", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 11}}},
				EditDelete{DeleteText: "b", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 11}}},
				EditDelete{DeleteText: " ", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 11}}},
				EditDelete{DeleteText: "c", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 11}}},
			},
			false},
		"new line at the end": {
			"abc\n",
			Position{Line: 3, Character: 10},
			[]Edit{
				EditDelete{DeleteText: "a", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 11}}},
				EditDelete{DeleteText: "b", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 11}}},
				EditDelete{DeleteText: "c", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 11}}},
				EditDelete{DeleteText: "\n", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 11}}},
			},
			false},
		"ERROR: new line in the middle": {"0123456789\n012三四", Position{Line: 3, Character: 10}, nil, true},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := deleteCharByChar(c.startPos, c.deleteText)
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
			"text with a wrong   spacing",
			Position{Line: 3, Character: 10},
			[]Edit{
				EditInsert{NewText: "text ", Position: Position{Line: 3, Character: 10}},
				EditInsert{NewText: "with ", Position: Position{Line: 3, Character: 15}},
				EditInsert{NewText: "a ", Position: Position{Line: 3, Character: 20}},
				EditInsert{NewText: "wrong ", Position: Position{Line: 3, Character: 22}},
				EditInsert{NewText: " ", Position: Position{Line: 3, Character: 28}},
				EditInsert{NewText: " ", Position: Position{Line: 3, Character: 29}},
				EditInsert{NewText: "spacing", Position: Position{Line: 3, Character: 30}},
			},
			false},
		"new line at the end": {
			"this is a text.\n",
			Position{Line: 3, Character: 10},
			[]Edit{
				EditInsert{NewText: "\n", Position: Position{Line: 3, Character: 10}},
				EditInsert{NewText: "this ", Position: Position{Line: 3, Character: 10}},
				EditInsert{NewText: "is ", Position: Position{Line: 3, Character: 15}},
				EditInsert{NewText: "a ", Position: Position{Line: 3, Character: 18}},
				EditInsert{NewText: "text.", Position: Position{Line: 3, Character: 20}},
			},
			false},
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

func TestDeleteWordByWord(t *testing.T) {
	cases := map[string]struct {
		newText  string
		startPos Position
		expected []Edit
		err      bool
	}{
		"this is a text.": {
			/**/ "this is a text.",
			Position{Line: 3, Character: 10},
			[]Edit{
				EditDelete{DeleteText: "this ", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 15}}},
				EditDelete{DeleteText: "is ", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 13}}},
				EditDelete{DeleteText: "a ", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 12}}},
				EditDelete{DeleteText: "text.", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 15}}},
			},
			false},
		"text with a wrong   spacing": {
			"text with a wrong   spacing",
			Position{Line: 3, Character: 10},
			[]Edit{
				EditDelete{DeleteText: "text ", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 15}}},
				EditDelete{DeleteText: "with ", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 15}}},
				EditDelete{DeleteText: "a ", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 12}}},
				EditDelete{DeleteText: "wrong ", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 16}}},
				EditDelete{DeleteText: " ", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 11}}},
				EditDelete{DeleteText: " ", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 11}}},
				EditDelete{DeleteText: "spacing", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 17}}},
			},
			false},
		"new line at the end": {
			"this is a text.\n",
			Position{Line: 3, Character: 10},
			[]Edit{
				EditDelete{DeleteText: "this ", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 15}}},
				EditDelete{DeleteText: "is ", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 13}}},
				EditDelete{DeleteText: "a ", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 12}}},
				EditDelete{DeleteText: "text.", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 15}}},
				EditDelete{DeleteText: "\n", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 11}}},
			},
			false},
		"ERROR: new line in the middle": {"0123456789\n012三四", Position{Line: 3, Character: 10}, nil, true},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := deleteWordByWord(c.startPos, c.newText)
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

// func TestSplitInsertByWord(t *testing.T) {
// 	cases := map[string]struct {
// 		edit     EditInsert
// 		expected []Edit
// 		err      bool
// 	}{
// 		"this is a text.": {
// 			EditInsert{NewText: "abc def\nghi\njk", Position: Position{Line: 3, Character: 10}},
// 			[]Edit{
// 				// \n should be added first
// 				EditInsert{NewText: "\n", Position: Position{Line: 3, Character: 10}},
// 				EditInsert{NewText: "abc ", Position: Position{Line: 3, Character: 10}},
// 				EditInsert{NewText: "def", Position: Position{Line: 3, Character: 14}},
// 			},
// 			false,
// 		},
// 	}

// 	for name, c := range cases {
// 		t.Run(name, func(t *testing.T) {
// 			result, err := splitInsertByWord(c.edit)
// 			if err != nil {
// 				if c.err {
// 					return // expected error
// 				}
// 				t.Fatalf("unexpected error: %s", err)
// 			}

// 			if c.err {
// 				t.Fatalf("Expected error: but succeeded with result = %+v", result)
// 			}

// 			diff := cmp.Diff(c.expected, result)
// 			if len(diff) > 0 {
// 				t.Errorf("%s", diff)
// 			}
// 		})
// 	}
// }
