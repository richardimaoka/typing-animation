package gitpkg

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func localRepoPath(orgname, reponame string) string {
	return fmt.Sprintf("/tmp/github.com/%s/%s", orgname, reponame)
}

func openOrClone(orgname, reponame string) (*git.Repository, error) {
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

func ToHash(hashString string) (plumbing.Hash, error) {
	if !plumbing.IsHash(hashString) {
		return plumbing.ZeroHash, fmt.Errorf("'%s' is an invalid git hash", hashString)
	}

	hash := plumbing.NewHash(hashString)
	return hash, nil
}

func CommitObject(repo *git.Repository, hashString string) (*object.Commit, error) {
	errorPrefix := "gitpkg.CommitObject failed"

	hash, err := ToHash(hashString)
	if err != nil {
		return nil, fmt.Errorf("%s, %s", errorPrefix, err)
	}

	commit, err := repo.CommitObject(hash)
	if err != nil {
		return nil, fmt.Errorf("%s, %s", errorPrefix, err)
	}

	return commit, err
}

func RepoFiles(orgname, reponame string) ([]string, error) {
	repo, err := openOrClone(orgname, reponame)
	if err != nil {
		return nil, err
	}

	headRef, err := repo.Head()
	if err != nil {
		return nil, err
	}

	commit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		return nil, err
	}

	iter, err := commit.Files()
	if err != nil {
		return nil, err
	}

	var files []string
	for file, err := iter.Next(); err == nil; file, err = iter.Next() {
		files = append(files, file.Name)
	}

	return files, nil
}

func RepoFileContents(orgname, reponame, filePath, commitHashStr string) (string, error) {
	repo, err := openOrClone(orgname, reponame)
	if err != nil {
		return "", err
	}

	var commitHash plumbing.Hash
	if commitHashStr == "" {
		headRef, err := repo.Head()
		if err != nil {
			return "", err
		}
		commitHash = headRef.Hash()
	} else {
		commitHash = plumbing.NewHash(commitHashStr)
	}

	commit, err := repo.CommitObject(commitHash)
	if err != nil {
		return "", err
	}

	file, err := commit.File(filePath)
	if err == object.ErrFileNotFound {
		return "", fmt.Errorf("file = '%s' not found, %s", filePath, err)
	} else if err != nil {
		return "", err
	}

	contents, err := file.Contents()
	if err != nil {
		return "", err
	}

	return contents, nil
}
