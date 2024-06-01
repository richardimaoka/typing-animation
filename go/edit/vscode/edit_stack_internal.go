package vscode

import (
	"errors"
	"unicode/utf8"
)

// Count the number of runes in line
// line should not contain '\n'
//
// If line is empty, this should return zero count
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
