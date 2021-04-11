package conf

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"splitter/version"
	"strings"
)

const (
	configName = "splitter"
)

var extensions = [2]string{"yaml", "yml"}

type AuthFunc func() http.AuthMethod

type Config struct {
	Root     `yaml:"root"`
	Packages `yaml:"packages"`
	Version  string   `yaml:"version"`
	Actions  []string `yaml:"actions"`
	version.Semver
	http.AuthMethod
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

func LoadConfig(name string, authFunc AuthFunc) (*Config, error) {
	if name != "" {
		return loadConfig(name, authFunc)
	}
	for _, ext := range extensions {
		filename := fmt.Sprintf("%s.%s", configName, ext)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			continue
		}
		return loadConfig(filename, authFunc)
	}
	return nil, fmt.Errorf("no suitable conf file found")
}

func loadConfig(filename string, authFunc AuthFunc ) (*Config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var c Config

	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}
	c.Semver = version.FromString(c.Version)
	c.AuthMethod = authFunc()

	if strings.HasPrefix(c.Root.Path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		c.Root.Path = filepath.Join(homeDir, c.Root.Path[1:])
	}
	for _, item := range c.Items {
		if item.Path == "" {
			item.Path = item.Remote
		}
	}

	return &c, nil
}
