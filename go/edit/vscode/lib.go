package vscode

import (
	"bufio"
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
	file       *os.File
	readWriter *bufio.ReadWriter
}

func NewFileHandler(filename string) (*FileHandler, error) {
	f, err := os.OpenFile(filename, syscall.O_RDWR, 0666)
	if err != nil {
		return nil, fmt.Errorf("NewFileHandler failed, %s", err)
	}

	rw := bufio.NewReadWriter(bufio.NewReader(f), bufio.NewWriter(f))

	return &FileHandler{
		file:       f,
		readWriter: rw,
	}, nil
}

func (f *FileHandler) byteIndexAt(position Position) (int, error) {
	if position.Line < 0 {
		return 0, fmt.Errorf("argument position.Line = %d, but it cannnot be a negative number", position.Line)
	}

	if position.Character < 0 {
		return 0, fmt.Errorf("argument position.Character = %d, but it cannnot be a negative number", position.Character)
	}

	var totalBytes int
	var line []byte

	// read up to the prev line of position.Line
	for i := 0; i < position.Line; i++ {
		line, _, err := f.readWriter.ReadLine()
		if err != nil {
			return 0, err
		}
		totalBytes += len(line)
		fmt.Println("aaa")
	}

	// read the position.Line
	line, _, err := f.readWriter.ReadLine()
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

	return totalBytes, nil
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
