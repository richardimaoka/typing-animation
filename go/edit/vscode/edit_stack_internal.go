package vscode

import (
	"fmt"
	"strings"
)

// Calculate the edit text's range end position.
//
// Regardless of the edit type, either insert, equal nor deletion, the end position is same,
// because it is the end of the text range, *NOT* the cursor position after edit.
//
//  1. If single line
//     text = "abcde", then offset pos = Position{ Line: current line, Char: current char + 5 }
//     --------12345
//
//  2. If multi line
//     text = "\n\nabcde", then offset pos = Position{ Line: current line + 2, Char: 5 }
//     ------------12345
//
// newText may contain '\n'
func editRangeEnd(currentPos Position, text string) (Position, error) {
	lines := strings.Split(text, "\n")

	if len(lines) == 1 {
		// 1. If single line text
		line := lines[0]
		runeCount, err := countRunesInLine(line)
		if err != nil {
			return Position{}, fmt.Errorf("failed to calculate position offset, %s", err)
		}

		return Position{
			Line:      currentPos.Line,
			Character: currentPos.Character + runeCount,
		}, nil

	} else {
		// 2. If multi-line text
		lastLine := lines[len(lines)-1]

		runeCount, err := countRunesInLine(lastLine)
		if err != nil {
			return Position{}, fmt.Errorf("failed to calculate position offset, %s", err)
		}

		return Position{
			Line:      currentPos.Line + len(lines) - 1,
			Character: runeCount,
		}, nil
	}
}

func diffToEdit(currentPos Position, diff Diff) (Edit, Position, error) {
	// regardless of diff type, range end position is same
	rangeEndPos, err := editRangeEnd(currentPos, diff.Text)
	if err != nil {
		return EditInsert{}, Position{}, err
	}

	switch diff.Type {
	case DiffInsert:
		return EditInsert{diff.Text, currentPos}, rangeEndPos, nil

	case DiffEqual:
		return nil, rangeEndPos, nil

	case DiffDelete:
		// Return the original currentPos, because the cursor doesn't move after deletion
		return EditDelete{diff.Text, Range{Start: currentPos, End: rangeEndPos}}, currentPos, nil

	default:
		return nil, Position{}, fmt.Errorf("diff type = %d is invalid", diff.Type)
	}
}
