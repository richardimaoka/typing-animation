package vscode

import (
	"errors"
	"fmt"
	"strings"
)

// Same as VS Code extention API's Position
// https://code.visualstudio.com/api/references/vscode-api#Range
type Range struct {
	Start Position
	End   Position
}

func (r Range) Validate() error {
	errs := []string{}

	// Both Line and Character is zero-based
	if err := r.Start.Validate(); err != nil {
		errs = append(errs, fmt.Sprintf("range start error, %s", err))
	}

	if err := r.End.Validate(); err != nil {
		errs = append(errs, fmt.Sprintf("range end error, %s", err))
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	if !r.Start.LessThanOrEqualTo(r.End) {
		return fmt.Errorf("invalid range, start %+v > end %+v", r.Start, r.End)
	}

	return nil
}
