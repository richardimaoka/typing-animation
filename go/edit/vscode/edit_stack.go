package vscode

import (
	"fmt"

	"github.com/richardimaoka/typing-animation/go/edit/monaco"
)

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

func (s *EditStack) CalcEdits( /*TODO: splitStrategy: Strategy */ ) ([]Edit, error) {
	pos := Position{0, 0}
	edits := []Edit{}

	for _, diff := range s.diffs {
		edit, newPos, err := diffToEdit(pos, diff)
		if err != nil {
			return nil, err
		}

		if edit != nil {
			edits = append(edits, edit)
		}

		pos = newPos
	}

	return edits, nil
}

func (s *EditStack) CalcMonacoEdits() ([]monaco.SingleEditOperation, error) {
	currentPos := Position{0, 0}
	edits := []monaco.SingleEditOperation{}

	for _, diff := range s.diffs {
		rangeEndPos, err := editRangeEnd(currentPos, diff.Text)
		if err != nil {
			return nil, err
		}

		mRange := monaco.Range{
			StartColumn:     currentPos.Character,
			StartLineNumber: currentPos.Line,
			EndColumn:       rangeEndPos.Character,
			EndLineNumber:   rangeEndPos.Line,
		}

		switch diff.Type {
		case DiffInsert:
			edits = append(edits, monaco.SingleEditOperation{Text: diff.Text, Range: mRange, Operation: "Insert"})
			currentPos = rangeEndPos

		case DiffEqual:
			currentPos = rangeEndPos

		case DiffDelete:
			edits = append(edits, monaco.SingleEditOperation{Text: "" /*empty text for delete*/, Range: mRange, Operation: "Delete"})
			// currentPos doesn't move after delete

		default:
			return nil, fmt.Errorf("diff type = %d is invalid", diff.Type)
		}
	}

	return edits, nil
}
