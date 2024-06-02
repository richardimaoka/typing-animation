package vscode

// Interface representing atomic edit
type Edit interface {
	IsEdit() bool
	Apply(filename string) error
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
