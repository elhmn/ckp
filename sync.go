/* ************************************************************************** */
/*                                                                            */
/*  sync.go                                                                   */
/*                                                                            */
/*   By: elhmn <www.elhmn.com>                                                */
/*             <nleme@live.fr>                                                */
/*                                                                            */
/*   Created:                                                 by elhmn        */
/*   Updated: Sat Mar 09 18:14:21 2019                        by bmbarga      */
/*                                                                            */
/* ************************************************************************** */

package main

import	(
	"fmt"
	"flag"
	"os/exec"
// 	"os"
	"log"
	"io/ioutil"
	yaml "gopkg.in/yaml.v2"
// 	"errors"
)

type sSyncFlag struct {
	alias	string
}

func	parseSyncFlags(args []string) (*sSyncFlag, *flag.FlagSet) {
	flags := &sSyncFlag{}
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	defer fs.Parse(args[1:])

	aUsage := "Sync script that has a specific alias"

	fs.StringVar(&flags.alias, "alias", "", aUsage)
	fs.StringVar(&flags.alias, "a", "", aUsage + "(shorthand)")
	return flags, fs
}

func	syncCommand(flags sSyncFlag) {
	list := make(tYaml)

	//Get script from yaml file
	{
		storePath := ckpDir + "/" + ckpStoreFileName

		content, err := ioutil.ReadFile(storePath)
		if err != nil {
			log.Fatal(err)
			return
		}

		//Get content on tYaml map
		if err := yaml.Unmarshal(content, list); err != nil {
			log.Fatal(err)
			return
		}
	}

	//Check if an alias exist in the yaml
	//Get bash zsh sh files
	//Add alias to user local zshrc bashrc or shrc
	// if it does not exist in the file

	//Move to ckpPath and clone the folder there
	cmd := exec.Command("bash", "-c", "echo I sync your script !")
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
		return
	}

	slurpErr, _ := ioutil.ReadAll(stderr)
	fmt.Printf("%s", slurpErr)

	slurpOut, _ := ioutil.ReadAll(stdout)
	fmt.Printf("%s", slurpOut)

	if err := cmd.Wait(); err != nil {
		return
	}
}

func	sync (args []string) {
	flags, _ := parseSyncFlags(args)

	syncCommand(*flags)
}
