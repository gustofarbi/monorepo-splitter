package conf

import (
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	c, err := LoadConfigWithVersion(
		"config.yaml",
		func() (http.AuthMethod, error) {
			return &http.BasicAuth{}, nil
		},
	)
	if err != nil {
		t.Fatalf("loading failed: %s", err)
	}
	if c.Root.Branch != "master" {
		t.Fatal("wrong root branch")
	}
	if c.Packages.Prefix != "packages" {
		t.Fatal("wrong packages prefix")
	}
	if len(c.Packages.Items) != 5 {
		t.Fatal("wrong number of packages")
	}
	if c.Packages.Branch != "master" {
		t.Fatal("wrong packages branch")
	}
	if c.Root.Remote != "origin" {
		t.Fatal("wrong root remote name")
	}
	if len(c.Actions) != 5 {
		t.Fatal("wrong number of actions")
	}
}
