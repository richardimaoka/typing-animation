package vscode

// Same as VS Code extention API's Position
type Position struct {
	Character int //The zero-based character value.
	Line      int //The zero-based line value.
}

func Insert(filename string, position Position, newText string) error {
	return nil
}
