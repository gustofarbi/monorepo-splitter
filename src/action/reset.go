package action

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"splitter/pkg"
)

type Reset struct {
	dryRun bool
}

func (r Reset) Act(collection *pkg.PackageCollection) error {
	err := os.Chdir(collection.Conf.Root.Path)
	if err != nil {
		return err
	}

	if r.dryRun {
		// nothing was written, so nothing to reset, just clean up
		cmd := exec.Command("git", "restore", ".")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			b, _ := ioutil.ReadAll(os.Stderr)
			return fmt.Errorf("cannot reset root package: %s", string(b))
		}
	} else {
		cmd := exec.Command("git", "reset", "--hard", "HEAD^")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			b, _ := ioutil.ReadAll(os.Stderr)
			return fmt.Errorf("cannot reset root package: %s", string(b))
		}
	}
	return nil
}

func (r Reset) Description() string {
	return "reset changes to original state"
}

func (r Reset) String() string {
	return "reset"
}
