/* ************************************************************************** */
/*                                                                            */
/*  main.go                                                                   */
/*                                                                            */
/*   By: elhmn <www.elhmn.com>                                                */
/*             <nleme@live.fr>                                                */
/*                                                                            */
/*   Created: Sun Mar  3 17:59:45 2019                        by elhmn        */
/*   Updated: Thu Mar 07 15:30:21 2019                        by bmbarga      */
/*                                                                            */
/* ************************************************************************** */

package	main

import	(
	"fmt"
	"os"
	"errors"
)

//type for command function call
type	fCall func ([]string)

//This is not great
var		knownCommands = map[string]fCall {
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
	cpkPath = "~/.cpk"
	cpkShellrc = "~/.zshrc"
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

	call(os.Args[:1])
};
