package composer

import (
	"os"
	"syscall"
	"testing"
)

func TestComposer_LoadComposer(t *testing.T) {
	c, err := LoadComposer("composer.json")
	if err != nil {
		t.Fatalf("loading failed: %s", err)
	}
	if len(c.Items.Require) != 5 {
		t.Fatal("wrong number of requires")
	}
	if len(c.Items.RequireDev) != 5 {
		t.Fatal("wrong number of require devs")
	}
	if len(c.Items.Replace) != 5 {
		t.Fatal("wrong number of replaces")
	}
	if c.Items.Name != "mp/myposter" {
		t.Fatal("wrong pkg name")
	}
}

func TestComposer_WriteToFile(t *testing.T) {
	defer syscall.Unlink("copy.json")
	c, _ := LoadComposer("composer.json")
	err := c.WriteToFile("copy.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat("copy.json"); err != nil {
		t.Fatal(err)
	}
}
