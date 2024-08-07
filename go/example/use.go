package example

import (
	"fmt"
	"os"
	"time"

	"github.com/richardimaoka/typing-animation/go/edit/vscode"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	text1 = `func newContext(ctx context.Context, resource, name string, gvk schema.GroupVersionKind) context.Context {
	oldInfo, found := genericapirequest.RequestInfoFrom(ctx)
	if !found {
		return ctx
	}
	newInfo := genericapirequest.RequestInfo{
		IsResourceRequest: true,
		Verb:              "get",
		Namespace:         oldInfo.Namespace,
		Resource:          resource,
		Name:              name,
		Parts:             []string{resource, name},`
	text2 = `func newContext(ctx context.Context, resource, name, namespace string, gvk schema.GroupVersionKind) context.Context {
	newInfo := genericapirequest.RequestInfo{
		IsResourceRequest: true,
		Verb:              "get",
		Namespace:         namespace,
		Resource:          resource,
		Name:              name,
		Parts:             []string{resource, name},`

//	text1 = `I am the very model of a modern Major-General,
//
// I've information vegetable, animal, and mineral,
// I know the kings of England, and I quote the fights historical,
// From Marathon to Waterloo, in order categorical.`
//
//	text2 = `I am the very model of a cartoon individual,
//
// My animation's comical, animal, and whimsical,
// I'm quite adept at funny gags, comedic theory I have read,
// From wicked puns and stupid jokes to anvils that drop on your head.`
)

func buildStack(diffs []diffmatchpatch.Diff) *vscode.EditStack {
	stack := vscode.NewEditStack()

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

func Experiment() {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(text1, text2, true)
	fmt.Printf("%+v\n", diffs)

	buildStack(diffs)
}

func ExperimentFiles() {
	before, err := os.ReadFile("./example/before.txt")
	if err != nil {
		panic(err)
	}

	after, err := os.ReadFile("./example/after.txt")
	if err != nil {
		panic(err)
	}

	resultFile := "./example/result.txt"
	if err = os.WriteFile(resultFile, before, 0666); err != nil {
		panic(err)
	}
	defer func() {
		if err = os.WriteFile(resultFile, before, 0666); err != nil {
			panic(err)
		}
	}()

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(before), string(after), false)

	for i, d := range diffs {
		fmt.Printf("[%2d]:type = %s, text = '%s'\n", i, d.Type.String(), d.Text)
	}

	stack := buildStack(diffs)
	edits, err := stack.CalcEdits()
	if err != nil {
		panic(err)
	}

	for _, e := range edits {
		e.ApplyToFile(resultFile)
		time.Sleep(300 * time.Millisecond)
	}
}
