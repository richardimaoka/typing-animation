package main

import (
	"fmt"

	"github.com/richardimaoka/typing-animation/go/diff"
	"github.com/richardimaoka/typing-animation/go/gitpkg"
)

func Experiment() {
	orgname := "spf13"
	reponame := "cobra"
	filePath := "command.go"

	repo, err := gitpkg.OpenOrClone(orgname, reponame)
	if err != nil {
		panic(err)
	}

	first := "9334a46bd6b3887f3561d705440038ec93b7f62e"  //Return an error in the case of unrunnable subcommand
	second := "51f06c7dd1e73470976107fc6931b21143b83676" //Correct all complaints from golint
	before, err := gitpkg.FileContentsInCommit(repo, first, filePath)
	if err != nil {
		panic(err)
	}

	after, err := gitpkg.FileContentsInCommit(repo, second, filePath)
	if err != nil {
		panic(err)
	}

	edits, err := diff.CalcEdits(before, after)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", edits)
}
