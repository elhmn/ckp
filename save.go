/* ************************************************************************** */
/*                                                                            */
/*  save.go                                                                   */
/*                                                                            */
/*   By: elhmn <www.elhmn.com>                                                */
/*             <nleme@live.fr>                                                */
/*                                                                            */
/*   Created:                                                 by elhmn        */
/*   Updated: Fri Mar 08 11:13:52 2019                        by bmbarga      */
/*                                                                            */
/* ************************************************************************** */

package main

import	(
	"fmt"
	"flag"
	yaml "gopkg.in/yaml.v2"
// 	"errors"
)

type sSaveFlag struct {
	file	string
	alias	string
	comment	string
}

func	parseSaveFlags(args []string) (*sSaveFlag, *flag.FlagSet) {
	flags := &sSaveFlag{ file: "" }
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	defer fs.Parse(args[1:])

	fUsage := "get the script from a file"
	aUsage := "add an alias to your script"
	cUsage := "add a comment to your script"

	fs.StringVar(&flags.file, "file", "", fUsage)
	fs.StringVar(&flags.file, "f", "", fUsage + "(shorthand)")
	fs.StringVar(&flags.alias, "alias", "", aUsage)
	fs.StringVar(&flags.file, "a", "", aUsage + "(shorthand)")
	fs.StringVar(&flags.comment, "comment", "", cUsage)
	fs.StringVar(&flags.file, "m", "", cUsage + "(shorthand)")

	return flags, fs
}

func	saveScript(flags sSaveFlag) {
	tmp, _ := yaml.Marshal("Hey") // Debug

	fmt.Println(flags) // Debug
	fmt.Println("tmp : " + string(tmp)) // Debug
}

func	save(args []string) {
	var script string
	flags, fs := parseSaveFlags(args)
	rest := fs.Args()

	// Get script
	{
		if flags.file == "filePath" {
			if len(rest) != 1 {
				fmt.Println("Usage : save {script} ")
				return
			}
			script = rest[0]
		} else {
			//Get script from file
			script = flags.file
		}
	}

	fmt.Println("script : " + script) // Debug
	saveScript(*flags)
}
