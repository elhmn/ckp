/* ************************************************************************** */
/*                                                                            */
/*  save.go                                                                   */
/*                                                                            */
/*   By: elhmn <www.elhmn.com>                                                */
/*             <nleme@live.fr>                                                */
/*                                                                            */
/*   Created:                                                 by elhmn        */
/*   Updated: Sun Mar 10 08:48:48 2019                        by bmbarga      */
/*                                                                            */
/* ************************************************************************** */

package main

import (
	"flag"
	"fmt"
	"github.com/rs/xid"
	yaml "gopkg.in/yaml.v2"
	"log"
	"os"
	// 	"errors"
)

type sSaveFlag struct {
	File    string
	Alias   string
	Comment string
}

type sScript struct {
	Alias   string
	Comment string
	Script  string
}

func (sc sScript) String() string {
	return fmt.Sprintf("alias: %s\ncomment: %s\n"+
		"script: \033[0;32m%s\033[0m\n",
		sc.Alias, sc.Comment, sc.Script)
}

type tYaml map[string]sScript

func parseSaveFlags(args []string) (*sSaveFlag, *flag.FlagSet) {
	flags := &sSaveFlag{}
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	defer fs.Parse(args[1:])

	fUsage := "get the script from a file"
	aUsage := "add an alias to your script"
	cUsage := "add a comment to your script"

	fs.StringVar(&flags.File, "file", "", fUsage)
	fs.StringVar(&flags.File, "f", "", fUsage+"(shorthand)")
	fs.StringVar(&flags.Alias, "alias", "", aUsage)
	fs.StringVar(&flags.Alias, "a", "", aUsage+"(shorthand)")
	fs.StringVar(&flags.Comment, "comment", "", cUsage)
	fs.StringVar(&flags.Comment, "m", "", cUsage+"(shorthand)")

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
			Comment: flags.Comment,
			Alias:   flags.Alias,
			Script:  "###" + script + "###",
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
}

func save(args []string) {
	var script string
	flags, fs := parseSaveFlags(args)
	rest := fs.Args()

	// Get script
	{
		if flags.File == "" {
			if len(rest) != 1 {
				fmt.Println("Usage : save {script} ")
				return
			}
			script = rest[0]
		} else {
			//Get script from file
			script = flags.File
		}
	}

	saveScript(*flags, script)
}
