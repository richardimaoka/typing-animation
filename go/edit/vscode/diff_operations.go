package vscode

type Operation int

const (
	// DiffDelete item represents a delete diff.
	DiffDelete Operation = -1
	// DiffInsert item represents an insert diff.
	DiffInsert Operation = 1
	// DiffEqual item represents an equal diff.
	DiffEqual Operation = 0
)

type Diff struct {
	Type Operation
	Text string
}
