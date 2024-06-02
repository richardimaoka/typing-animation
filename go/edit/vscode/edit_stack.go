package vscode

type EditStack struct {
	diffs []Diff
	// currentPosition Position
}

func NewEditStack() *EditStack {
	return &EditStack{}
}

func (s *EditStack) AppendDelete(text string) {
	delete := Diff{Type: DiffDelete, Text: text}
	s.diffs = append(s.diffs, delete)
}

func (s *EditStack) AppendEqual(text string) {
	equal := Diff{Type: DiffEqual, Text: text}
	s.diffs = append(s.diffs, equal)
}

func (s *EditStack) AppendInsert(text string) {
	insert := Diff{Type: DiffInsert, Text: text}
	s.diffs = append(s.diffs, insert)
}

func (s *EditStack) CalcEdits() ([]Edit, error) {
	edits := []Edit{}

	pos := Position{0, 0}
	var edit Edit = nil
	var err error

	for _, diff := range s.diffs {
		edit, pos, err = diffToEdit(pos, diff)
		if err != nil {
			return nil, err
		}

		edits = append(edits, edit)
	}

	return edits, nil
}
