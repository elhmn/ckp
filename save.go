/* ************************************************************************** */
/*                                                                            */
/*  save.go                                                                   */
/*                                                                            */
/*   By: elhmn <www.elhmn.com>                                                */
/*             <nleme@live.fr>                                                */
/*                                                                            */
/*   Created:                                                 by elhmn        */
/*   Updated: Wed Oct 09 12:10:18 2019                        by bmbarga      */
/*                                                                            */
/* ************************************************************************** */

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/rs/xid"
	yaml "gopkg.in/yaml.v2"
	// 	"errors"
)

type sSaveFlag struct {
	File    string
	Alias   string
	Comment string
	Group   string
}

type sScript struct {
	File         string
	Alias        string
	Comment      string
	Script       string
	Group        string
	CreationTime string
}

func (sc sScript) String() string {
	return fmt.Sprintf("alias: %s\ncomment: %s\n"+
		"script: \033[0;32m%s\033[0m\n"+
		"file: \033[0;32m%s\033[0m\n"+
		"group: \033[0;32m%s\033[0m\n"+
		"creationTime: \033[0;32m%s\033[0m\n",
		sc.Alias, sc.Comment, sc.Script, sc.File, sc.Group, sc.CreationTime)
}

type tYaml map[string]sScript

func parseSaveFlags(args []string) (*sSaveFlag, *flag.FlagSet) {
	flags := &sSaveFlag{}
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	defer fs.Parse(args[1:])

	fUsage := "get the script from a file"
	aUsage := "add an alias to your script"
	cUsage := "add a comment to your script"
	gUsage := "specify a group for your script"

	fs.StringVar(&flags.File, "file", "", fUsage)
	fs.StringVar(&flags.File, "f", "", fUsage+"(shorthand)")
	fs.StringVar(&flags.Alias, "alias", "", aUsage)
	fs.StringVar(&flags.Alias, "a", "", aUsage+"(shorthand)")
	fs.StringVar(&flags.Comment, "comment", "", cUsage)
	fs.StringVar(&flags.Comment, "m", "", cUsage+"(shorthand)")
	fs.StringVar(&flags.Group, "group", "", gUsage)
	fs.StringVar(&flags.Group, "g", "", gUsage+"(shorthand)")

	return flags, fs
}

func saveScript(flags sSaveFlag, script string) {

	file, err := os.OpenFile(ckpStorePath,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//Create a unique id for the script
	guid := xid.New()

	//Create the yaml structure to store the script
	yml := make(tYaml)
	{
		yml[guid.String()] = sScript{
			CreationTime: time.Now().String(),
			Comment:      flags.Comment,
			Alias:        flags.Alias,
			File:         flags.File,
			Group:        flags.Group,
			Script:       "###" + script + "###",
		}
	}

	//Append script yaml to the store file
	scriptYaml, _ := yaml.Marshal(yml)
	{
		if _, err := file.WriteString("#script ===\n" + string(scriptYaml)); err != nil {
			log.Fatal(err)
		}
	}

	if err := file.Close(); err != nil {
		log.Fatal(err)
	}

	//Show successfull message (must be replaced)
	fmt.Printf("New code entry successfully saved:\n\n")
	re, err := regexp.Compile(`###(.*)###`)
	fmt.Println("\033[0;33mid: " + guid.String() + "\033[0m")
	output := re.ReplaceAllString(yml[guid.String()].String(), `$1`)
	fmt.Println(output)
}

func save(args []string) {
	script := ""
	flags, fs := parseSaveFlags(args)
	rest := fs.Args()
	if len(rest) >= 1 {
		script = rest[0]
	}
	saveScript(*flags, script)
}
