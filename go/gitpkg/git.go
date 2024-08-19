package gitpkg

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func Open(orgname, reponame string) (*git.Repository, error) {
	errorPrefix := "gitpkg.Open failed"

	// TODO: specific error for org/repo non-existent even in GitHub

	localPath := localRepoPath(orgname, reponame)
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		return nil, fmt.Errorf("%s, %s", errorPrefix, err)
	}

	return repo, nil
}

func Clone(orgname, reponame string) (*git.Repository, error) {
	errorPrefix := "gitpkg.Clone failed"

	localPath := localRepoPath(orgname, reponame)
	repo, err := git.PlainClone(localPath, false, &git.CloneOptions{
		URL:      fmt.Sprintf("https://github.com/%s/%s", orgname, reponame),
		Progress: os.Stdout,
	})
	if err != nil {
		return nil, fmt.Errorf("%s, %s", errorPrefix, err)
	}

	return repo, nil
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

func CommitObject(repo *git.Repository, hashString string) (*object.Commit, error) {
	errorPrefix := "gitpkg.CommitObject failed"

	commit, err := commitObjectInternal(repo, hashString)
	if err != nil {
		return nil, fmt.Errorf("%s, %s", errorPrefix, err)
	}

	return commit, err
}

func CommitsForFile(orgname, reponame, filepath string) ([]*object.Commit, error) {
	errorPrefix := "gitpkg.CommitsForFile failed"

	repo, err := Open(orgname, reponame)
	if err != nil {
		return nil, err
	}

	commits, err := commitsForFileInternal(repo, filepath)
	if err != nil {
		return nil, fmt.Errorf("%s, %s", errorPrefix, err)
	}

	return commits, err
}

func FileInCommit(orgname, reponame, filePath, hashString string) (*object.File, error) {
	errorPrefix := "gitpkg.FileInCommit failed"

	repo, err := Open(orgname, reponame)
	if err != nil {
		return nil, err
	}

	file, err := fileInCommitInternal(repo, hashString, filePath)
	if err != nil {
		return nil, fmt.Errorf("%s, %s", errorPrefix, err)
	}

	return file, err
}

func FileContentsInCommit(repo *git.Repository, hashString, filePath string) (string, error) {
	errorPrefix := "gitpkg.FileContentsInCommit failed"

	contents, err := fileContentsInCommitInternal(repo, hashString, filePath)
	if err != nil {
		return "", fmt.Errorf("%s, %s", errorPrefix, err)
	}

	return contents, err
}

func RepoFiles(orgname, reponame string) ([]string, error) {
	repo, err := Open(orgname, reponame)
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

func RepoBranches(orgname, reponame string) ([]string, error) {
	repo, err := Open(orgname, reponame)
	if err != nil {
		return nil, err
	}

	iter, err := repo.Branches()
	if err != nil {
		return nil, err
	}

	var branches []string
	for branch, err := iter.Next(); err == nil; branch, err = iter.Next() {
		branches = append(branches, branch.Name().Short())
	}

	return branches, nil
}

func RepoFileContents(orgname, reponame, filePath, commitHashStr string) (string, error) {
	repo, err := OpenOrClone(orgname, reponame)
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
