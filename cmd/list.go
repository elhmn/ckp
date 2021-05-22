package cmd

import (
	"fmt"
	"time"

	"github.com/elhmn/ckp/internal/config"
	"github.com/elhmn/ckp/internal/store"
	"github.com/spf13/cobra"
)

//NewListCommand stores everything that written after --code or --solution flag
func NewListCommand(conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "list will display your code snippets and solutions",
		Long: `list will display the code snippets and solutions you have stored

	example: ckp list
	Will list your first 10 code snippets and solutions

	example: ckp list --limit 20
	Will list your first 20 code snippets and solutions

	example: ckp list --all
	Will list all your code snippets and solutions

	example: ckp list --code
	Will list your first 10 code snippets only

	example: ckp list --solution
	Will list your 10 first solutions only
`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := listCommand(cmd, args, conf); err != nil {
				fmt.Fprintf(conf.OutWriter, "Error: %s\n", err)
				return
			}
		},
	}

	command.PersistentFlags().IntP("limit", "l", 10, `limit the number of element listed`)
	command.PersistentFlags().BoolP("code", "c", false, `list your code records only`)
	command.PersistentFlags().BoolP("solution", "s", false, `list your solutions only`)
	command.PersistentFlags().BoolP("all", "a", false, `list all your code and solutions`)

	return command
}

func listCommand(cmd *cobra.Command, args []string, conf config.Config) error {
	cmd.Flags().Parse(args)
	flags := cmd.Flags()

	//Get data from flags
	limit, err := flags.GetInt("limit")
	if err != nil {
		return fmt.Errorf("could not parse `limit` flag: %s", err)
	}
	code, err := flags.GetBool("code")
	if err != nil {
		return fmt.Errorf("could not parse `code` flag: %s", err)
	}
	solution, err := flags.GetBool("solution")
	if err != nil {
		return fmt.Errorf("could not parse `solution` flag: %s", err)
	}
	all, err := flags.GetBool("all")
	if err != nil {
		return fmt.Errorf("could not parse `all` flag: %s", err)
	}

	//get store data
	storeFile, err := config.GetStoreFilePath(conf)
	if err != nil {
		return fmt.Errorf("failed to get the store file path: %s", err)
	}

	storeData, _, err := store.LoadStore(storeFile)
	if err != nil {
		return fmt.Errorf("failed to laod store: %s", err)
	}

	list := listScripts(storeData.Scripts, code, solution, all, limit)

	fmt.Fprintln(conf.OutWriter, list)
	return nil
}

func getField(field, value string) string {
	if value != "" {
		return fmt.Sprintf("%s: %s\n", field, value)
	}
	return ""
}

func listScripts(scripts []store.Script, isCode, isSolution, shouldListAll bool, limit int) string {
	list := ""
	size := len(scripts)

	//if --all was specified set the limit to the size of the list of scripts
	if shouldListAll {
		limit = size
	}

	for i := 0; i < limit && i < size; i++ {
		s := scripts[i]
		//if the script is a solution
		if s.Solution.Content != "" {
			if isCode {
				continue
			}
			list += getField("ID", s.ID)
			list += getField("CreationTime", s.CreationTime.Format(time.RFC1123))
			list += getField("UpdateTime", s.UpdateTime.Format(time.RFC1123))
			list += fmt.Sprintf("  Type: Solution\n")
			list += getField("  Comment", s.Comment)
			list += getField("  Solution", s.Solution.Content)
		} else {
			if isSolution {
				continue
			}
			list += getField("ID", s.ID)
			list += getField("CreationTime", s.CreationTime.Format(time.RFC1123))
			list += getField("UpdateTime", s.UpdateTime.Format(time.RFC1123))
			list += fmt.Sprintf("  Type: Code\n")
			list += getField("  Alias", s.Code.Alias)
			list += getField("  Comment", s.Comment)
			list += getField("  Code", s.Code.Content)
		}
		list += "\n"
	}
	return list
}
