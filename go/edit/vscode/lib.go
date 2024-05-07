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
// https://code.visualstudio.com/api/references/vscode-api#Position
type Position struct {
	Line      int //The zero-based line value.
	Character int //The zero-based character value.
}

type Range struct {
	Start Position
	End   Position
}

type FileHandler struct {
	file *os.File
}

func (p Position) Validate() bool {
	// Both is zero-based
	return p.Character >= 0 && p.Line >= 0
}

func (p Position) LessThanOrEqualTo(target Position) bool {
	if p.Line < target.Line {
		return true
	} else if p.Line == target.Line {
		return p.Character <= target.Character
	} else {
		// p.Line > target.Line
		return false
	}
}

func (r Range) Validate() bool {
	return r.Start.Validate() && r.End.Validate()
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
	// copy the position.Line up to the position.Character
	bytesAt := 0
	for j := 0; j <= position.Character; j++ {
		r, size := utf8.DecodeRune(line[bytesAt:])
		if size == 0 {
			return fmt.Errorf("argument position.Character = %d, but it cannnot be greater than %d, the number of characters at line = %d", position.Character, j, position.Line)
		}
		builder.WriteRune(r)
		bytesAt += size
	}
	// insert newText
	builder.WriteString(newText)
	// copy the rest of the line
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

func Delete(filename string, delRange Range) error {
	// 1. Validate arguments
	if delRange.Start.Line < 0 {
		return fmt.Errorf("argument delRange.Start.Line = %d, but it cannnot be a negative number", delRange.Start.Line)
	}
	if delRange.Start.Character < 0 {
		return fmt.Errorf("argument delRange.Start.Character = %d, but it cannnot be a negative number", delRange.Start.Character)
	}

	if delRange.End.Line < 0 {
		return fmt.Errorf("argument delRange.End.Line = %d, but it cannnot be a negative number", delRange.End.Line)
	}
	if delRange.End.Character < 0 {
		return fmt.Errorf("argument delRange.End.Character = %d, but it cannnot be a negative number", delRange.End.Character)
	}

	if delRange.Start.Line > delRange.End.Line {
		return fmt.Errorf("argument delRange.Start.Line = %d is greater than delRange.End.Line = %d", delRange.Start.Line, delRange.End.Line)
	}
	if delRange.Start.Line == delRange.End.Line && delRange.Start.Character > delRange.End.Character {
		return fmt.Errorf(
			"argument delRange.Start.Character = %d is greater than delRange.End.Character = %d on the same line %d, but ",
			delRange.Start.Character, delRange.End.Character, delRange.End.Line,
		)
	}

	// 2. Open file
	f, err := os.OpenFile(filename, syscall.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	bufReader := bufio.NewReader(f)
	var builder strings.Builder

	// 2. Copy up to the start line - 1
	for i := 0; i < delRange.Start.Line; i++ {
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

	// 3. Process the start line
	// read the position.Line
	line, isPrefix, err := bufReader.ReadLine()
	if err != nil {
		return err
	}
	if isPrefix {
		return fmt.Errorf("line is too long! This function does not handle `isPrefix = true` returned by bufio.ReadLine()")
	}
	// copy the position.Line up to the prev char of start
	bytesAt := 0
	for c := 0; c < delRange.Start.Character; c++ {
		r, size := utf8.DecodeRune(line[bytesAt:])
		if size == 0 {
			return fmt.Errorf("argument delRange.Start.Character = %d, but it cannnot be greater than %d, the number of characters at line = %d", delRange.Start.Character, c, delRange.Start.Line)
		}
		builder.WriteRune(r)
		bytesAt += size
	}

	// 4. If the start line == end line
	if delRange.Start.Line == delRange.End.Line {
		bytesAt := 0
		for j := delRange.Start.Character; j <= delRange.End.Character; j++ {
			_, size := utf8.DecodeRune(line[bytesAt:])
			if size == 0 {
				//TODO: change error
				return fmt.Errorf("argument delRange.End.Character = %d, but it cannnot be greater than %d, the number of characters at line = %d", delRange.End.Character, j, delRange.End.Line)
			}
			// incremental skip
			bytesAt += size
		}

		// copy the rest of the line
		builder.WriteString(string(line[bytesAt:]))
		builder.WriteString("\n")
	}

	// 5. Delete up to the end line - 1
	for i := delRange.Start.Line; i < delRange.End.Line; i++ {
		_, isPrefix, err := bufReader.ReadLine()
		if err != nil {
			return err
		}
		if isPrefix {
			return fmt.Errorf("line is too long! This function does not handle `isPrefix = true` returned by bufio.ReadLine()")
		}
	}

	// // 4. Copy the rest, until the end of file
	// for i := delRange.End.Line; ; i++ {
	// 	line, isPrefix, err := bufReader.ReadLine()
	// 	if err == io.EOF {
	// 		break
	// 	} else if err != nil {
	// 		return err
	// 	}
	// 	if isPrefix {
	// 		return fmt.Errorf(":ine is too long! This function does not handle `isPrefix = true` returned by bufio.ReadLine()")
	// 	}
	// 	builder.WriteString(string(line))
	// 	builder.WriteString("\n")
	// }

	// // 5. Write back to the file
	// f.WriteAt([]byte(builder.String()), 0)

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
