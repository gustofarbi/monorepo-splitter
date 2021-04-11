package pkg

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"
	"path/filepath"
	"splitter/utils/composer"
	"splitter/utils/conf"
	"splitter/utils/version"
)

type Package struct {
	Composer   *composer.Composer
	Repo       *git.Repository
	Tags       *version.SemverCollection
	Path       string
	RemoteName string
	RemoteUrl  string
}

func loadRootPackage(conf *conf.Config) (*Package, error) {
	rootComposer, err := composer.LoadComposer(conf.Root.Path)
	if err != nil {
		return nil, err
	}
	rootRepo, err := git.PlainOpen(conf.Root.Path)
	if err != nil {
		return nil, err
	}
	workTree, err := rootRepo.Worktree()
	if err != nil {
		return nil, err
	}
	err = workTree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/heads/" + conf.Root.Branch),
	})
	if err != nil {
		return nil, fmt.Errorf("error checking out root repo: %s", err)
	}
	err = workTree.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return nil, err
	}
	return &Package{
		Composer:   rootComposer,
		Repo:       rootRepo,
		Tags:       nil,
		Path:       conf.Root.Path,
		RemoteName: conf.Root.Remote,
	}, nil
}

func loadPackage(pkg *conf.Pkg, cnf *conf.Config) (*Package, error) {
	tags := version.NewSemverCollection()
	repo, err := git.Init(memory.NewStorage(), nil)
	if err != nil {
		return nil, err
	}
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: pkg.Remote,
		URLs: []string{pkg.Url},
	})
	if err != nil && err != git.ErrRemoteExists {
		return nil, err
	}

	err = repo.Fetch(&git.FetchOptions{RemoteName: pkg.Remote})
	if err != nil && err != transport.ErrEmptyRemoteRepository {
		return nil, err
	}
	iter, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	err = iter.ForEach(func(ref *plumbing.Reference) error {
		tag := version.FromTag(ref.Name().String())
		tags.Add(tag)
		return nil
	})
	localPath := filepath.Join(cnf.Root.Path, cnf.Packages.Prefix, pkg.Path)
	packageComposer, err := composer.LoadComposer(localPath)
	if err != nil {
		return nil, err
	}
	return &Package{
		Composer:   packageComposer,
		Repo:       repo,
		Tags:       tags,
		Path:       filepath.Join(cnf.Packages.Prefix, pkg.Path),
		RemoteName: pkg.Remote,
		RemoteUrl:  pkg.Url,
	}, nil
}
