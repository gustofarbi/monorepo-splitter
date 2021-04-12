package action

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"log"
	"os"
	"os/exec"
	"splitter/pkg"
)

type CommitChanges struct{}

func (t CommitChanges) Act(collection *pkg.PackageCollection) {
	err := os.Chdir(collection.RootPackage.Path)
	if err != nil {
		panic(err)
	}

	// stage all changes
	cmd := exec.Command("git", "add", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	// commit changes
	cmd = exec.Command(
		"git",
		"commit",
		"-m",
		fmt.Sprintf("'prepare release %s'", collection.Conf.Semver.GitTag()),
	)
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	return
	// todo: for some reason go-git takes forever to stage changes
	rootRepo := collection.RootPackage.Repo
	workTree, err := rootRepo.Worktree()
	if err != nil {
		log.Fatalf("cannot acquire worktree: %s", err)
	}
	err = os.Chdir(collection.RootPackage.Path)
	if err != nil {
		log.Fatalf("cannot change directory: %s", err)
	}
	_, err = workTree.Add(".")
	if err != nil {
		log.Fatalf("cannot stage changes: %s", err)
	}
	commit, err := workTree.Commit("prepare release", &git.CommitOptions{})
	if err != nil {
		log.Fatalf("cannot commit changes: %s", err)
	}
	_, err = rootRepo.CommitObject(commit)
	if err != nil {
		log.Fatalf("cannot commit object: %s", err)
	}
	head, _ := rootRepo.Head()
	_, err = rootRepo.CreateTag(collection.Conf.Semver.String(), head.Hash(), &git.CreateTagOptions{
		Message: "prepare release",
	})
	if err != nil {
		log.Fatalf("cannot create tag: %s", err)
	}
	po := &git.PushOptions{
		RemoteName: collection.RootPackage.RemoteName,
		RefSpecs:   []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")},
		Auth:       collection.Conf.RootAuth,
	}
	err = collection.RootPackage.Repo.Push(po)
	if err != nil {
		log.Fatalln("cannot push tags into root repo", err)
	}
}

func (t CommitChanges) Description() string {
	return "add to git, commit changes and tag the new release"
}

func (t CommitChanges) String() string {
	return "push-changes"
}
