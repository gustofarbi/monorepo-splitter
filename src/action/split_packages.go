package action

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/gustofarbi/lite/splitter"
	"log"
	"os"
	"os/exec"
	"splitter/utils/pkg"
)

type SplitPackages struct{}

func (s SplitPackages) Act(collection *pkg.PackageCollection) {
	rootLocalPath := collection.RootPackage.Path
	for _, singlePkg := range collection.Packages {
		_, err := collection.RootPackage.Repo.CreateRemote(&config.RemoteConfig{
			Name: singlePkg.RemoteName,
			URLs: []string{singlePkg.RemoteUrl},
		})
		if err != nil && err != git.ErrRemoteExists {
			log.Fatalf("cannot create remote %s: %s", singlePkg.RemoteName, err)
		}
		err = os.Chdir(rootLocalPath)
		if err != nil {
			log.Fatalf("cannot change directory: %s", err)
		}
		prefix := &splitter.Prefix{
			From: singlePkg.Path,
		}
		splitConfig := &splitter.Config{
			Prefixes:   []*splitter.Prefix{prefix},
			Origin:     "HEAD",
			Path:       ".",
			GitVersion: "latest",
		}

		result := &splitter.Result{}
		err = splitter.Split(splitConfig, result)
		if err != nil {
			fmt.Println("cannot split: " + err.Error())
		}
		// needs to be done via cmdline because of this https://github.com/go-git/go-git/issues/105
		err = os.Chdir(singlePkg.Path)
		if err != nil {
			log.Fatalf("cannot change directory: %s", err)
		}

		repo := singlePkg.Repo
		cmd := exec.Command(
			"git",
			"push",
			singlePkg.RemoteName,
			fmt.Sprintf("%s:refs/heads/%s", result.Head().String(), collection.Conf.Packages.Branch),
			"-f",
		)
		println(cmd.String())
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			b, _ := cmd.Output()
			log.Fatalf("pushing splitted changes to remote failed on %s: %s", singlePkg.RemoteName, string(b))
		}
		err = repo.Fetch(&git.FetchOptions{
			RemoteName: singlePkg.RemoteName,
		})
		t, err := repo.Object(plumbing.AnyObject, plumbing.NewHash(result.Head().String()))
		if err != nil {
			panic(err)
		}
		_, err = repo.CreateTag(collection.Conf.Semver.GitTag(), t.ID(), &git.CreateTagOptions{
			Message: "prepare release",
		})
		if err != nil {
			log.Fatalf("cannot create tag: %s", err)
		}
		po := &git.PushOptions{
			RemoteName: singlePkg.RemoteName,
			RefSpecs:   []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")},
			Auth:       collection.Auth,
		}
		err = repo.Push(po)
		if err != nil {
			panic(err)
		}
	}
}

func (s SplitPackages) Description() string {
	return "split packages into separate repositories and push the changes to their remotes"
}
