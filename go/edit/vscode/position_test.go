package vscode

import "testing"

func TestValidatePosition(t *testing.T) {
	cases := map[string]struct {
		success bool
		pos     Position
	}{
		"negative char": {false, Position{Character: -1, Line: 10}},
		"zero char":     {true, Position{Character: 0, Line: 10}},
		"negative line": {false, Position{Character: 1, Line: -10}},
		"zero line":     {true, Position{Character: 1, Line: 0}},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := c.pos.Validate()
			if err != nil && c.success {
				t.Fatalf("unexpected error %s", err)
			}
			if err == nil && !c.success {
				t.Fatalf("error expected but not encounterd")
			}
		})
	}
}

func TestLessThanOrEqualToPosition(t *testing.T) {
	cases := map[string]struct {
		expected bool
		p1       Position
		p2       Position
	}{
		"p1 line above p2 line":           {true, Position{Character: 5, Line: 10}, Position{Character: 1, Line: 12}},
		"p1 line below p2 line":           {false, Position{Character: 1, Line: 10}, Position{Character: 5, Line: 8}},
		"p1 left to p2 on the same line":  {true, Position{Character: 0, Line: 10}, Position{Character: 1, Line: 10}},
		"p1 right to p2 on the same line": {false, Position{Character: 5, Line: 10}, Position{Character: 1, Line: 10}},
		"p1 and p2 same position":         {true, Position{Character: 0, Line: 10}, Position{Character: 0, Line: 10}},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			result := c.p1.LessThanOrEqualTo(c.p2)
			if c.expected != result {
				t.Errorf("%t expected but result = %t", c.expected, result)
			}
		})
	}
}
