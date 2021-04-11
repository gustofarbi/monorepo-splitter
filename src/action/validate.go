package action

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"log"
	"splitter/pkg"
	"splitter/version"
)

type Validate struct{}

func (v Validate) Act(collection *pkg.PackageCollection) {
	validateRootPackage(collection.RootPackage.Repo, collection.Conf.Root.Branch)
	for _, singlePackage := range collection.Packages {
		validateSinglePackage(
			singlePackage.Repo,
			singlePackage.RemoteName,
			singlePackage.RemoteUrl,
			collection.Conf.Semver,
		)
	}
}

func (v Validate) Description() string {
	return "validate configuration"
}

func validateRootPackage(repo *git.Repository, branch string) {
	workTree, err := repo.Worktree()
	if err != nil {
		log.Fatalf("error obtaining worktree: %+v", err)
	}
	err = workTree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/heads/" + branch),
	})
	if err != nil {
		log.Fatalf("error checking out root repo: %+v", err)
	}
	err = workTree.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		log.Fatalf("error pulling: %+v", err)
	}
}

func validateSinglePackage(repo *git.Repository, remote, url string, newVersion version.Semver) {
	_, err := repo.CreateRemote(&config.RemoteConfig{
		Name: remote,
		URLs: []string{url},
	})
	if err != nil && err != git.ErrRemoteExists {
		log.Fatalf("error creating remote %s %s: %+v", remote, url, err)
	}

	err = repo.Fetch(&git.FetchOptions{RemoteName: remote})
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
