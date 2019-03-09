/* ************************************************************************** */
/*                                                                            */
/*  list.go                                                                   */
/*                                                                            */
/*   By: elhmn <www.elhmn.com>                                                */
/*             <nleme@live.fr>                                                */
/*                                                                            */
/*   Created:                                                 by elhmn        */
/*   Updated: Sat Mar 09 17:46:42 2019                        by bmbarga      */
/*                                                                            */
/* ************************************************************************** */

package main

import	(
	"fmt"
	"flag"
	"regexp"
	"io/ioutil"
	"log"
	yaml "gopkg.in/yaml.v2"
// 	"errors"
)

type sListFlag struct {
}

func	parseListFlags(args []string) (*sListFlag, *flag.FlagSet) {
	flags := &sListFlag{}
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	defer fs.Parse(args[1:])

	return flags, fs
}

func	showList(list tYaml, flags sListFlag) {
	re, err := regexp.Compile(`###(.*)###`)

	if err != nil {
		log.Fatal(err)
	}

	for id, elem :=  range list {
		fmt.Println("\033[0;33mid: " + id + "\033[0m");
		script := re.ReplaceAllString(elem.String(), `$1`)
		fmt.Println(script)
	}
}

func	listScripts(flags sListFlag) {
	storePath := ckpDir + "/" + ckpStoreFileName

	content, err := ioutil.ReadFile(storePath)
	if err != nil {
		log.Fatal(err)
		return
	}

	//Get content on tYaml map
	list := make(tYaml)
	if err := yaml.Unmarshal(content, list); err != nil {
		log.Fatal(err)
		return
	}

	showList(list, flags)
}

func	list (args []string) {
	flags, _ := parseListFlags(args)

	listScripts(*flags)
}
