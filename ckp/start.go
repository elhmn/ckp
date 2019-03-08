/* ************************************************************************** */
/*                                                                            */
/*  start.go                                                                  */
/*                                                                            */
/*   By: elhmn <www.elhmn.com>                                                */
/*             <nleme@live.fr>                                                */
/*                                                                            */
/*   Created:                                                 by elhmn        */
/*   Updated: Fri Mar 08 07:26:32 2019                        by bmbarga      */
/*                                                                            */
/* ************************************************************************** */

package main

import	(
	"fmt"
	"flag"
	"os/exec"
	"os"
	"log"
	"io/ioutil"
// 	"errors"
)

type sStartFlag struct {
}

func	parseStartFlags(args []string) (*sStartFlag, *flag.FlagSet) {
	flags := &sStartFlag{}
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	defer fs.Parse(args[1:])

	return flags, fs
}

func	fetchRemoteRepository(flags sStartFlag, remote string) {
	//Create ckpPath file if it does not exist
	if _, err := os.Stat(ckpDir); os.IsNotExist(err) {
		err := os.Mkdir(ckpDir, os.ModePerm)
		if err != nil {
			log.Fatal("Error : " + err.Error())
			return
		}
		fmt.Println(ckpDir + " was created")
	}

	//Move to ckpPath and clone the folder there
	cmd := exec.Command("bash", "-c",
		"cd " + ckpDir +
		" && " +
		"git clone " + remote + " " + ckpRepoName +
		" && echo " + remote + " > " + ckpRemoteFileName)

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

func	start (args []string) {
	flags, fs := parseStartFlags(args)
	rest := fs.Args()

	if len(rest) == 0 {
		fmt.Println("Usage : start {remote}")
		return
	}

	remote := rest[0]
	fetchRemoteRepository(*flags, remote)
}
