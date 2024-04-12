package vscode

import (
	"fmt"
	"os"
	"syscall"
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

func (f *FileHandler) Insert(position Position, newText string) error {
	// f. filepointer
	return nil
}

func (f *FileHandler) Delete(position Position, newText string) error {
	return nil
}

func (f *FileHandler) Close() error {
	return f.file.Close()
}
