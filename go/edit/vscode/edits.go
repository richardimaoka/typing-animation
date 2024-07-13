package vscode

type SplitStrategy int

const (
	SplitByLine SplitStrategy = 1
	SplitByWord SplitStrategy = 2
	SplitByChar SplitStrategy = 3
)

// Interface representing atomic edit
type Edit interface {
	Apply(filename string) error
	Split(strategy SplitStrategy) ([]Edit, error)
}

// Concrete edit types
type EditInsert struct {
	NewText  string
	Position Position
}

type EditDelete struct {
	DeleteText  string // DeleteText is necessary for word-by-word split, and char-by-char split
	DeleteRange Range

	// Why does it have *both* DeleteText and DeleteRange?
	// reason: avoid erorr handling every time getting end pos.
	// Theoretically, below also works:
	//
	//   EditDelete struct {
	//     DeleteText string
	//     StartPos   Position
	//   }
	//
	// However, it requires calculation of end position on the fly,
	// which can result in error. (at least from type signature perspective)
	// So, better to store the entire Range instead.
}

func (e EditInsert) Apply(filename string) error {
	return InsertInFile(filename, e.Position, e.NewText)
}

func (e EditDelete) Apply(filename string) error {
	return DeleteInFile(filename, e.DeleteRange)
}

func (e EditInsert) Split(strategy SplitStrategy) ([]Edit, error) {
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

func (e EditDelete) Split(strategy SplitStrategy) ([]Edit, error) {
	switch strategy {
	case SplitByLine:
		return splitDeleteByLine(e)
	case SplitByWord:
		return splitDeleteByWord(e)
	case SplitByChar:
		return splitDeleteByChar(e)
	default:
		return nil, nil
	}
}
