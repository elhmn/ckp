package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"sort"

	yaml "gopkg.in/yaml.v2"
	// 	"errors"
)

type sListFlag struct {
}

func parseListFlags(args []string) (*sListFlag, *flag.FlagSet) {
	flags := &sListFlag{}
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	defer fs.Parse(args[1:])

	return flags, fs
}

func showList(list tYaml, flags sListFlag) {
	re, err := regexp.Compile(`###(.*)###`)

	if err != nil {
		log.Fatal(err)
	}

	var ids []string
	for id := range list {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		fmt.Println("\033[0;33mid: " + id + "\033[0m")
		script := re.ReplaceAllString(list[id].String(), `$1`)
		fmt.Println(script)
	}
}

func listScripts(flags sListFlag) {
	content, err := ioutil.ReadFile(ckpStorePath)
	if err != nil {
		log.Fatal(err)
	}

	//Get content on tYaml map
	list := make(tYaml)
	if err := yaml.Unmarshal(content, list); err != nil {
		log.Fatal(err)
	}

	showList(list, flags)
}

func list(args []string) {
	flags, _ := parseListFlags(args)

	listScripts(*flags)
}
