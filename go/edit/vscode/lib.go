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
	result, err := insertInternal(file, position, newText)
	if err != nil {
		return fmt.Errorf("%s, %s", errorPrefix, err)
	}

	// 4. Write to file
	if err := writeFromBeginning(file, result); err != nil {
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
	result, err := deleteInternal(file, delRange)
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
