package main_test

import (
	"os"
	"strings"
	"testing"
)

//https://github.com/richardimaoka/article-nextjs-todo-editable/commit/13036a97e1f232e1044bcf6fdfc2b961f5efb5c4
func Test(t *testing.T) {
	cases := map[string]struct {
		originalFile string
	}{
		"before after": {"testdata/test.txt"},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			// 1. Preparation
			originalContents, err := os.ReadFile(c.originalFile)
			if err != nil {
				t.Fatalf("failed to read file = '%s'", c.originalFile)
			}
			// Copy to temp file
			tempFile := strings.Replace(c.originalFile, ".txt", "_temp.txt", 1)
			if err = os.WriteFile(tempFile, originalContents, 0666); err != nil {
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
			// err = vscode.Insert(tempFile, c.pos, c.newText)
			// if err != nil {
			// 	if c.err {
			// 		return // expected error
			// 	}
			// 	t.Fatal(err)
			// }
			// if c.err {
			// 	t.Fatal("expected error but succeeded")
			// }

			// 3. Check results
			//    Read from temp file
			// resultedContents, err := os.ReadFile(tempFile)
			// if err != nil {
			// 	t.Fatal(err)
			// }
			// // Read from golden file
			// goldenFile := strings.Replace(c.originalFile, ".txt", "_golden.txt", 1)
			// expectedContents, err := os.ReadFile(goldenFile)
			// if err != nil {
			// 	t.Fatal(err)
			// }
			// // Comparison
			// if string(expectedContents) != string(resultedContents) {
			// 	t.Errorf("%s", cmp.Diff(string(expectedContents), string(resultedContents)))
			// }
		})
	}
}
