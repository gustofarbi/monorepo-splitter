package action

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"os"
	"os/exec"
	"splitter/pkg"
	"splitter/version"
	"splitter/version/semver"
)

type Validate struct{}

func (v Validate) Act(collection *pkg.PackageCollection) error {
	if err := validateRootPackage(collection); err != nil {
		return fmt.Errorf("invalid root package: %s", err)
	}
	if collection.Conf.PackageAuth == nil {
		if auth, err := collection.Conf.PackageAuthFunc(); err != nil {
			return fmt.Errorf("cannot authenticate: %s", err)
		} else {
			collection.Conf.PackageAuth = auth
		}
	}
	for _, singlePackage := range collection.Packages {
		if err := validateSinglePackage(
			singlePackage.Repo,
			singlePackage.RemoteName,
			singlePackage.RemoteUrl,
			collection.Conf.VersionValue,
			collection.Conf.PackageAuth,
		); err != nil {
			return fmt.Errorf("invalid package %s: %s", singlePackage.RemoteName, err)
		}
	}

	return nil
}

func (v Validate) Description() string {
	return "validate configuration"
}

func validateRootPackage(collection *pkg.PackageCollection) error {
	if err := os.Chdir(collection.RootPackage.Path); err != nil {
		return fmt.Errorf("cannot change dir: %s", err)
	}
	cmd := exec.Command("git", "status", "--porcelain")
	buf := bytes.Buffer{}
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cannot get repo status: %s", err)
	}
	if buf.String() != "" {
		return errors.New("root repo contains unstaged changes")
	}

	if collection.Conf.Root.Branch != "" {
		cmd = exec.Command("git", "checkout", collection.Conf.Root.Branch)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("cannot checkout root branch: %s: %s", collection.Conf.Root.Branch, err)
		}
	}
	cmd = exec.Command("git", "pull")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cannot pull root branch: %s", err)
	}

	return nil
}

func validateSinglePackage(
	repo *git.Repository,
	remote, url string,
	newVersion version.Version,
	auth transport.AuthMethod,
) error {
	_, err := repo.CreateRemote(&config.RemoteConfig{
		Name: remote,
		URLs: []string{url},
	})
	if err != nil && err != git.ErrRemoteExists {
		return fmt.Errorf("error creating remote %s %s: %s", remote, url, err)
	}

	err = repo.Fetch(&git.FetchOptions{RemoteName: remote, Auth: auth})
	if err != nil && err != transport.ErrEmptyRemoteRepository {
		return fmt.Errorf("error fetching remote %s: %s", remote, err)
	}
	iter, err := repo.Tags()
	if err != nil {
		return fmt.Errorf("error fetching tags from %s: %s", remote, err)
	}
	tags := semver.NewSemverCollection()
	if err = iter.ForEach(func(ref *plumbing.Reference) error {
		if tag, err := semver.FromTag(ref.Name().String()); err == nil {
			tags.Add(tag)
		}
		return nil
	}); err != nil {
		return err
	}
	highest := tags.GetHighest()
	if highest.IsGreater(newVersion) {
		return fmt.Errorf("version from config %s is higher than an existing tag %s", newVersion, highest)
	}

	return nil
}

func (v Validate) String() string {
	return "validate"
}
