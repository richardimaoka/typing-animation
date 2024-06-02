package vscode

import (
	"errors"
	"strings"
	"unicode/utf8"
)

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
