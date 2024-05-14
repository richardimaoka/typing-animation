package vscode_test

import (
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/richardimaoka/typing-animation/go/edit/vscode"
)

func TestInsert(t *testing.T) {
	cases := map[string]struct {
		originalFile string
		pos          vscode.Position
		newText      string
		err          bool
	}{
		// "beginning":                {"testdata/insert/digits1.txt", vscode.Position{Line: 2, Character: 3}, "abc", false},
		// "in the middle":            {"testdata/insert/digits2.txt", vscode.Position{Line: 0, Character: 0}, "aaaa", false},
		// "in the middle Japanese":   {"testdata/insert/Japanese.txt", vscode.Position{Line: 3, Character: 6}, "すばらしい", false}, // And この文章のいくつかのpartは(char = 6 is「の」)
		// "in the middle, multiline": {"testdata/insert/digits2.txt", vscode.Position{Line: 0, Character: 0}, "aaaa", false},
		// "at the end":               {"testdata/insert/digits2.txt", vscode.Position{Line: 0, Character: 0}, "aaaa", false},
		// "at the end, newline":      {"testdata/insert/digits2.txt", vscode.Position{Line: 0, Character: 0}, "aaaa", false},

		"ERROR: negative line":       {"testdata/insert/1st_line_beginning.txt" /******/, vscode.Position{Line: -1, Character: 3}, "inserted ", true},
		"1st line, at the beginning": {"testdata/insert/1st_line_beginning.txt" /******/, vscode.Position{Line: 0, Character: 0}, "inserted ", false},
		"1st line, in the middle":    {"testdata/insert/1st_line_middle.txt" /*********/, vscode.Position{Line: 0, Character: 3}, " inserted ", false},
		"1st line, at the end":       {"testdata/insert/1st_line_end.txt" /************/, vscode.Position{Line: 0, Character: 10}, " at the end", false},
		"ERROR: 1st line, after end": {"testdata/insert/1st_line_end.txt" /************/, vscode.Position{Line: 2, Character: 11}, " at the end", true},
		// "middle line":                          {"testdata/insert/middle_line.txt" /*************/, vscode.Position{Line: 2, Character: 3}, "abc", false},
		// "middle line, Japanese":                {"testdata/insert/middle_line_Japanese.txt" /****/, vscode.Position{Line: 2, Character: 3}, "abc", false},
		// "middle line, at the end":              {"testdata/insert/middle_line_at_the_end.txt" /**/, vscode.Position{Line: 2, Character: 3}, "abc", false},
		// "middle line insert-multi-line text 1": {"testdata/insert/multi_line_insert1.txt" /******/, vscode.Position{Line: 2, Character: 3}, "abc", false},
		// "middle line insert-multi-line text 2": {"testdata/insert/multi_line_insert2.txt" /******/, vscode.Position{Line: 2, Character: 3}, "abc", false},
		// "last line":                            {"testdata/insert/1st_line2.txt" /***************/, vscode.Position{Line: 2, Character: 3}, "abc", false},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			// 1. Preparation
			//    Read input file
			expectedContents, err := os.ReadFile(c.originalFile)
			if err != nil {
				t.Fatal(err)
			}
			// Copy to temp file
			tempFile := strings.Replace(c.originalFile, ".txt", "_temp.txt", 1)
			err = os.WriteFile(tempFile, expectedContents, 0666)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				err = os.Remove(tempFile)
				if err != nil {
					t.Fatal(err)
				}
			}()

			// 2. Target operation
			//    Insert to temp file
			err = vscode.Insert(tempFile, c.pos, c.newText)
			if err != nil {
				if c.err {
					return // expected error
				}
				t.Fatal(err)
			}
			if c.err {
				t.Fatal("expected error but succeeded")
			}

			// 3. Check results
			//    Read from temp file
			resultedContents, err := os.ReadFile(tempFile)
			if err != nil {
				t.Fatal(err)
			}
			// Read from golden file
			goldenFile := strings.Replace(c.originalFile, ".txt", "_golden.txt", 1)
			expectedContents, err = os.ReadFile(goldenFile)
			if err != nil {
				t.Fatal(err)
			}
			// Comparison
			if string(expectedContents) != string(resultedContents) {
				t.Errorf("%s", cmp.Diff(string(expectedContents), string(resultedContents)))
			}
		})
	}
}

// func TestDelete(t *testing.T) {
// 	r := Range{Position{Line: 2, Character: 2}, Position{Line: 2, Character: 3}}
// 	err := Delete("testdata/inert/delete.txt", r)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
