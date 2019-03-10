/* ************************************************************************** */
/*                                                                            */
/*  stop.go                                                                   */
/*                                                                            */
/*   By: elhmn <www.elhmn.com>                                                */
/*             <nleme@live.fr>                                                */
/*                                                                            */
/*   Created:                                                 by elhmn        */
/*   Updated: Fri Mar 08 07:37:43 2019                        by bmbarga      */
/*                                                                            */
/* ************************************************************************** */

package main


import	(
	"fmt"
	"flag"
	"os"
	"log"
// 	"errors"
)

type sStopFlag struct {
}

func	parseStopFlags(args []string) (*sStopFlag, *flag.FlagSet) {
	flags := &sStopFlag{}
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	defer fs.Parse(args[1:])

	return flags, fs
}

func	clearCpk(flags sStopFlag) {
	if _, err := os.Stat(ckpDir); err == nil {
		err := os.RemoveAll(ckpDir)
		if err != nil {
			log.Fatal("Error : " + err.Error())
		}
		fmt.Println(ckpDir + " was removed")
	}
	//remove aliases from .zshrc
	//And cleanup other things
}

func	stop (args []string) {
	flags, _ := parseStopFlags(args)

	clearCpk(*flags)
}
