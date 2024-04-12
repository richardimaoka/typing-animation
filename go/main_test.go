package main_test

import "testing"

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
