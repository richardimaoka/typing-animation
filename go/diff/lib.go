package diff

import (
	"fmt"

	"github.com/richardimaoka/typing-animation/go/edit/monaco"
	"github.com/richardimaoka/typing-animation/go/edit/vscode"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func createStack(before, after string) *vscode.EditStack {
	stack := vscode.NewEditStack()
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(before, after, true)

	for _, d := range diffs {
		switch d.Type {
		case diffmatchpatch.DiffDelete:
			stack.AppendDelete(d.Text)
		case diffmatchpatch.DiffEqual:
			stack.AppendEqual(d.Text)
		case diffmatchpatch.DiffInsert:
			stack.AppendInsert(d.Text)
		default:
			panic(fmt.Sprintf("unexpected diff match patch type = %d (%s)", d.Type, d.Type.String()))
		}
	}

	return stack
}

func CalcEdits(before, after string) ([]vscode.Edit, error) {
	stack := createStack(before, after)

	edits, err := stack.CalcEdits()
	if err != nil {
		return nil, fmt.Errorf("diff.CalcEdits failed, %s", err)
	}

	return edits, nil
}

func CalcMonacoEdits(before, after string) ([]monaco.SingleEditOperation, error) {
	stack := createStack(before, after)

	edits, err := stack.CalcMonacoEdits()
	if err != nil {
		return nil, fmt.Errorf("diff.CalcMonacoEdits failed, %s", err)
	}

	return edits, nil
}
