package composer

import (
	"fmt"
	json "github.com/json-iterator/go"
	"os"
	"path/filepath"
	"strings"
)

const composerFilename = "composer.json"

const (
	require    = "require"
	requireDev = "require-dev"
	replace    = "replace"
	name       = "name"
	config     = "config"

	VendorDir = "vendor-dir"
)

//todo do this via json

type Composer struct {
	Items struct {
		Require    map[string]string      `json:"require,omitempty"`
		RequireDev map[string]string      `json:"require-dev,omitempty"`
		Replace    map[string]string      `json:"replace,omitempty"`
		Config     map[string]interface{} `json:"config"`
		Name       string                 `json:"name"`
	}
	Rest map[string]interface{}
}

func LoadComposer(path string) (*Composer, error) {
	var composerPath string
	if strings.HasSuffix(path, composerFilename) {
		composerPath = path
	} else {
		composerPath = filepath.Join(path, composerFilename)
	}
	if _, err := os.Stat(composerPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("composer.json file not found at %s", composerPath)
	}
	b, err := os.ReadFile(composerPath)
	if err != nil {
		return nil, err
	}

	var c Composer
	err = json.Unmarshal(b, &c.Items)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &c.Rest)
	if err != nil {
		return nil, err
	}

	delete(c.Rest, require)
	delete(c.Rest, requireDev)
	delete(c.Rest, replace)
	delete(c.Rest, name)
	delete(c.Rest, config)

	return &c, nil
}

func (c *Composer) WriteToFile(path string) error {
	if c.Items.Require != nil {
		c.Rest[require] = c.Items.Require
	}
	if c.Items.RequireDev != nil {
		c.Rest[requireDev] = c.Items.RequireDev
	}
	if c.Items.Replace != nil {
		c.Rest[replace] = c.Items.Replace
	}
	if c.Items.Config != nil {
		c.Rest[config] = c.Items.Config
	}
	c.Rest[name] = c.Items.Name

	b, err := json.MarshalIndent(c.Rest, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, os.FileMode(0644))

}
