package conf

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	c, err := LoadConfig("conf.yaml")
	if err != nil {
		t.Fatalf("loading failed: %s", err)
	}
	if c.Root.Branch != "master" {
		t.Fatal("wrong root branch")
	}
	if c.Root.Path != "~/workspace/split-demo" {
		t.Fatal("wrong root local path")
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
	if c.Version != "1.5.3" {
		t.Fatal("wrong version")
	}
}
