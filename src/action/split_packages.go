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
	"splitter/pkg"
)

type SplitPackages struct{}

func (s SplitPackages) Act(collection *pkg.PackageCollection) {
	rootPath := collection.RootPackage.Path

	// fetch git credentials
	if collection.Conf.PackageAuth == nil {
		collection.Conf.PackageAuth = collection.Conf.PackageAuthFunc()
	}

	for _, singlePkg := range collection.Packages {
		createRemote(collection, singlePkg, rootPath)
		result := getSplitResult(singlePkg)

		// needs to be done via cmdline because of this https://github.com/go-git/go-git/issues/105
		if err := os.Chdir(singlePkg.Path); err != nil {
			panic(fmt.Sprintf("cannot change directory: %+v", err))
		}

		repo := singlePkg.Repo
		cmd := exec.Command(
			"git",
			"push",
			singlePkg.RemoteName,
			fmt.Sprintf("%s:refs/heads/%s", result.Head().String(), collection.Conf.Packages.Branch),
			"-f",
		)
		log.Println(cmd.String())
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			panic(err)
		}
		if err := repo.Fetch(&git.FetchOptions{
			RemoteName: singlePkg.RemoteName,
			Force:      true,
			Progress:   os.Stdout,
			Auth:       collection.Conf.PackageAuth,
		}); err != nil {
			panic(err)
		}
		t, err := repo.Object(plumbing.AnyObject, plumbing.NewHash(result.Head().String()))
		if err != nil {
			panic(err)
		}
		if _, err = repo.CreateTag(collection.Conf.Semver.String(), t.ID(), &git.CreateTagOptions{
			Message: "prepare release",
		}); err != nil {
			panic(fmt.Sprintf("cannot create tag: %+v", err))
		}
		po := &git.PushOptions{
			RemoteName: singlePkg.RemoteName,
			RefSpecs:   []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")},
			Auth:       collection.Conf.PackageAuth,
		}
		if err = repo.Push(po); err != nil {
			panic(fmt.Sprintf("cannot push to remote repository: %+v", err))
		}
	}
}

func (s SplitPackages) Description() string {
	return "split packages into separate repositories and push the changes to their remotes"
}

func (s SplitPackages) String() string {
	return "split-packages"
}

func createRemote(collection *pkg.PackageCollection, singlePkg *pkg.Package, rootPath string) {
	_, err := collection.RootPackage.Repo.CreateRemote(&config.RemoteConfig{
		Name: singlePkg.RemoteName,
		URLs: []string{singlePkg.RemoteUrl},
	})
	if err != nil && err != git.ErrRemoteExists {
		panic(fmt.Sprintf("cannot create remote %s: %s", singlePkg.RemoteName, err))
	}
	if err = os.Chdir(rootPath); err != nil {
		panic(fmt.Sprintf("cannot change directory: %+v", err))
	}
}

func getSplitResult(singlePkg *pkg.Package) *splitter.Result {
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
	if err := splitter.Split(splitConfig, result); err != nil {
		panic("cannot split: " + err.Error())
	}
	return result
}