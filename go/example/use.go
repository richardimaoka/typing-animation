package example

import (
	"fmt"

	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	text1 = `I am the very model of a modern Major-General,
I've information vegetable, animal, and mineral,
I know the kings of England, and I quote the fights historical,
From Marathon to Waterloo, in order categorical.`
	text2 = `I am the very model of a cartoon individual,
My animation's comical, animal, and whimsical,
I'm quite adept at funny gags, comedic theory I have read,
From wicked puns and stupid jokes to anvils that drop on your head.`
)

func Experiment() {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(text1, text2, true)
	fmt.Printf("%+v\n", diffs)
}
