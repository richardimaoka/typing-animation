package gitpkg

import (
	"fmt"
	"testing"
)

func TestCommitsForFileInternal(t *testing.T) {
	orgname := "spf13"
	reponame := "cobra"
	filePath := "command.go"

	repo, err := OpenOrClone(orgname, reponame)
	if err != nil {
		t.Fatal(err)
	}

	commits, err := commitsForFileInternal(repo, filePath)
	if err != nil {
		t.Fatal(err)
	}

	for i, c := range commits {
		hash := string([]rune(c.Hash.String())[:7])
		parent, _ := c.Parent(0)
		parentHash := string([]rune(parent.Hash.String())[:7])

		fmt.Printf("%3d %s %s %s %s\n", i, hash, parentHash, c.Author.When, string([]rune(c.Message)[:15]))
	}
}
