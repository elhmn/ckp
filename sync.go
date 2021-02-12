package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type sSyncFlag struct {
	alias string
}

func parseSyncFlags(args []string) (*sSyncFlag, *flag.FlagSet) {
	flags := &sSyncFlag{}
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	defer fs.Parse(args[1:])

	aUsage := "Sync script that has a specific alias"

	fs.StringVar(&flags.alias, "alias", "", aUsage)
	fs.StringVar(&flags.alias, "a", "", aUsage+"(shorthand)")
	return flags, fs
}

func syncCommand(flags sSyncFlag) {
	list := make(tYaml)

	//Get script from yaml file
	{
		content, err := ioutil.ReadFile(ckpStorePath)
		if err != nil {
			log.Fatal(err)
		}

		//Get content on tYaml map
		if err := yaml.Unmarshal(content, list); err != nil {
			log.Fatal(err)
		}
	}

	//Check if an alias exist in the yaml
	{
		//Open ckpAliasFile
		aliasFilePath := ckpDir + "/" + ckpAliasFile
		aliasFile, err := os.OpenFile(aliasFilePath,
			os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer aliasFile.Close()

		//Get lines
		lines := []string{}
		{
			scanner := bufio.NewScanner(aliasFile)
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		}

	AliasLoop:
		for id, elem := range list {
			if elem.Alias != "" {
				//Check if an alias already exist
				{
					for _, line := range lines {
						if strings.Contains(line, elem.Alias) {
							fmt.Printf("\033[0;32m%s\033[0m already exist in %s\n", elem.Alias, aliasFilePath)
							continue AliasLoop
						}
					}
				}

				//Add alias to ckpAliasFile
				{
					re, err := regexp.Compile(`###(.*)###`)
					if err != nil {
						log.Fatal(err)
					}
					script := re.ReplaceAllString(elem.Script, `$1`)

					if _, err := aliasFile.WriteString("alias " + elem.Alias + "='" + script + "'\n"); err != nil {
						log.Fatal(err)
					}
				}

				fmt.Printf("\033[0;33m%s\033[0m : \033[0;32m%s\033[0m : was added in %s !\n", id, elem.Alias, aliasFilePath) // Debug
			}
		}

		//Get bash zsh sh files
	RcFileLoop:
		for _, rc := range ckpRcFiles {
			rcFilePath := ckpUsr.HomeDir + "/" + rc
			source := "source " + aliasFilePath
			fmt.Println(rcFilePath) // Debug

			rcFile, err := os.OpenFile(rcFilePath,
				os.O_RDWR|os.O_APPEND, 0644)
			if err != nil {
				fmt.Println(err)
				continue RcFileLoop
			}

			scanner := bufio.NewScanner(rcFile)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(line, source) {
					rcFile.Close()
					fmt.Printf("Run '\033[0;33msource %s' \033[0m to update your shell\n", rcFilePath)
					continue RcFileLoop
				}
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}

			if _, err := rcFile.WriteString("source " + aliasFilePath + "\n"); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("'%s' added to %s\n", source, rcFilePath)
			fmt.Printf("Run '\033[0;33msource %s' \033[0m to update your shell", rcFilePath)
			rcFile.Close()
		}
	}
}

func sync(args []string) {
	flags, _ := parseSyncFlags(args)

	syncCommand(*flags)
}
