package main_test

import "testing"

//https://github.com/richardimaoka/article-nextjs-todo-editable/commit/13036a97e1f232e1044bcf6fdfc2b961f5efb5c4
func Test(t *testing.T) {
	cases := []struct {
		inputFile  string
		goldenFile string
	}{
		{"", ""},
	}

	for _, c := range cases {
		t.Run(c.inputFile, func(t *testing.T) {

			//copy before.txt
			//change
			//write to tempfile
			//delete tempfile
		})
	}
}
