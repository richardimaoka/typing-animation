package vscode

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
)

func Insert(reader io.Reader, position Position, newText string) (string, error) {
	fromReader := bufio.NewReader(reader)
	var toBuilder strings.Builder

	// 1. Copy up to position.Line-1
	if position.Line > 0 {
		if err := copyUpToLine(fromReader, &toBuilder, position.Line-1); err != nil {
			return "", err
		}
	}

	// 2. Process the position.Line
	if err := processLine(fromReader, &toBuilder, position, newText); err != nil {
		return "", err
	}

	// 3. Copy the rest, until the end of file
	if err := copyUntilEOF(fromReader, &toBuilder); err != nil {
		return "", err
	}

	return toBuilder.String(), nil
}

func Delete(reader io.Reader, delRange Range) (string, error) {
	fromReader := bufio.NewReader(reader)
	var toBuilder strings.Builder

	// 1. Copy up to position.Line-1
	if delRange.Start.Line > 0 {
		if err := copyUpToLine(fromReader, &toBuilder, delRange.Start.Line-1); err != nil {
			return "", err
		}
	}

	// 2. Process from start line to end line
	if err := processLinesOnRange(fromReader, &toBuilder, delRange); err != nil {
		return "", err
	}

	// 3. Copy the rest, until the end of file
	if err := copyUntilEOF(fromReader, &toBuilder); err != nil {
		return "", err
	}

	return toBuilder.String(), nil
}

func InsertInFile(filename string, position Position, newText string) error {
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

	// 3. Core logic
	result, err := Insert(file, position, newText)
	if err != nil {
		return fmt.Errorf("%s, %s", errorPrefix, err)
	}

	// 4. Write to file
	if err := writeFromBeginning(file, result); err != nil {
		return err
	}

	return nil
}

func DeleteInFile(filename string, delRange Range) error {
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

	// 3. Core logic
	result, err := Delete(file, delRange)
	if err != nil {
		return fmt.Errorf("%s, %s", errorPrefix, err)
	}

	// 4. Write to file
	//   os.Truncate() is needed because the deletion makes file shorter,
	//   otherwise, there will be a residual contents at the end of the file
	if err := os.Truncate(file.Name(), 0); err != nil {
		return err
	}
	if err := writeFromBeginning(file, result); err != nil {
		return err
	}

	return nil
}
