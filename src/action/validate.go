package action

import (
	"bytes"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"log"
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
	err := os.Chdir(collection.RootPackage.Path)
	if err != nil {
		log.Fatalf("cannot change dir: %+v", err)
	}
	cmd := exec.Command("git", "status", "--porcelain")
	buf := bytes.Buffer{}
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("cannot get repo status: %+v", err)
	}
	if buf.String() != "" {
		panic("root repo contains unstaged changes")
	}
	cmd = exec.Command("git", "checkout", collection.Conf.Root.Branch)
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	cmd = exec.Command("git", "pull")
	err = cmd.Run()
	if err != nil {
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
		log.Fatalf("error creating remote %s %s: %+v", remote, url, err)
	}

	err = repo.Fetch(&git.FetchOptions{RemoteName: remote, Auth: auth})
	if err != nil && err != transport.ErrEmptyRemoteRepository {
		log.Fatalf("error fetching remote %s: %+v", remote, err)
	}
	iter, err := repo.Tags()
	if err != nil {
		log.Fatalf("error fetching tags from %s: %+v", remote, err)
	}
	tags := version.NewSemverCollection()
	err = iter.ForEach(func(ref *plumbing.Reference) error {
		tag := version.FromTag(ref.Name().String())
		tags.Add(tag)
		return nil
	})
	highest := tags.GetHighest()
	if highest.IsGreater(newVersion) {
		log.Fatalf("version from config %s is higher than an existing tag %s", newVersion, highest)
	}
}

func (v Validate) String() string {
	return "validate"
}
