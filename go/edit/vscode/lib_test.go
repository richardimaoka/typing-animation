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

		"ERROR: negative line":                  {"testdata/insert/1st_line_beginning.txt" /****/, vscode.Position{Line: -1, Character: 3}, "inserted ", true},
		"1st line, at the beginning":            {"testdata/insert/1st_line_beginning.txt" /****/, vscode.Position{Line: 0, Character: 0}, "inserted ", false},
		"1st line, in the middle":               {"testdata/insert/1st_line_middle.txt" /*******/, vscode.Position{Line: 0, Character: 3}, " inserted ", false},
		"1st line, at the end":                  {"testdata/insert/1st_line_end.txt" /**********/, vscode.Position{Line: 0, Character: 10}, " at the end", false},
		"ERROR: 1st line, after end":            {"testdata/insert/1st_line_end.txt" /**********/, vscode.Position{Line: 2, Character: 11}, " at the end", true},
		"middle line":                           {"testdata/insert/middle_line_English.txt" /***/, vscode.Position{Line: 2, Character: 4}, " inserted ", false},
		"middle line, Japanese":                 {"testdata/insert/middle_line_Japanese.txt" /**/, vscode.Position{Line: 2, Character: 4}, " inserted ", false},
		"middle line insert-multi-line text":    {"testdata/insert/middle_line_multilne.txt" /**/, vscode.Position{Line: 2, Character: 4}, " inserted \nnext line", false},
		"last line, at the end":                 {"testdata/insert/last_line_no_newline.txt" /**/, vscode.Position{Line: 5, Character: 10}, " inserted ", false},
		"ERROR: last line, at the end":          {"testdata/insert/last_line_no_newline.txt" /**/, vscode.Position{Line: 5, Character: 11}, " inserted ", true},
		"ERROR: after last line":                {"testdata/insert/last_line_no_newline.txt" /**/, vscode.Position{Line: 6, Character: 11}, " inserted ", true},
		"last line, at the end, newline":        {"testdata/insert/last_line_newline.txt" /*****/, vscode.Position{Line: 5, Character: 10}, " inserted ", false},
		"ERROR: last line, newline, at the end": {"testdata/insert/last_line_newline.txt" /*****/, vscode.Position{Line: 5, Character: 11}, " inserted ", true},
		"last line, true end, newline":          {"testdata/insert/last_line_trueend.txt" /*****/, vscode.Position{Line: 6, Character: 0}, " inserted ", false},
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

func TestDelete(t *testing.T) {
	cases := map[string]struct {
		originalFile string
		pos          vscode.Position
		newText      string
		err          bool
	}{
		// invalid range
		// beginning
		// at the end
		// multi-line delete
		// error delete after end

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
