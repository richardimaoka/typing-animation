package vscode

import (
	"fmt"
	"os"
	"syscall"
)

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
