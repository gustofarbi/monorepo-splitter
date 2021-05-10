package action

import (
	"fmt"
	"os"
	"os/exec"
	"splitter/pkg"
)

type CommitChanges struct{}

func (t CommitChanges) Act(collection *pkg.PackageCollection) {
	if err := os.Chdir(collection.RootPackage.Path); err != nil {
		panic(err)
	}

	// stage all changes
	cmd := exec.Command("git", "add", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	// commit changes
	cmd = exec.Command(
		"git",
		"commit",
		"-m",
		fmt.Sprintf("'prepare release %s'", collection.Conf.Semver.GitTag()),
	)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func (t CommitChanges) Description() string {
	return "add to git, commit changes and tag the new release"
}

func (t CommitChanges) String() string {
	return "push-changes"
}
