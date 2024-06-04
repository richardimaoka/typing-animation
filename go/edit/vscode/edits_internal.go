package vscode

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// Return edits, split by char, to add a line from currentPos
// line may only contain '\n' at the end, but not in the middle
//
// If line contains '\n' in the middle, this returns an error
// If line is empty, this should return the count of zero
func addCharByChar(currentPos Position, line string) ([]Edit, error) {
	if len(line) == 0 {
		return nil, nil
	}

	edits := []Edit{}

	// If line ends in '\n', add '\n' first otherwise the typing animation looks unnatural
	lineWithoutNL, hasNewLine := strings.CutSuffix(line, "\n")
	if hasNewLine {
		edits = append(edits, EditInsert{Position: currentPos, NewText: "\n"})
	}

	lineBytes := []byte(lineWithoutNL)

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
func deleteCharByChar(startPos Position, line string) ([]Edit, error) {
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
				DeleteText: string(r),
				DeleteRange: Range{
					Start: startPos,
					End: Position{
						Line:      startPos.Line,
						Character: startPos.Character + 1,
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
				DeleteText: word,
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

func splitDeleteByWord(delete EditDelete) ([]Edit, error) {
	pos := delete.DeleteRange.Start
	lines := strings.Split(delete.DeleteText, "\n")

	var edits []Edit
	for _, l := range lines {
		lineEdits, err := deleteWordByWord(pos, l)
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
	lines := strings.SplitAfter(insert.NewText, "\n")

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
	pos := delete.DeleteRange.Start
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
