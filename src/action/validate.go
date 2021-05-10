package action

import (
	"bytes"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"os"
	"os/exec"
	"splitter/pkg"
	"splitter/version"
)

type Validate struct{}

func (v Validate) Act(collection *pkg.PackageCollection) {
	validateRootPackage(collection)
	if collection.Conf.PackageAuth == nil {
		collection.Conf.PackageAuth = collection.Conf.PackageAuthFunc()
	}
	for _, singlePackage := range collection.Packages {
		validateSinglePackage(
			singlePackage.Repo,
			singlePackage.RemoteName,
			singlePackage.RemoteUrl,
			collection.Conf.Semver,
			collection.Conf.PackageAuth,
		)
	}
}

func (v Validate) Description() string {
	return "validate configuration"
}

func validateRootPackage(collection *pkg.PackageCollection) {
	if err := os.Chdir(collection.RootPackage.Path); err != nil {
		panic(fmt.Sprintf("cannot change dir: %+v", err))
	}
	cmd := exec.Command("git", "status", "--porcelain")
	buf := bytes.Buffer{}
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(fmt.Sprintf("cannot get repo status: %+v", err))
	}
	if buf.String() != "" {
		panic("root repo contains unstaged changes")
	}
	cmd = exec.Command("git", "checkout", collection.Conf.Root.Branch)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	cmd = exec.Command("git", "pull")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func validateSinglePackage(
	repo *git.Repository,
	remote, url string,
	newVersion version.Semver,
	auth transport.AuthMethod) {
	_, err := repo.CreateRemote(&config.RemoteConfig{
		Name: remote,
		URLs: []string{url},
	})
	if err != nil && err != git.ErrRemoteExists {
		panic(fmt.Sprintf("error creating remote %s %s: %+v", remote, url, err))
	}

	err = repo.Fetch(&git.FetchOptions{RemoteName: remote, Auth: auth})
	if err != nil && err != transport.ErrEmptyRemoteRepository {
		panic(fmt.Sprintf("error fetching remote %s: %+v", remote, err))
	}
	iter, err := repo.Tags()
	if err != nil {
		panic(fmt.Sprintf("error fetching tags from %s: %+v", remote, err))
	}
	tags := version.NewSemverCollection()
	err = iter.ForEach(func(ref *plumbing.Reference) error {
		tag := version.FromTag(ref.Name().String())
		tags.Add(tag)
		return nil
	})
	highest := tags.GetHighest()
	if highest.IsGreater(newVersion) {
		panic(fmt.Sprintf("version from config %s is higher than an existing tag %s", newVersion, highest))
	}
}

func (v Validate) String() string {
	return "validate"
}
