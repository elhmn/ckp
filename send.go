/* ************************************************************************** */
/*                                                                            */
/*  send.go                                                                   */
/*                                                                            */
/*   By: elhmn <www.elhmn.com>                                                */
/*             <nleme@live.fr>                                                */
/*                                                                            */
/*   Created: Sun Mar 10 08:45:07 2019                        by elhmn        */
/*   Updated: Sun Mar 10 09:56:18 2019                        by bmbarga      */
/*                                                                            */
/* ************************************************************************** */

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	// 	"errors"
)

type sSendFlag struct {
}

func parseSendFlags(args []string) (*sSendFlag, *flag.FlagSet) {
	flags := &sSendFlag{}
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	defer fs.Parse(args[1:])

	return flags, fs
}

func sendScript(flags sSendFlag) {
	remoteFilePath := ckpDir + "/" + ckpRemoteFileName
	repoDir := ckpDir + "/" + ckpRepoName
	content, err := ioutil.ReadFile(remoteFilePath)
	if err != nil {
		log.Fatal(err)
	}
	remote := string(content)

	cmd := exec.Command("bash", "-c", "cd "+repoDir+
		" && git add "+ckpStoreFileName+
		" && git commit -m 'Update "+ckpStoreFileName+"'"+
		" && git push origin master "+
		"&& echo "+remote)
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
	}

	slurpErr, _ := ioutil.ReadAll(stderr)
	fmt.Printf("%s", slurpErr)

	slurpOut, _ := ioutil.ReadAll(stdout)
	fmt.Printf("%s", slurpOut)

	if err := cmd.Wait(); err != nil {
		return
	}
}

func send(args []string) {
	flags, _ := parseSendFlags(args)

	sendScript(*flags)
}
