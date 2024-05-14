package vscode

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
	"unicode/utf8"
)

// Copy lines from fromReader up to uptoLine (zero-based), and copy the lines to toBuilder.
// from must be set to the initial position of the input, otherwise the behavior is not guaranteed.
// uptoLine must be zero or a positive number.
// If `uptoLine` is greater than the number of lines in fromReader, then error.
func copyUpToLine(fromReader *bufio.Reader, toBuilder *strings.Builder, uptoLine int) error {
	if uptoLine < 0 {
		return fmt.Errorf("uptoLine = %d is a negative number", uptoLine)
	}

	for i := 0; i <= uptoLine; /* since uptoLine is zero-based, upToLine is included */ i++ {
		line, err := fromReader.ReadBytes('\n')
		if err == io.EOF {
			if i == uptoLine {
				toBuilder.WriteString(string(line))
				break
			} else {
				return fmt.Errorf("trying to copy up to line = %d, but there are only %d lines", uptoLine, i)
			}
		} else if err != nil {
			return err
		}
		toBuilder.WriteString(string(line))
	}
	return nil
}

// Insert newText at chartAt on the line.
// If line has '\n', '\n' must be at the end of line, otherwise, behavior is not guaranteed
// If charAt is greater than the end of line, error is returned, otherwise, no error.
func insertInLine(charAt int, newText string, line []byte) (string, error) {
	var builder strings.Builder

	// copy the position.Line up to the position.Character
	byteOffset := 0
	for i := 0; i < charAt; i++ {
		r, size := utf8.DecodeRune(line[byteOffset:])
		if r == '\n' {
			return "", fmt.Errorf("trying to insert '%s' at char = %d, but encountered new-line at chart = %d", newText, charAt, i)
		}
		if r == utf8.RuneError {
			if size == 0 {
				return "", fmt.Errorf("trying to insert '%s' at char = %d, but reached the end of line at char = %d", newText, charAt, i)
			} else {
				return "", fmt.Errorf("trying to insert '%s' at char = %d, but encountered decoding error at char =  %d", newText, charAt, i)
			}
		}

		builder.WriteRune(r)
		byteOffset += size
	}

	// insert newText
	builder.WriteString(newText)

	// copy the rest of the line
	builder.WriteString(string(line[byteOffset:]))

	return builder.String(), nil
}

// Take the current line from fromReader, and insert newText to the line, then copy the updated line to toBuilder.
// It doesn't matter whether the line ends in '\n' or EOF
func processLine(fromReader *bufio.Reader, toBuilder *strings.Builder, pos Position, newText string) error {
	if err := pos.Validate(); err != nil {
		return fmt.Errorf("processing line = %d failed, %s", pos.Line, err)
	}

	// Firstly, read the line to process
	line, err := fromReader.ReadBytes('\n')
	if err == io.EOF {
		// Ignore EOF and continue - it doesn't matter whether the line ends in '\n' or EOF
	} else if err != nil {
		return fmt.Errorf("processing line = %d failed, %s", pos.Line, err)
	}

	// Secondly, write the udpted line
	//   Update line and assign it to a temp variable, updatedLine
	updatedLine, err := insertInLine(pos.Character, newText, line)
	if err != nil {
		return fmt.Errorf("processing line = %d failed, %s", pos.Line, err)
	}
	//   Copy totoBuilder from updatedLine
	_, err = toBuilder.WriteString(updatedLine)
	if err != nil {
		return fmt.Errorf("processing line = %d failed, %s", pos.Line, err)
	}

	return nil
}

func processRange(fromReader *bufio.Reader, toBuilder *strings.Builder, delRange Range) error {
	_ = copyUpToLine(fromReader, toBuilder, delRange.End.Line-delRange.Start.Line)
	return nil
}

// Read all the remaining lines from fromReader, then copy them to toBuilder
func copyUntilEOF(fromReader *bufio.Reader, toBuilder *strings.Builder) error {
	for {
		line, err := fromReader.ReadBytes('\n')
		if err == io.EOF {
			if len(line) != 0 {
				// If EOF, and line has something to copy
				toBuilder.WriteString(string(line))
			}
			// break out of the loop, as EOF
			break
		} else if err != nil {
			// For other errors, return error
			return fmt.Errorf("Insert() error, %s", err)
		}

		// not io.EOF yet, keep looping
		toBuilder.WriteString(string(line))
	}

	return nil
}

func writeFromBeginning(file *os.File, text string) error {
	n, err := file.WriteAt([]byte(text), 0)
	if n == 0 {
		return fmt.Errorf("no byte is written to file")
	}

	// return err on error, nil on success
	return err
}

func insertInternal(rwFile *os.File, position Position, newText string) error {
	fromReader := bufio.NewReader(rwFile)
	var toBuilder strings.Builder

	// 1. Copy up to position.Line-1
	if position.Line > 0 {
		if err := copyUpToLine(fromReader, &toBuilder, position.Line-1); err != nil {
			return err
		}
	}

	// 2. Process the position.Line
	if err := processLine(fromReader, &toBuilder, position, newText); err != nil {
		return err
	}

	// 3. Copy the rest, until the end of file
	if err := copyUntilEOF(fromReader, &toBuilder); err != nil {
		return err
	}

	// 4. Write to file
	if err := writeFromBeginning(rwFile, toBuilder.String()); err != nil {
		return err
	}

	return nil
}

func Insert(filename string, position Position, newText string) error {
	errorPrefix := fmt.Errorf("Insert() error in file = '%s'", filename)
	// 1. Validate arguments
	if err := position.Validate(); err != nil {
		return fmt.Errorf("%s, %s", errorPrefix, err)
	}

	// 2. Open file
	file, err := os.OpenFile(filename, syscall.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("%s, %s", errorPrefix, err)
	}
	defer file.Close()

	// 3. Internal logic
	if err := insertInternal(file, position, newText); err != nil {
		return fmt.Errorf("%s, %s", errorPrefix, err)
	}

	return nil
}

func deleteInternal(rwFile *os.File, delRange Range) error {
	fromReader := bufio.NewReader(rwFile)
	var toBuilder strings.Builder

	// 1. Copy up to position.Line-1
	if delRange.Start.Line > 0 {
		if err := copyUpToLine(fromReader, &toBuilder, delRange.Start.Line-1); err != nil {
			return err
		}
	}

	// 2. Process the position.Line
	if err := processRange(fromReader, &toBuilder, delRange); err != nil {
		return err
	}

	// 3. Copy the rest, until the end of file
	if err := copyUntilEOF(fromReader, &toBuilder); err != nil {
		return err
	}

	// 4. Write to file
	if err := writeFromBeginning(rwFile, toBuilder.String()); err != nil {
		return err
	}

	return nil
}

func Delete(filename string, delRange Range) error {
	errorPrefix := fmt.Errorf("Delete() error in file = '%s'", filename)

	// 1. Validate arguments
	if err := delRange.Validate(); err != nil {
		return fmt.Errorf("%s, %s", errorPrefix, err)
	}

	// 2. Open file
	file, err := os.OpenFile(filename, syscall.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("%s, %s", errorPrefix, err)
	}
	defer file.Close()

	// 3. Internal logic
	if err := deleteInternal(file, delRange); err != nil {
		return fmt.Errorf("%s, %s", errorPrefix, err)
	}

	return nil
}
