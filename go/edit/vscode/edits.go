package vscode

type Edit interface {
	IsEdit() bool
}

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


