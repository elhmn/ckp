package main

import (
	"errors"
	"fmt"
	"os"
	"os/user"
)

//type for command function call
type fCall func([]string)

//This is not great
var knownCommands = map[string]fCall{
	"start": start,
	"stop":  stop,
	"save":  save,
	"sync":  sync,
	"send":  send,
	"list":  list,
	"help":  help,
}

func help(args []string) {
	fmt.Println(ckpUsage)
}

//Environment variables
var (
	ckpUsr, _         = user.Current()
	ckpDir            = ckpUsr.HomeDir + "/.ckp"
	ckpRepoName       = "repo"
	ckpRemoteFileName = "remote"
	ckpStoreFileName  = "store.ckp"
	ckpAliasFile      = "ckp_aliases"
	ckpShellrc        = ckpUsr.HomeDir + "/.zshrc"
	ckpStorePath      = ckpDir + "/" + ckpRepoName + "/" + ckpStoreFileName
	ckpRcFiles        = []string{
		".zshrc",
		".shrc",
		".bashrc",
	}
	ckpUsage = `usage: ckp help commands

A tool to manage your scripts.

positional arguments:
	{save,start,stop,sync,list}
        start		Clone your remote solution repoitory and init ckp
	stop		Remove ckp instance
	save		Save your scripts locally
	send		Send your local scripts to a remote server
	sync		Add your aliased scripts to your local .rc file
	list		List local scripts`
)

func getCommandCall(args []string) (fCall, error) {
	programName := args[0]

	if l := len(args); l < 2 {
		return nil, errors.New(ckpUsage)
	}

	cmdName := args[1]
	call, ok := knownCommands[cmdName]
	if !ok {
		return nil, errors.New(programName + " can't handle command " + cmdName)
	}

	return call, nil
}

func main() {
	call, err := getCommandCall(os.Args)

	if err != nil {
		fmt.Println(err)
		return
	}

	call(os.Args[1:])
}
