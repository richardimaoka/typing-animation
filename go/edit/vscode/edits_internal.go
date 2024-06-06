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
func insertLineByChar(currentPos Position, line string) ([]Edit, error) {
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
			return nil, errors.New("encountered new-line in the middle")
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
// line may only contain '\n' at the end, but not in the middle
//
// If line contains '\n' in the middle, this returns an error
// If line is empty, this should return the count of zero
func deleteLineByChar(startPos Position, line string) ([]Edit, error) {
	if len(line) == 0 {
		return nil, nil
	}

	edits := []Edit{}

	lineBytes := []byte(line)

	byteOffset := 0
	for c := 0; ; c++ {
		r, size := utf8.DecodeRune(lineBytes[byteOffset:])
		if r == '\n' && byteOffset != len(lineBytes)-1 {
			return nil, errors.New("encountered new-line in the middle")
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

// Return edis, split by word, to insert line from currentPos
// line may only contain '\n' at the end, but not in the middle
//
// If line contains '\n' in the middle, this returns an error
// If line is empty, this should return the count of zero
func insertLineByWord(currentPos Position, line string) ([]Edit, error) {
	if len(line) == 0 {
		return nil, nil
	}

	edits := []Edit{}

	// If line ends in '\n', add '\n' first otherwise the typing animation looks unnatural
	lineWithoutNL, hasNewLine := strings.CutSuffix(line, "\n")
	if hasNewLine {
		edits = append(edits, EditInsert{Position: currentPos, NewText: "\n"})
	}

	lineWords := strings.SplitAfter(lineWithoutNL, " ")

	pos := currentPos
	for _, word := range lineWords {
		if word == "" {
			continue
		}

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
// line may only contain '\n' at the end, but not in the middle
//
// If line contains '\n' in the middle, this returns an error
// If line is empty, this should return the count of zero
func deleteLineByWord(currentPos Position, line string) ([]Edit, error) {
	if len(line) == 0 {
		return nil, nil
	} else if line == "\n" {
		return []Edit{EditDelete{
			DeleteText:  "\n",
			DeleteRange: Range{Start: currentPos, End: Position{Line: currentPos.Line + 1, Character: 0}},
		}}, nil
	}

	edits := []Edit{}

	// Necessary to cut the last '\n', since countRunesInLine() expects no '\n' in the line
	lineWithoutNL, hasNewLine := strings.CutSuffix(line, "\n")
	lineWords := strings.SplitAfter(lineWithoutNL, " ")

	for _, word := range lineWords {
		if word == "" {
			continue
		}

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
	}

	if hasNewLine {
		edits = append(edits, EditDelete{
			DeleteText: "\n",
			DeleteRange: Range{
				Start: currentPos,
				End: Position{
					Line:      currentPos.Line + 1,
					Character: 0,
				},
			},
		})
	}

	return edits, nil
}

// Split
func splitInsertByLine(insert EditInsert) ([]Edit, error) {
	pos := insert.Position
	lines := strings.SplitAfter(insert.NewText, "\n")

	var edits []Edit
	for _, l := range lines {
		// if NewText ends in '\n', the last line is ""
		if l != "" {
			edits = append(edits, EditInsert{Position: pos, NewText: l})
			pos = Position{Line: pos.Line + 1, Character: 0}
		}
	}

	return edits, nil
}

func splitDeleteByLine(delete EditDelete) ([]Edit, error) {
	start := delete.DeleteRange.Start
	end := delete.DeleteRange.End

	lines := strings.SplitAfter(delete.DeleteText, "\n")

	var edits []Edit
	for lc := 0; lc < len(lines); lc++ {
		line := lines[lc]

		// if NewText ends in '\n', the last line is ""
		if line == "" {
			continue
		}

		var lineStart Position
		if lc == 0 {
			// line at the start position
			lineStart = Position{Line: start.Line + lc, Character: start.Character}
		} else {
			lineStart = Position{Line: start.Line + lc, Character: 0}
		}

		var lineEnd Position
		if lc == len(lines)-1 {
			// line at the end position
			lineEnd = Position{Line: start.Line + lc, Character: end.Character}
		} else {
			lineEnd = Position{Line: start.Line + lc + 1, Character: 0}
		}

		edits = append(edits, EditDelete{DeleteText: line, DeleteRange: Range{Start: lineStart, End: lineEnd}})
	}

	return edits, nil
}

func splitInsertByWord(insert EditInsert) ([]Edit, error) {
	pos := insert.Position
	lines := strings.SplitAfter(insert.NewText, "\n")

	var edits []Edit
	for _, l := range lines {
		lineEdits, err := insertLineByWord(pos, l)
		if err != nil {
			return nil, err
		}

		edits = append(edits, lineEdits...)
		pos = Position{Line: pos.Line + 1, Character: 0}
	}

	return edits, nil
}

func splitDeleteByWord(delete EditDelete) ([]Edit, error) {
	stratPos := delete.DeleteRange.Start
	lines := strings.SplitAfter(delete.DeleteText, "\n")

	var edits []Edit
	for _, l := range lines {
		lineEdits, err := deleteLineByWord(stratPos, l)
		if err != nil {
			return nil, err
		}

		edits = append(edits, lineEdits...)
	}

	return edits, nil
}

func splitInsertByChar(insert EditInsert) ([]Edit, error) {
	pos := insert.Position
	lines := strings.SplitAfter(insert.NewText, "\n")

	var edits []Edit
	for _, l := range lines {
		lineEdits, err := insertLineByChar(pos, l)
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
		lineEdits, err := deleteLineByChar(pos, l)
		if err != nil {
			return nil, err
		}

		edits = append(edits, lineEdits...)
		pos = Position{Line: pos.Line + 1, Character: 0}
	}

	return edits, nil
}
