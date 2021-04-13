package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os/user"
	"path/filepath"
)

func main() {
	currentUser, err := user.Current()
	check(err)
	println("open repo")
	repoPath := filepath.Join(currentUser.HomeDir, "workspace/split-demo")
	repo, err := git.PlainOpen(repoPath)
	check(err)
	println("acquire worktree")
	worktree, err := repo.Worktree()
	check(err)
	println("stage changes")
	_, err = worktree.Add(".")
	check(err)
	println("commit changes")
	commit, err := worktree.Commit("add stuff", &git.CommitOptions{})
	check(err)
	println("commit objects")
	_, err = repo.CommitObject(commit)
	check(err)
	auth := &http.BasicAuth{
		Username: "gustofarbi",
		Password: "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",
	}
	println("push changes")
	err = repo.Push(&git.PushOptions{Auth: auth})
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
