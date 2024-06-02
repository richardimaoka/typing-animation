package vscode

type EditSplitStrategy int

const (
	SplitByLine EditSplitStrategy = 1
	SplitByWord EditSplitStrategy = 2
	SplitByChar EditSplitStrategy = 3
)

// Interface representing atomic edit
type Edit interface {
	IsEdit() bool
	Apply(filename string) error
	Split(strategy EditSplitStrategy) ([]Edit, error)
}

// Concrete edit types
type EditInsert struct {
	NewText  string
	Position Position
}

type EditDelete struct {
	DelRange Range
}

func (e EditInsert) IsEdit() bool {

	return true
}

func (e EditDelete) IsEdit() bool {
	return true
}

func (e EditInsert) Apply(filename string) error {
	return Insert(filename, e.Position, e.NewText)
}

func (e EditDelete) Apply(filename string) error {
	return Delete(filename, e.DelRange)
}

func (e EditInsert) Split(strategy EditSplitStrategy) ([]Edit, error) {
	switch strategy {
	case SplitByLine:
		return splitInsertByLine(e)
	case SplitByWord:
		return splitInsertByWord(e)
	case SplitByChar:
		return splitInsertByChar(e)
	default:
		return nil, nil
	}
}

func (e EditDelete) Split(strategy EditSplitStrategy) ([]Edit, error) {
	switch strategy {
	// case SplitByLine:
	// 	return splitInsertByLine(e)
	default:
		return nil, nil
	}
}