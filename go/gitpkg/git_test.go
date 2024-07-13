package gitpkg

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	orgname := "spf13"
	reponame := "cobra"
	//  filePath :=

	repo, err := openOrClone(orgname, reponame)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		currentCommitHash string
		nextCommitHash    string
	}{
		{"b312f0a", "6d00909"},
	}

	for _, c := range cases {
		t.Run(c.currentCommitHash, func(t *testing.T) {
			commit, err := CommitObject(repo, c.currentCommitHash)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Print(commit)
			// currentCommit =
			// nextCommit =
		})
	}
}
