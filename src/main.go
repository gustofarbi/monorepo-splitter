package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/utils/ioutil"
	"golang.org/x/term"
	"os"
	"path"
	"splitter/action"
	"splitter/conf"
	"splitter/pkg"
)

var (
	config    = flag.String("c", "", "configfile to use")
	overwrite = flag.Bool("o", false, "overwrite configs cache")

	dryRun = flag.Bool(
		"d",
		false,
		"instead of pushing changes, this option resets all repositories to previous state",
	)
)

type SplitterConfig struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func main() {
	flag.Parse()
	if *dryRun {
		fmt.Println("!!! running dry !!!")
	}
	collection, err := loadCollection()
	if err != nil {
		fmt.Printf("could not load collection: %s\n", err)
		return
	}
	if pipeline, err := action.NewPipeline(collection.Conf.Actions, *dryRun); err != nil {
		fmt.Printf("could not create a pipeline: %s\n", err)
		return
	} else {
		if err := pipeline.Run(collection); err != nil {
			panic(err)
		}
	}
}

func loadAuth() (http.AuthMethod, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot get home-dir: %s\n", err)
	}
	var sc SplitterConfig
	splitterCachePath := path.Join(homeDir, ".splitter", "config")
	if _, err := os.Stat(splitterCachePath); err == nil || *overwrite {
		var errCacheFile error
		bytes, errCacheFile := os.ReadFile(splitterCachePath)
		if errCacheFile != nil {
			return nil, errors.New("cannot read config cache file")
		} else if errCacheFile = json.Unmarshal(bytes, &sc); errCacheFile != nil {
			return nil, errors.New("cannot unmarshal config cache file")
		}
		if errCacheFile != nil {
			_ = os.Remove(splitterCachePath)
		} else {
			return &http.BasicAuth{
				Username: sc.Username,
				Password: sc.Token,
			}, nil
		}
	}
	var username string
	fmt.Printf("enter github credentials\nusername: ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		username = scanner.Text()
		break
	}
	fmt.Print("token: ")
	b, err := term.ReadPassword(0)
	if err != nil {
		return nil, fmt.Errorf("error reading token: %s", err)
	}
	sc = SplitterConfig{
		Username: username,
		Token:    string(b),
	}
	fmt.Println()
	err = os.Mkdir(path.Dir(splitterCachePath), 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot create directory %s: %s", path.Dir(splitterCachePath), err)
	}
	f, err := os.Create(splitterCachePath)
	if err != nil {
		return nil, fmt.Errorf("cannot write config file: %s", err)
	}
	defer ioutil.CheckClose(f, &err)
	err = json.NewEncoder(f).Encode(sc)
	if err != nil {
		return nil, fmt.Errorf("cannot write config into file: %s", err)
	}
	return &http.BasicAuth{
		Username: sc.Username,
		Password: sc.Token,
	}, nil
}

func loadCollection() (*pkg.PackageCollection, error) {
	cnf, err := conf.LoadConfig(*config, loadAuth)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %s", err)
	}

	collection, err := pkg.FromConfig(cnf)
	if err != nil {
		return nil, fmt.Errorf("error loading packages collection: %s", err)
	}

	return collection, nil
}
