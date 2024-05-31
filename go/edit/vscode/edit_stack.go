package vscode

type EditStack struct {
	diffs []Diff
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

func (s *EditStack) CalcEdits() []Edit {
	return nil
}
