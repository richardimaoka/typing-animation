package vscode

import (
	"testing"
)

func TestCountRunesInLine(t *testing.T) {
	cases := map[string]struct {
		line     string
		expected int
		err      bool
	}{
		"ASCII":                        {"0123456789", 10, false},
		"Japanese":                     {"012三四五六七八九", 10, false},
		"ERROR new line":               {"012三四五六七八九\n", 0, true},
		"ERROR new line in the middle": {"012三四五\n六七八九", 0, true},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			line := []byte(c.line)
			result, err := countRunesInLine(line)
			if err != nil {
				if c.err {
					return // expected error
				}
				t.Fatalf("unexpected error: %s", err)
			}

			if c.err {
				t.Fatalf("Expected error: but succeeded with result = %d", result)
			}
			if c.expected != result {
				t.Errorf("Result = %d is different from expected = %d", result, c.expected)
			}
		})
	}
}
