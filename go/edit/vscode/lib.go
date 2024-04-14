package vscode

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
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

func Insert(filename string, position Position, newText string) error {
	// 1. Validate arguments
	if position.Line < 0 {
		return fmt.Errorf("argument position.Line = %d, but it cannnot be a negative number", position.Line)
	}

	if position.Character < 0 {
		return fmt.Errorf("argument position.Character = %d, but it cannnot be a negative number", position.Character)
	}

	// 2. Open file
	f, err := os.OpenFile(filename, syscall.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	bufReader := bufio.NewReader(f)
	var builder strings.Builder

	// 2. Copy up to the prev line of position.Line
	for i := 0; i < position.Line; i++ {
		line, isPrefix, err := bufReader.ReadLine()
		if err != nil {
			return err
		}
		if isPrefix {
			return fmt.Errorf(":ine is too long! This function does not handle `isPrefix = true` returned by bufio.ReadLine()")
		}
		builder.WriteString(string(line))
		builder.WriteString("\n")
	}

	// 3. Process the position.Line
	// read the position.Line
	line, isPrefix, err := bufReader.ReadLine()
	if err != nil {
		return err
	}
	if isPrefix {
		return fmt.Errorf("line is too long! This function does not handle `isPrefix = true` returned by bufio.ReadLine()")
	}
	// read the position.Line up to the position.Character
	bytesAt := 0
	for i := 0; i <= position.Character; i++ {
		r, size := utf8.DecodeRune(line[bytesAt:])
		if size == 0 {
			return fmt.Errorf("argument position.Character = %d, but it cannnot be greater than %d, the number of characters at line = %d", position.Character, i, position.Line)
		}
		builder.WriteRune(r)
		bytesAt += size
	}
	// insert newText
	builder.WriteString(newText)
	builder.WriteString(string(line[bytesAt:]))
	builder.WriteString("\n")

	// 4. Copy the rest, until the end of file
	for i := 0; i < position.Line; i++ {
		line, isPrefix, err := bufReader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if isPrefix {
			return fmt.Errorf(":ine is too long! This function does not handle `isPrefix = true` returned by bufio.ReadLine()")
		}
		if len(line) == 0 {
			break
		}
		builder.WriteString(string(line))
		builder.WriteString("\n")
	}

	// 5. Write back to the file
	f.WriteAt([]byte(builder.String()), 0)

	return nil
}

func (f *FileHandler) Insert2(position Position, newText string) error {
	offset, err := f.offset(position)
	if err != nil {
		return fmt.Errorf("Insert() failed, %s", err)
	}

	oldText := make([]byte, len(newText))
	n, err := f.file.ReadAt(oldText, int64(offset))
	if err != nil {
		return fmt.Errorf("Insert failed, %s", err)
	}
	if n != len(newText) {
		return fmt.Errorf("Insert failed, couldn't read %d bytes upon preparation", len(newText))
	}

	updatedText := string(oldText) + newText
	_, err = f.file.WriteAt([]byte(updatedText), int64(offset))
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
