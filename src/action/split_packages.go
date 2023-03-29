package action

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/gustofarbi/lite/splitter"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"splitter/pkg"
	"splitter/version"
)

type SplitPackages struct {
	dryRun bool
}

func (s SplitPackages) Act(collection *pkg.PackageCollection) error {
	rootPath := collection.RootPackage.Path

	// fetch git credentials
	if collection.Conf.PackageAuth == nil {
		if auth, err := collection.Conf.PackageAuthFunc(); err != nil {
			return fmt.Errorf("cannot authenticate: %s", err)
		} else {
			collection.Conf.PackageAuth = auth
		}
	}

	for _, singlePkg := range collection.Packages {
		if err := createRemote(collection, singlePkg, rootPath); err != nil {
			return fmt.Errorf("cannot create remote %s: %s", singlePkg.RemoteName, err)
		}
		result, err := getSplitResult(singlePkg)
		if err != nil {
			return err
		}

		// needs to be done via cmdline because of this https://github.com/go-git/go-git/issues/105
		if err = os.Chdir(singlePkg.Path); err != nil {
			return fmt.Errorf("cannot change directory: %s", err)
		}

		if s.dryRun {
			fmt.Printf(
				"%s %s %s %s %s\n",
				"git",
				"push",
				singlePkg.RemoteName,
				fmt.Sprintf("%s:refs/heads/%s", result.Head().String(), collection.Conf.Packages.Branch),
				"-f",
			)
			if err = reset(); err != nil {
				return err
			}
		} else {
			if err = pushAndTag(
				singlePkg,
				result,
				collection.Conf.Packages.Branch,
				collection.Conf.PackageAuth,
				collection.Conf.VersionValue,
			); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s SplitPackages) Description() string {
	if s.dryRun {
		return "reset packages to previous state"
	} else {
		return "split packages into separate repositories and push the changes to their remotes"
	}
}

func (s SplitPackages) String() string {
	return "split-packages"
}

func createRemote(collection *pkg.PackageCollection, singlePkg *pkg.Package, rootPath string) error {
	_, err := collection.RootPackage.Repo.CreateRemote(&config.RemoteConfig{
		Name: singlePkg.RemoteName,
		URLs: []string{singlePkg.RemoteUrl},
	})
	if err != nil && err != git.ErrRemoteExists {
		return fmt.Errorf("cannot create remote %s: %s", singlePkg.RemoteName, err)
	}
	if err = os.Chdir(rootPath); err != nil {
		return fmt.Errorf("cannot change directory: %+v", err)
	}

	return nil
}

func getSplitResult(singlePkg *pkg.Package) (*splitter.Result, error) {
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
		return nil, fmt.Errorf("cannot split package %s: %s", singlePkg.RemoteName, err)
	}
	return result, nil
}

func pushAndTag(
	singlePkg *pkg.Package,
	result *splitter.Result,
	targetBranch string,
	packageAuth http.AuthMethod,
	version version.Version,
) error {
	cmd := exec.Command(
		"git",
		"push",
		singlePkg.RemoteName,
		fmt.Sprintf("%s:refs/heads/%s", result.Head().String(), targetBranch),
		"-f",
	)
	fmt.Println(cmd.String())
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error pushing to remote %s: %s", singlePkg.RemoteName, err)
	}
	repo := singlePkg.Repo
	if err := repo.Fetch(&git.FetchOptions{
		RemoteName: singlePkg.RemoteName,
		Force:      true,
		Progress:   os.Stdout,
		Auth:       packageAuth,
	}); err != nil {
		return fmt.Errorf("error fetching from remote %s: %s", singlePkg.RemoteName, err)
	}
	t, err := repo.Object(plumbing.AnyObject, plumbing.NewHash(result.Head().String()))
	if err != nil {
		return err
	}
	if _, err = repo.CreateTag(version.String(), t.ID(), &git.CreateTagOptions{
		Message: "prepare release",
	}); err != nil {
		if err.Error() == "tag already exists" {
			return nil
		}
		return fmt.Errorf("cannot create tag: %s", err)
	}
	po := &git.PushOptions{
		RemoteName: singlePkg.RemoteName,
		RefSpecs:   []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")},
		Auth:       packageAuth,
	}
	if err = repo.Push(po); err != nil {
		log.Printf("cannot push tag to remote repository: %s", err)
	}

	return nil
}

func reset() error {
	cmd := exec.Command("git", "reset", "--hard", "HEAD^")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		b, _ := ioutil.ReadAll(os.Stderr)
		return fmt.Errorf("cannot reset root package: %s", string(b))
	}

	return nil
}
