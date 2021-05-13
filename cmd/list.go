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
	Will list both your first 10 code snippets and solutions

	example: ckp list --limit 20
	Will list both your first 20 code snippets and solutions

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

	command.PersistentFlags().Int64P("limit", "l", 10, `ckp list --limit 20`)
	command.PersistentFlags().BoolP("code", "c", false, `ckp list --code`)
	command.PersistentFlags().BoolP("solution", "s", false, `ckp list --solution`)

	return command
}

func listCommand(cmd *cobra.Command, args []string, conf config.Config) error {
	cmd.Flags().Parse(args)
	flags := cmd.Flags()

	//Get data from flags
	limit, err := flags.GetInt64("limit")
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

	//get store data
	storeFile, err := config.GetStoreFilePath(conf)
	if err != nil {
		return fmt.Errorf("failed to get the store file path: %s", err)
	}

	storeData, _, err := store.LoadStore(storeFile)
	if err != nil {
		return fmt.Errorf("failed to laod store: %s", err)
	}

	list := listScripts(storeData.Scripts, code, solution, limit)

	fmt.Fprintln(conf.OutWriter, list)
	return nil
}

func listScripts(scripts []store.Script, isCode, isSolution bool, limit int64) string {
	list := ""
	for _, s := range scripts {
		//if the script is a solution
		if s.Solution.Content != "" {
			if isCode {
				continue
			}
			list += fmt.Sprintf("ID: %s\n", s.ID)
			list += fmt.Sprintf("CreationTime: %s\n", s.CreationTime.Format(time.RFC1123))
			list += fmt.Sprintf("UpdateTime: %s\n", s.UpdateTime.Format(time.RFC1123))
			list += fmt.Sprintf("  Type: Solution\n")
			list += fmt.Sprintf("  Comment: %s\n", s.Comment)
			list += fmt.Sprintf("  Solution: %s\n", s.Solution.Content)
		} else {
			if isSolution {
				continue
			}
			list += fmt.Sprintf("ID: %s\n", s.ID)
			list += fmt.Sprintf("CreationTime: %s\n", s.CreationTime.Format(time.RFC1123))
			list += fmt.Sprintf("UpdateTime: %s\n", s.UpdateTime.Format(time.RFC1123))
			list += fmt.Sprintf("  Type: Code\n")
			list += fmt.Sprintf("  Alias: %s\n", s.Code.Alias)
			list += fmt.Sprintf("  Comment: %s\n", s.Comment)
			list += fmt.Sprintf("  Code: %s\n", s.Code.Content)
		}
		list += "\n"
	}
	return list
}
