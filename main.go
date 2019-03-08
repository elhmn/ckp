/* ************************************************************************** */
/*                                                                            */
/*  main.go                                                                   */
/*                                                                            */
/*   By: elhmn <www.elhmn.com>                                                */
/*             <nleme@live.fr>                                                */
/*                                                                            */
/*   Created: Sun Mar  3 17:59:45 2019                        by elhmn        */
/*   Updated: Fri Mar 08 12:45:51 2019                        by bmbarga      */
/*                                                                            */
/* ************************************************************************** */

package	main

import	(
	"fmt"
	"os"
	"os/user"
	"errors"
)

//type for command function call
type	fCall func ([]string)

//This is not great
var		knownCommands = map[string]fCall {
	"start": start,
	"stop": stop,
	"save": save,
	"run": run,
	"help": func ([]string) { },
	"sync": func ([]string) { } ,
	"send": func ([]string) { },
	"fetch": func ([]string) { },
	"list": func ([]string) { },
}

//Environment variables
var (
	usr, _ = user.Current()
	ckpDir = usr.HomeDir + "/.ckp"
	ckpRepoName = "repo"
	ckpRemoteFileName = "remote"
	ckpStoreFileName = "store.ckp"
	ckpShellrc = usr.HomeDir + "/.zshrc"
)

func	getCommandCall(args []string) (fCall, error) {
	programName := args[0]

	if l := len(args); l < 2 {
		return nil, errors.New(programName + " must have at least 2 arguments")
	}

	cmdName := args[1]
	call, ok := knownCommands[cmdName]
	if !ok {
		return nil, errors.New(programName + " can't handle command " + cmdName)
	}

	return call, nil
}

func	main() {
	call, err := getCommandCall(os.Args)

	if err != nil {
		fmt.Println("Failed to parse command")
		return
	}

	call(os.Args[1:])
};
