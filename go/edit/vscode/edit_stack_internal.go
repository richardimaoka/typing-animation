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
func countRunesInLine(line []byte) (int, error) {
	if len(line) == 0 {
		return 0, nil
	}

	byteOffset := 0
	runeCount := 0
	for ; ; runeCount++ {
		r, size := utf8.DecodeRune(line[byteOffset:])
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

func countRunesInLineS(line string) (int, error) {
	return countRunesInLine([]byte(line))
}

func positionAfterAdd(currentPos Position, newText string) (Position, error) {
	linesToAdd := strings.Split(newText, "\n")

	if len(linesToAdd) == 1 {
		line := linesToAdd[0]
		runeCount, err := countRunesInLineS(line)
		if err != nil {
			return Position{}, fmt.Errorf("failed to calculate new position after add %s", err)
		}

		return Position{
			Line:      currentPos.Line,
			Character: currentPos.Character + runeCount,
		}, nil

	} else {
		lastLine := linesToAdd[len(linesToAdd)-1]

		runeCount, err := countRunesInLineS(lastLine)
		if err != nil {
			return Position{}, fmt.Errorf("failed to calculate new position after add %s", err)
		}

		return Position{
			Line:      currentPos.Line + len(linesToAdd) - 1,
			Character: runeCount,
		}, nil
	}
}

// func addWordByWord(newText string) []string {
// 	return nil
// }

func addCharByChar(currentPos Position, lineString string) ([]EditInsert, error) {
	line := []byte(lineString)

	if len(line) == 0 {
		return []EditInsert{}, nil
	}

	edits := []EditInsert{}
	byteOffset := 0
	for c := 0; ; c++ {
		r, size := utf8.DecodeRune(line[byteOffset:])
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

// func addLineByLine(newText string) []string {
// 	return nil
// }
