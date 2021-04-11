package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"splitter/utils/version"
	"strings"
)

const (
	configName = "splitter"
)

var extensions = [2]string{"yaml", "yml"}

type Config struct {
	Root     `yaml:"root"`
	Packages `yaml:"packages"`
	Version  string `yaml:"version"`
	version.Semver
}

type Root struct {
	Branch string `yaml:"branch"`
	Path   string `yaml:"path"`
	Remote string `yaml:"remote"`
}

type Packages struct {
	Prefix string `yaml:"prefix"`
	Items  []*Pkg `yaml:"items"`
	Branch string `yaml:"branch"`
}

type Pkg struct {
	Remote string `yaml:"remote"`
	Url    string `yaml:"url"`
	Path   string `yaml:"path"`
}

func LoadConfig(names ...string) (*Config, error) {
	if len(names) > 0 && names[0] != "" {
		return loadConfig(names[0])
	}
	for _, ext := range extensions {
		filename := fmt.Sprintf("%s.%s", configName, ext)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			continue
		}
		return loadConfig(filename)
	}
	return nil, fmt.Errorf("no suitable conf file found")
}

func loadConfig(filename string) (*Config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var c Config
	err = yaml.Unmarshal(b, &c)
	c.Semver = version.FromString(c.Version)
	if strings.HasPrefix(c.Root.Path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		c.Root.Path = filepath.Join(homeDir, c.Root.Path[1:])
	}
	if err != nil {
		return nil, err
	}
	for _, item := range c.Items {
		if item.Path == "" {
			item.Path = item.Remote
		}
	}

	return &c, nil
}
