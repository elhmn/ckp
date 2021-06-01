package cmd

import (
	"fmt"
	"strings"
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
		Short: "pushes your changes to your remote repository",
		Long: `pushes your changed to your remote repository

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

	historyStoreFilePath, err := config.GetHistoryFilePath(conf)
	if err != nil {
		return fmt.Errorf("failed get history store file path: %s", err)
	}

	spin.Suffix = " pulling remote changes..."
	err = pullRemoteChanges(conf, dir, storeFilePath, historyStoreFilePath)
	if err != nil {
		return fmt.Errorf("failed to pull remote changes: %s", err)
	}
	spin.Suffix = " remote changes pulled"

	spin.Suffix = " pushing local changes..."
	err = pushLocalChanges(conf, dir, commitDefaultAction, storeFilePath, historyStoreFilePath)
	if err != nil {
		return fmt.Errorf("failed to push local changes: %s", err)
	}
	spin.Suffix = " local changes pushed"

	fmt.Fprintf(conf.OutWriter, "\nckp store was pushed successfully\n")
	return nil
}

func pullRemoteChanges(conf config.Config, dir string, files ...string) error {
	hasChanges := false
	hasStashed := false

	out, err := conf.Exec.DoGit(dir, "fetch", "origin", "main")
	if err != nil {
		return fmt.Errorf("failed to fetch origin/main: %s: %s", err, out)
	}

	args := append([]string{"diff", "origin/main", "--"}, files...)
	out, err = conf.Exec.DoGit(dir, args...)
	if err != nil {
		return fmt.Errorf("failed to check for local changes: %s: %s", err, out)
	}

	if out != "" {
		hasChanges = true
	}

	if hasChanges {
		out, err = conf.Exec.DoGit(dir, "stash")
		if err != nil {
			return fmt.Errorf("failed to stash changes: %s: %s", err, out)
		}

		if !strings.Contains(out, "No local changes to save") {
			hasStashed = true
		}
	}

	out, err = conf.Exec.DoGit(dir, "pull", "--rebase", "origin", "main")
	if err != nil {
		return fmt.Errorf("failed to pull remote changes: %s: %s", err, out)
	}

	if hasStashed {
		out, err = conf.Exec.DoGit(dir, "stash", "apply")
		//if there is an error and that the error is not related
		if err != nil && !strings.Contains(out, "No stash entries found") {
			return fmt.Errorf("failed to apply stash changes: %s: %s", err, out)
		}
	}
	return nil
}

func pushLocalChanges(conf config.Config, dir, action string, files ...string) error {
	out, err := conf.Exec.DoGit(dir, "fetch", "origin", "main")
	if err != nil {
		return fmt.Errorf("failed to fetch origin/main: %s: %s", err, out)
	}

	args := append([]string{"diff", "origin/main", "--"}, files...)
	out, err = conf.Exec.DoGit(dir, args...)
	if err != nil {
		return fmt.Errorf("failed to check for local changes: %s: %s", err, out)
	}
	//abort if `file` does not have changes
	if out == "" {
		return nil
	}

	args = append([]string{"add"}, files...)
	out, err = conf.Exec.DoGit(dir, args...)
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
