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
	"splitter/conf"
	"splitter/pkg"
)

var (
	config = flag.String("c", "", "configfile to use")
)

func main() {
	flag.Parse()
	collection := loadCollection()
	pipeline := action.NewPipeline(collection.Conf.Actions)
	pipeline.Run(collection)
}

func loadAuth() http.AuthMethod {
	var username string
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
	fmt.Println()
	return &http.BasicAuth{
		Username: username,
		Password: string(b),
	}
}

func loadCollection() *pkg.PackageCollection {
	cnf, err := conf.LoadConfig(*config, loadAuth)
	checkError(err)
	collection, err := pkg.FromConfig(cnf)
	checkError(err)
	return collection
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
