package action

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"log"
	"os"
	"splitter/pkg"
)

type TagRelease struct{}

func (t TagRelease) Act(collection *pkg.PackageCollection) {
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
	// needs to be done via cmdline because of this https://github.com/go-git/go-git/issues/105
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
		Auth:       collection.Conf.AuthMethod,
	}
	err = collection.RootPackage.Repo.Push(po)
	if err != nil {
		log.Fatalln("cannot push tags into root repo", err)
	}
}

func (t TagRelease) Description() string {
	return "add to git, commit changes and tag the new release"
}

func (t TagRelease) String() string {
	return "tag-release"
}
