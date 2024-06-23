package gitpkg

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

func localRepoPath(orgname, reponame string) string {
	return fmt.Sprintf("/tmp/github.com/%s/%s", orgname, reponame)
}

func OpenOrClone(orgname, reponame string) (*git.Repository, error) {
	localPath := localRepoPath(orgname, reponame)

	repo, err := git.PlainOpen(localPath)
	if err == git.ErrRepositoryNotExists {
		githubPath := fmt.Sprintf("https://github.com/%s/%s", orgname, reponame)

		repo, err = git.PlainClone(localPath, false, &git.CloneOptions{
			URL:      githubPath,
			Progress: os.Stdout,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to clone %s, %s", githubPath, err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to open %s/%s, %s", orgname, reponame, err)
	}

	return repo, nil
}
