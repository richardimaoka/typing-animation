package vscode

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

// Count the number of runes in line
// line should not contain '\n'
//
// If line has '\n', this returns an error
// If line is empty, this should return the count of zero
func countRunesInLine(lineString string) (int, error) {
	if len(lineString) == 0 {
		return 0, nil
	}

	lineBytes := []byte(lineString)

	byteOffset := 0
	runeCount := 0
	for ; ; runeCount++ {
		r, size := utf8.DecodeRune(lineBytes[byteOffset:])
		if r == '\n' {
			return 0, errors.New("encountered new-line")
		}
		if r == utf8.RuneError {
			if size == 0 {
				// reached the end of line
				break
			} else {
				return 0, errors.New("encountered decoding error")
			}
		}

		byteOffset += size
	}

	return runeCount, nil
}

// Calculate the position offset by text
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
func offsetPosition(currentPos Position, text string) (Position, error) {
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

// func addWordByWord(newText string) []string {
// 	return nil
// }

// Return []EditInsert slice, split by char, to add line from currentPos
// line should not contain \'\n'
//
// If line contains '\n', this returns an error
// If line is empty, this should return the count of zero
func addCharByChar(currentPos Position, line string) ([]Edit, error) {
	if len(line) == 0 {
		return nil, nil
	}

	lineBytes := []byte(line)

	edits := []Edit{}
	byteOffset := 0
	for c := 0; ; c++ {
		r, size := utf8.DecodeRune(lineBytes[byteOffset:])
		if r == '\n' {
			return nil, errors.New("encountered new-line")
		}
		if r == utf8.RuneError {
			if size == 0 {
				// reached the end of line
				break
			} else {
				return nil, errors.New("encountered decoding error")
			}
		}

		edits = append(edits,
			EditInsert{
				NewText: string(r),
				Position: Position{
					Line:      currentPos.Line,
					Character: currentPos.Character + c},
			},
		)
		byteOffset += size
	}

	return edits, nil
}

// Return []EditInsert slice, split by word, to add line from currentPos
// line should not contain \'\n'
//
// If line contains '\n', this returns an error
// If line is empty, this should return the count of zero
func addWordByWord(currentPos Position, lineString string) ([]Edit, error) {
	pos := currentPos

	if len(lineString) == 0 {
		return nil, nil
	}

	lineWords := strings.SplitAfter(lineString, " ")

	edits := []Edit{}
	for _, word := range lineWords {
		edits = append(edits,
			EditInsert{
				NewText:  word,
				Position: pos,
			},
		)

		c, err := countRunesInLine(word)
		if err != nil {
			return nil, err
		}
		pos.Character = pos.Character + c // pos.Line remain same
	}

	return edits, nil
}

func diffToEdit(currentPos Position, diff Diff) (Edit, Position, error) {
	// regardless of diff type, edit end position is same
	editEndPos, err := offsetPosition(currentPos, diff.Text)
	if err != nil {
		return EditInsert{}, Position{}, err
	}

	switch diff.Type {
	case DiffInsert:
		return EditInsert{diff.Text, currentPos}, editEndPos, nil

	case DiffEqual:
		return nil, editEndPos, nil

	case DiffDelete:
		// Return the original currentPos, because the cursor doesn't move after deletion
		return EditDelete{Range{Start: currentPos, End: editEndPos}}, currentPos, nil

	default:
		return nil, Position{}, fmt.Errorf("diff type = %d is invalid", diff.Type)
	}
}

func splitInsertByLine(insert EditInsert) ([]Edit, error) {
	pos := insert.Position
	lines := strings.Split(insert.NewText, "\n")

	var edits []Edit
	for _, l := range lines {
		edits = append(edits, EditInsert{Position: pos, NewText: l})
		pos = Position{Line: pos.Line + 1, Character: 0}
	}

	return edits, nil
}

func splitInsertByWord(insert EditInsert) ([]Edit, error) {
	pos := insert.Position
	lines := strings.Split(insert.NewText, "\n")

	var edits []Edit
	for _, l := range lines {
		lineEdits, err := addWordByWord(pos, l)
		if err != nil {
			return nil, err
		}

		edits = append(edits, lineEdits...)
		pos = Position{Line: pos.Line + 1, Character: 0}
	}

	return edits, nil
}

func splitInsertByChar(insert EditInsert) ([]Edit, error) {
	pos := insert.Position
	lines := strings.Split(insert.NewText, "\n")

	var edits []Edit
	for _, l := range lines {
		lineEdits, err := addCharByChar(pos, l)
		if err != nil {
			return nil, err
		}

		edits = append(edits, lineEdits...)
		pos = Position{Line: pos.Line + 1, Character: 0}
	}

	return edits, nil
}
