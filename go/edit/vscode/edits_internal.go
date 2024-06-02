package vscode

import (
	"errors"
	"strings"
	"unicode/utf8"
)

func (e EditDelete) StartPos() Position {
	return e.DeleteRange.Start
}

// Return edits, split by char, to add line from currentPos
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

// Return edits, split by char, to delete line from currentPos
// line should not contain \'\n'
//
// If line contains '\n', this returns an error
// If line is empty, this should return the count of zero
func deleteCharByChar(currentPos Position, line string) ([]Edit, error) {
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
			EditDelete{
				DeleteRange: Range{
					Start: currentPos,
					End: Position{
						Line:      currentPos.Line,
						Character: currentPos.Character + 1,
					},
				},
			},
		)
		byteOffset += size
	}

	return edits, nil
}

// Return edis, split by word, to add line from currentPos
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

// Return edis, split by word, to delete line from currentPos
// line should not contain \'\n'
//
// If line contains '\n', this returns an error
// If line is empty, this should return the count of zero
func deleteWordByWord(currentPos Position, lineString string) ([]Edit, error) {
	pos := currentPos

	if len(lineString) == 0 {
		return nil, nil
	}

	lineWords := strings.SplitAfter(lineString, " ")

	edits := []Edit{}
	for _, word := range lineWords {
		c, err := countRunesInLine(word)
		if err != nil {
			return nil, err
		}

		edits = append(edits,
			EditDelete{
				DeleteRange: Range{
					Start: currentPos,
					End: Position{
						Line:      currentPos.Line,
						Character: currentPos.Character + c,
					},
				},
			},
		)

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

func splitDeleteByLine(delete EditDelete) ([]Edit, error) {
	start := delete.DeleteRange.Start
	end := delete.DeleteRange.End

	var edits []Edit
	for l := start.Line; l <= end.Line; l++ {
		var lineStart Position
		if l == start.Line {
			lineStart = Position{Line: l, Character: start.Character}
		} else {
			lineStart = Position{Line: l, Character: 0}
		}

		var lineEnd Position
		if l == end.Line {
			lineEnd = Position{Line: l, Character: end.Character}
		} else {
			lineEnd = Position{Line: l + 1, Character: 0}
		}

		edits = append(edits, EditDelete{DeleteRange: Range{Start: lineStart, End: lineEnd}})
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

// func splitDeleteByWord(delete EditDelete) ([]Edit, error) {
// 	pos := delete.Position
// 	lines := strings.Split(delete.NewText, "\n")

// 	var edits []Edit
// 	for _, l := range lines {
// 		lineEdits, err := addWordByWord(pos, l)
// 		if err != nil {
// 			return nil, err
// 		}

// 		edits = append(edits, lineEdits...)
// 		pos = Position{Line: pos.Line + 1, Character: 0}
// 	}

// 	return edits, nil
// }

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

func splitDeleteByChar(delete EditDelete) ([]Edit, error) {
	pos := delete.StartPos()
	lines := strings.Split(delete.DeleteText, "\n")

	var edits []Edit
	for _, l := range lines {
		lineEdits, err := deleteCharByChar(pos, l)
		if err != nil {
			return nil, err
		}

		edits = append(edits, lineEdits...)
		pos = Position{Line: pos.Line + 1, Character: 0}
	}

	return edits, nil
}
