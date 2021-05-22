package cmd_test

import (
	"bytes"
	"testing"

	"github.com/elhmn/ckp/cmd"
)

func TestListCommand(t *testing.T) {
	t.Run("make sure that is runs successfully with limit 12", func(t *testing.T) {
		conf := createConfig()
		if err := setupFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		command := cmd.NewListCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		command.SetArgs([]string{"-l", "12"})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})

	t.Run("make sure that is runs successfully with --all flag set", func(t *testing.T) {
		conf := createConfig()
		setupFolder(conf)
		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		command := cmd.NewListCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		command.SetArgs([]string{"--all"})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		deleteFolder(conf)
	})
}
