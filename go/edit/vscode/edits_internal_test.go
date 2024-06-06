package vscode

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInsertLineByChar(t *testing.T) {
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
			result, err := insertLineByChar(c.currentPos, c.newText)
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

func TestDeleteLineByChar(t *testing.T) {
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
			result, err := deleteLineByChar(c.startPos, c.deleteText)
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

func TestInsertLineByWord(t *testing.T) {
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
			result, err := insertLineByWord(c.currentPos, c.newText)
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

func TestDeleteLineByWord(t *testing.T) {
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
			result, err := deleteLineByWord(c.startPos, c.newText)
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

func TestSplitInsertByLine(t *testing.T) {
	cases := map[string]struct {
		edit     EditInsert
		expected []Edit
		err      bool
	}{
		"single line": {
			EditInsert{NewText: "123456789", Position: Position{Line: 3, Character: 10}},
			[]Edit{
				EditInsert{NewText: "123456789", Position: Position{Line: 3, Character: 10}},
			},
			false,
		},
		`single line, ends in \n`: {
			EditInsert{NewText: "123456789\n", Position: Position{Line: 3, Character: 10}},
			[]Edit{
				EditInsert{NewText: "123456789\n", Position: Position{Line: 3, Character: 10}},
			},
			false,
		},
		`consecutive\n`: {
			EditInsert{NewText: "123456\n\n\n789\n", Position: Position{Line: 3, Character: 10}},
			[]Edit{
				EditInsert{NewText: "123456\n", Position: Position{Line: 3, Character: 10}},
				EditInsert{NewText: "\n", Position: Position{Line: 4, Character: 0}},
				EditInsert{NewText: "\n", Position: Position{Line: 5, Character: 0}},
				EditInsert{NewText: "789\n", Position: Position{Line: 6, Character: 0}},
			},
			false,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := splitInsertByLine(c.edit)
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

func TestSplitDeletetByLine(t *testing.T) {
	cases := map[string]struct {
		edit     EditDelete
		expected []Edit
		err      bool
	}{
		"single line": {
			EditDelete{DeleteText: "123456789", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 20}}},
			[]Edit{
				EditDelete{DeleteText: "123456789", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 20}}},
			},
			false,
		},
		`single line, ends in \n`: {
			EditDelete{DeleteText: "123456789\n", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 4, Character: 0}}},
			[]Edit{
				EditDelete{DeleteText: "123456789\n", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 4, Character: 0}}},
			},
			false,
		},
		`consecutive\n`: {
			EditDelete{DeleteText: "123456\n\n\n789\n", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 7, Character: 0}}},
			[]Edit{
				EditDelete{DeleteText: "123456\n", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 4, Character: 0}}},
				EditDelete{DeleteText: "\n", DeleteRange: Range{Start: Position{Line: 4, Character: 0}, End: Position{Line: 5, Character: 0}}},
				EditDelete{DeleteText: "\n", DeleteRange: Range{Start: Position{Line: 5, Character: 0}, End: Position{Line: 6, Character: 0}}},
				EditDelete{DeleteText: "789\n", DeleteRange: Range{Start: Position{Line: 6, Character: 0}, End: Position{Line: 7, Character: 0}}},
			},
			false,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := splitDeleteByLine(c.edit)
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

func TestSplitDeletetByWord(t *testing.T) {
	cases := map[string]struct {
		edit     EditDelete
		expected []Edit
		err      bool
	}{
		"single line, single word": {
			EditDelete{DeleteText: "123456789", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 19}}},
			[]Edit{
				EditDelete{DeleteText: "123456789", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 19}}},
			},
			false,
		},
		"single line, multi words": {
			EditDelete{DeleteText: "123 456 789", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 21}}},
			[]Edit{
				EditDelete{DeleteText: "123 ", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 14}}},
				EditDelete{DeleteText: "456 ", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 14}}},
				EditDelete{DeleteText: "789", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 13}}},
			},
			false,
		},
		// `single line, ends in \n`: {
		// 	EditDelete{DeleteText: "123 456 789\n", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 4, Character: 0}}},
		// 	[]Edit{
		// 		EditDelete{DeleteText: "123 ", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 14}}},
		// 		EditDelete{DeleteText: "456 ", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 3, Character: 14}}},
		// 		EditDelete{DeleteText: "789\n", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 4, Character: 0}}},
		// 	},
		// 	false,
		// },
		// `consecutive\n`: {
		// 	EditDelete{DeleteText: "123456\n\n\n789\n", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 7, Character: 0}}},
		// 	[]Edit{
		// 		EditDelete{DeleteText: "123456\n", DeleteRange: Range{Start: Position{Line: 3, Character: 10}, End: Position{Line: 4, Character: 0}}},
		// 		EditDelete{DeleteText: "\n", DeleteRange: Range{Start: Position{Line: 4, Character: 0}, End: Position{Line: 5, Character: 0}}},
		// 		EditDelete{DeleteText: "\n", DeleteRange: Range{Start: Position{Line: 5, Character: 0}, End: Position{Line: 6, Character: 0}}},
		// 		EditDelete{DeleteText: "789\n", DeleteRange: Range{Start: Position{Line: 6, Character: 0}, End: Position{Line: 7, Character: 0}}},
		// 	},
		// 	false,
		// },
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := splitDeleteByWord(c.edit)
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
