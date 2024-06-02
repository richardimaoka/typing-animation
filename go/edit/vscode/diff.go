package vscode

type DiffOperation int

const (
	// DiffDelete item represents a delete diff.
	DiffDelete DiffOperation = -1
	// DiffInsert item represents an insert diff.
	DiffInsert DiffOperation = 1
	// DiffEqual item represents an equal diff.
	DiffEqual DiffOperation = 0
)

type Diff struct {
	Type DiffOperation
	Text string
}
