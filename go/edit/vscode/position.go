package vscode

import (
	"fmt"
	"strings"
)

// Same as VS Code extention API's Position
// https://code.visualstudio.com/api/references/vscode-api#Position
type Position struct {
	Line      int //The zero-based line value.
	Character int //The zero-based character value.
}

func (p Position) Validate() error {
	errors := []string{}

	// Both Line and Character is zero-based
	if p.Character < 0 {
		errors = append(errors, fmt.Sprintf("negative character = %d", p.Character))
	}

	if p.Line < 0 {
		errors = append(errors, fmt.Sprintf("negative line = %d", p.Line))
	}

	if len(errors) > 0 {
		return fmt.Errorf("invalid position, %s", strings.Join(errors, ", "))
	}

	return nil
}

func (p Position) LessThanOrEqualTo(target Position) bool {
	if p.Line < target.Line {
		return true
	} else if p.Line == target.Line {
		return p.Character <= target.Character
	} else {
		// p.Line > target.Line
		return false
	}
}
