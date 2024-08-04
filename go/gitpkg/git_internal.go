package gitpkg

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// private top-level variable, to manipulate in testing within the same package
var githubDir string = "/tmp/github.com"

func localRepoPath(orgname, reponame string) string {
	return fmt.Sprintf("%s/%s/%s", githubDir, orgname, reponame)
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

func commitsForFileInternal(repo *git.Repository, filepath string) ([]*object.Commit, error) {
	// wt, err := repo.Worktree()
	// if err != nil {
	// 	return nil, err
	// }

	// err = wt.Checkout(&git.CheckoutOptions{
	// 	Branch: plumbing.ReferenceName(branchName),
	// })
	// if err != nil {
	// 	return nil, err
	// }

	headRef, err := repo.Head()
	if err != nil {
		return nil, err
	}

	commit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		return nil, err
	}

	var commits []*object.Commit
	for {
		parent, err := commit.Parent(0)
		if err == object.ErrParentNotFound {
			break
		} else if err != nil {
			return nil, err
		}

		patch, err := parent.Patch(commit)
		if err != nil {
			return nil, err
		}

		filePatches := patch.FilePatches()
		for _, fp := range filePatches {
			from, to := fp.Files()
			if from == nil {
				// If the patch creates a new file, "from" will be nil.
				if to.Path() == filepath {
					commits = append(commits, commit) //added by this commit
				}
			} else if to == nil {
				// If the patch deletes a file, "to" will be nil.
				if from.Path() == filepath {
					commits = append(commits, commit) //deleted by this commit
				}
			} else {
				if from.Path() == filepath {
					commits = append(commits, commit) //updated by this commit
				}
			}
		}

		// for the next loop iteration
		commit = parent
	}

	return commits, err
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
