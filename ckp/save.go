/* ************************************************************************** */
/*                                                                            */
/*  save.go                                                                   */
/*                                                                            */
/*   By: elhmn <www.elhmn.com>                                                */
/*             <nleme@live.fr>                                                */
/*                                                                            */
/*   Created:                                                 by elhmn        */
/*   Updated: Thu Mar 07 19:06:46 2019                        by bmbarga      */
/*                                                                            */
/* ************************************************************************** */

package main

import	(
	"fmt"
	"flag"
// 	"errors"
)

func	parseSaveFlags(args []string) {
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	var tmp string
	fs.StringVar(&tmp, "file", "filePath", "-file=filePath")
	fs.Parse(args[0:])
	fmt.Println("tmp : ", tmp);
}

func	save(args []string) {
	parseSaveFlags(args)
}
