package action

import (
	"fmt"
	"os"
	"os/exec"
	"splitter/pkg"
)

type CommitChanges struct{}

func (t CommitChanges) Act(collection *pkg.PackageCollection) error {
	if err := os.Chdir(collection.RootPackage.Path); err != nil {
		return fmt.Errorf("cannot change directory to: %s: %+v", collection.RootPackage.Path, err)
	}

	// stage all changes
	cmd := exec.Command("git", "add", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error adding changes to git: %+v", err)
	}

	// commit changes
	cmd = exec.Command(
		"git",
		"commit",
		"-m",
		fmt.Sprintf("'prepare release %s'", collection.Conf.VersionValue.GitTag()),
		"--no-verify",
	)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error commiting changes: %+v", err)
	}

	return nil
}

func (t CommitChanges) Description() string {
	return "add to git and commit changes"
}

func (t CommitChanges) String() string {
	return "commit-changes"
}
