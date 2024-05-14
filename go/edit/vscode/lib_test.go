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
	}{
		"beginning":                {"testdata/insert/digits1.txt", vscode.Position{Line: 2, Character: 3}, "abc"},
		"in the middle":            {"testdata/insert/digits2.txt", vscode.Position{Line: 0, Character: 0}, "aaaa"},
		"in the middle Japanese":   {"testdata/insert/Japanese.txt", vscode.Position{Line: 3, Character: 6}, "すばらしい"}, // And この文章のいくつかのpartは(char = 6 is「の」)
		"in the middle, multiline": {"testdata/insert/digits2.txt", vscode.Position{Line: 0, Character: 0}, "aaaa"},
		"at the end":               {"testdata/insert/digits2.txt", vscode.Position{Line: 0, Character: 0}, "aaaa"},
		"at the end, newline":      {"testdata/insert/digits2.txt", vscode.Position{Line: 0, Character: 0}, "aaaa"},
		// "Japanese file n": {"testdata/inert/Japanese_n.txt", vscode.Position{Line: 3, Character: 3}, "abc"},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			// Read input file
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

			// Insert to temp file
			err = vscode.Insert(tempFile, c.pos, c.newText)
			if err != nil {
				t.Fatal(err)
			}

			// Read from temp file
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
