package vscode

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"syscall"
	"unicode/utf8"
)

// Same as VS Code extention API's Position
type Position struct {
	Character int //The zero-based character value.
	Line      int //The zero-based line value.
}

type Range struct {
	Start Position
	End   Position
}

type FileHandler struct {
	file *os.File
}

func NewFileHandler(filename string) (*FileHandler, error) {
	f, err := os.OpenFile(filename, syscall.O_RDWR, 0666)
	if err != nil {
		return nil, fmt.Errorf("NewFileHandler failed, %s", err)
	}

	return &FileHandler{
		file: f,
	}, nil
}

func (f *FileHandler) offset(position Position) (int, error) {
	if position.Line < 0 {
		return 0, fmt.Errorf("argument position.Line = %d, but it cannnot be a negative number", position.Line)
	}

	if position.Character < 0 {
		return 0, fmt.Errorf("argument position.Character = %d, but it cannnot be a negative number", position.Character)
	}

	r := bufio.NewReader(f.file)

	var totalBytes int
	var line []byte

	// read up to the prev line of position.Line
	for i := 0; i < position.Line; i++ {
		line, _, err := r.ReadLine()
		if err != nil {
			return 0, err
		}
		totalBytes += len(line)
	}

	// read the position.Line
	line, _, err := r.ReadLine()
	if err != nil {
		return 0, err
	}

	// read the position.Line up to the position.Character
	for i := 0; i <= position.Character; {
		_, size := utf8.DecodeRune(line[i:])
		if size == 0 {
			return 0, fmt.Errorf("argument position.Character = %d, but it cannnot be greater than %d, the number of characters at line = %d", position.Character, i, position.Line)
		}
		i += size
		totalBytes += size
	}

	_, err = f.file.Seek(0, 0)
	if err != nil {
		return 0, errors.New("failed to reset offset at the end")
	}

	return totalBytes, nil
}

func (f *FileHandler) Insert(position Position, newText string) error {
	offset, err := f.offset(position)
	if err != nil {
		return fmt.Errorf("Insert() failed, %s", err)
	}

	bytes := make([]byte, len(newText))
	n, err := f.file.ReadAt(bytes, int64(offset))
	if err != nil {
		return fmt.Errorf("Insert failed, %s", err)
	}
	if n != len(newText) {
		return fmt.Errorf("Insert failed, couldn't read %d bytes upon preparation", len(newText))
	}

	_, err = f.file.WriteAt([]byte(newText), int64(offset))
	if err != nil {
		return fmt.Errorf("Insert() failed, %s", err)
	}

	// f. filepointer
	return nil
}

func (f *FileHandler) Delete(position Position, newText string) error {
	return nil
}

func (f *FileHandler) Close() error {
	return f.file.Close()
}
