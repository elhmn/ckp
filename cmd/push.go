package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/elhmn/ckp/internal/config"
	"github.com/spf13/cobra"
)

const (
	commitAddAction     = "add"
	commitEditAction    = "edit"
	commitRemoveAction  = "rm"
	commitDefaultAction = "push"
)

//NewPushCommand create new cobra command for the push command
func NewPushCommand(conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "push",
		Short: "push your changes to your remote repository",
		Long: `push your changed to your remote repository

		example: ckp push
`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := pushCommand(conf); err != nil {
				fmt.Fprintf(conf.OutWriter, "Error: %s\n", err)
				return
			}
		},
	}

	return command
}

func pushCommand(conf config.Config) error {
	//Setup spinner
	spin := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spin.Start()
	defer spin.Stop()

	dir, err := config.GetStoreDirPath(conf)
	if err != nil {
		return fmt.Errorf("failed get repository path: %s", err)
	}

	storeFilePath, err := config.GetStoreFilePath(conf)
	if err != nil {
		return fmt.Errorf("failed get store file path: %s", err)
	}

	spin.Suffix = " pulling remote changes..."
	err = pullRemoteChanges(conf, dir, storeFilePath)
	if err != nil {
		return fmt.Errorf("failed to pull remote changes: %s", err)
	}
	spin.Suffix = " remote changes pulled"

	spin.Suffix = " pushing local changes..."
	err = pushLocalChanges(conf, dir, storeFilePath, commitDefaultAction)
	if err != nil {
		return fmt.Errorf("failed to push local changes: %s", err)
	}
	spin.Suffix = " local changes pushed"

	fmt.Fprintf(conf.OutWriter, "\nckp store was pushed successfully\n")
	return nil
}

func pullRemoteChanges(conf config.Config, dir, file string) error {
	hasLocalChanges := false

	out, err := conf.Exec.DoGit(dir, "fetch", "origin", "main")
	if err != nil {
		return fmt.Errorf("failed to fetch origin/main: %s: %s", err, out)
	}

	out, err = conf.Exec.DoGit(dir, "diff", "origin/main", "--", file)
	if err != nil {
		return fmt.Errorf("failed to check for local changes: %s: %s", err, out)
	}

	if out != "" {
		hasLocalChanges = true
	}

	if hasLocalChanges {
		out, err = conf.Exec.DoGit(dir, "stash")
		if err != nil {
			return fmt.Errorf("failed to stash changes: %s: %s", err, out)
		}
	}

	out, err = conf.Exec.DoGit(dir, "pull", "--rebase", "origin", "main")
	if err != nil {
		return fmt.Errorf("failed to pull remote changes: %s: %s", err, out)
	}

	if hasLocalChanges {
		out, err = conf.Exec.DoGit(dir, "stash", "apply")
		if err != nil {
			return fmt.Errorf("failed to apply stash changes: %s: %s", err, out)
		}
	}
	return nil
}

func pushLocalChanges(conf config.Config, dir, file string, action string) error {
	out, err := conf.Exec.DoGit(dir, "fetch", "origin", "main")
	if err != nil {
		return fmt.Errorf("failed to fetch origin/main: %s: %s", err, out)
	}

	out, err = conf.Exec.DoGit(dir, "diff", "origin/main", "--", file)
	if err != nil {
		return fmt.Errorf("failed to check for local changes: %s: %s", err, out)
	}
	//abort if `file` does not have local changes
	if out == "" {
		return nil
	}

	out, err = conf.Exec.DoGit(dir, "add", file)
	if err != nil {
		return fmt.Errorf("failed to add changes: %s: %s", err, out)
	}

	out, err = conf.Exec.DoGit(dir, "commit", "-m", getCommitMessage(action))
	if err != nil {
		return fmt.Errorf("failed to commit changes: %s: %s", err, out)
	}

	out, err = conf.Exec.DoGitPush(dir, "origin", "main")
	if err != nil {
		return fmt.Errorf("failed to push store: %s: %s", err, out)
	}

	return nil
}

func getCommitMessage(action string) string {
	switch action {
	case commitAddAction:
		return "ckp: add entry"
	case commitEditAction:
		return "ckp: edit entry"
	case commitRemoveAction:
		return "ckp: remove entry"
	}
	return "update: update store"
}
