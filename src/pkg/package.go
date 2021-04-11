package pkg

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"path/filepath"
	"splitter/composer"
	"splitter/conf"
)

type Package struct {
	Composer   *composer.Composer
	Repo       *git.Repository
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

	return &Package{
		Composer:   rootComposer,
		Repo:       rootRepo,
		Path:       conf.Root.Path,
		RemoteName: conf.Root.Remote,
	}, nil
}

func loadPackage(pkg *conf.Pkg, cnf *conf.Config) (*Package, error) {
	repo, err := git.Init(memory.NewStorage(), nil)
	if err != nil {
		return nil, err
	}
	localPath := filepath.Join(cnf.Root.Path, cnf.Packages.Prefix, pkg.Path)
	packageComposer, err := composer.LoadComposer(localPath)
	if err != nil {
		return nil, err
	}
	return &Package{
		Composer:   packageComposer,
		Repo:       repo,
		Path:       filepath.Join(cnf.Packages.Prefix, pkg.Path),
		RemoteName: pkg.Remote,
		RemoteUrl:  pkg.Url,
	}, nil
}
