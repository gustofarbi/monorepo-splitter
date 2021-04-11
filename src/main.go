package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"golang.org/x/term"
	"log"
	"os"
	"splitter/action"
	"splitter/utils/conf"
	"splitter/utils/pkg"
)

var (
	config = flag.String("conf", "", "conf file to use")
)

func main() {
	var username, password string
	fmt.Printf("enter git credentials\nusername: ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		username = scanner.Text()
		break
	}
	fmt.Print("password: ")
	b, err := term.ReadPassword(0)
	if err != nil {
		log.Fatalln("error reading password: ", err)
	}
	password = string(b)
	auth := &http.BasicAuth{
		Username: username,
		Password: password,
	}
	cnf, err := conf.LoadConfig(*config)
	checkError(err)
	collection, err := pkg.FromConfig(cnf, auth)
	checkError(err)
	pipeline := action.NewPipeline(
		//action.UpdateReplaceRelease{},
		action.SetPackagesDependencies{},
		action.WriteChanges{},
		action.TagRelease{},
		action.SplitPackages{},
	)
	pipeline.Act(collection)
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
