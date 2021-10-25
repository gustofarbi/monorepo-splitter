package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/utils/ioutil"
	"golang.org/x/term"
	"log"
	"os"
	"path"
	"splitter/action"
	"splitter/conf"
	"splitter/pkg"
)

var (
	config    = flag.String("c", "", "configfile to use")
	overwrite = flag.Bool("o", false, "overwrite configs cache")
)

type SplitterConfig struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func main() {
	flag.Parse()
	collection := loadCollection()
	pipeline := action.NewPipeline(collection.Conf.Actions)
	pipeline.Run(collection)
}

func loadAuth() http.AuthMethod {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("cannot get home-dir: %+v\n", err)
	}
	var sc SplitterConfig
	splitterCachePath := path.Join(homeDir, ".splitter", "config")
	if _, err := os.Stat(splitterCachePath); err == nil || *overwrite {
		var errCacheFile error
		bytes, errCacheFile := os.ReadFile(splitterCachePath)
		if errCacheFile != nil {
			fmt.Println("cannot read config cache file")
		} else if errCacheFile = json.Unmarshal(bytes, &sc); errCacheFile != nil {
			fmt.Println("cannot unmarshal config cache file")
		}
		if errCacheFile != nil {
			_ = os.Remove(splitterCachePath)
		} else {
			return &http.BasicAuth{
				Username: sc.Username,
				Password: sc.Token,
			}
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
		log.Fatalln("error reading token: ", err)
	}
	sc = SplitterConfig{
		Username: username,
		Token:    string(b),
	}
	fmt.Println()
	err = os.Mkdir(path.Dir(splitterCachePath), 0755)
	if err != nil {
		fmt.Printf("cannot create directory %s: %+v\n", path.Dir(splitterCachePath), err)
	}
	f, err := os.Create(splitterCachePath)
	if err != nil {
		log.Printf("cannot write config file: %+v\n", err)
	}
	defer ioutil.CheckClose(f, &err)
	err = json.NewEncoder(f).Encode(sc)
	if err != nil {
		log.Printf("cannot write config into file: %+v\n", err)
	}
	return &http.BasicAuth{
		Username: sc.Username,
		Password: sc.Token,
	}
}

func loadCollection() *pkg.PackageCollection {
	cnf, err := conf.LoadConfig(*config, loadAuth)
	if err != nil {
		log.Fatalf("error loading config: %+v", err)
	}

	collection, err := pkg.FromConfig(cnf)
	if err != nil {
		log.Fatalf("error loading packages collection: %+v", err)
	}

	return collection
}
