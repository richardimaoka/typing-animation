package gitpkg

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func localRepoPath(orgname, reponame string) string {
	return fmt.Sprintf("/tmp/github.com/%s/%s", orgname, reponame)
}

func toHash(hashString string) (plumbing.Hash, error) {
	if !plumbing.IsHash(hashString) {
		return plumbing.ZeroHash, fmt.Errorf("'%s' is an invalid git hash", hashString)
	}

	hash := plumbing.NewHash(hashString)
	return hash, nil
}

func commitObjectInternal(repo *git.Repository, hashString string) (*object.Commit, error) {
	hash, err := toHash(hashString)
	if err != nil {
		return nil, err
	}

	commit, err := repo.CommitObject(hash)
	if err != nil {
		return nil, err
	}

	return commit, err
}

func fileInCommitInternal(repo *git.Repository, hashString, filePath string) (*object.File, error) {
	commit, err := commitObjectInternal(repo, hashString)
	if err != nil {
		return nil, err
	}

	file, err := commit.File(filePath)
	if err != nil {
		return nil, fmt.Errorf("error in file = '%s', %s", filePath, err)
	}

	return file, err
}

func fileContentsInCommitInternal(repo *git.Repository, hashString, filePath string) (string, error) {
	file, err := fileInCommitInternal(repo, hashString, filePath)
	if err != nil {
		return "", err
	}

	contents, err := file.Contents()
	if err != nil {
		return "", err
	}

	return contents, err
}
