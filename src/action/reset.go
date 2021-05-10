package action

import (
	"io/ioutil"
	"os"
	"os/exec"
	"splitter/pkg"
)

type Reset struct{}

func (r Reset) Act(collection *pkg.PackageCollection) {
	err := os.Chdir(collection.Conf.Root.Path)
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("git", "reset", "--hard", "HEAD^")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		b, _ := ioutil.ReadAll(os.Stderr)
		panic(string(b))
	}
}

func (r Reset) Description() string {
	return "reset changes to original state"
}

func (r Reset) String() string {
	return "reset"
}
