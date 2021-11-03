package cmd_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/elhmn/ckp/cmd"
	"github.com/elhmn/ckp/internal/config"
)

func TestListCommand(t *testing.T) {
	t.Run("make sure that is runs successfully with limit 12", func(t *testing.T) {
		conf := createConfig(t)
		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		if err := setupFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}

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
		conf := createConfig(t)
		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		if err := setupFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		command := cmd.NewListCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		command.SetArgs([]string{"--all"})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})

	t.Run("make sure that is runs successfully on history, with --all flag set", func(t *testing.T) {
		conf := createConfig(t)
		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		if err := setupFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		command := cmd.NewListCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		command.SetArgs([]string{"--all", "--from-history"})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})
}

func BenchmarkListFromHistory(b *testing.B) {
	conf := config.NewDefaultConfig(config.Options{Version: "0.0.0+dev"})
	writer := &bytes.Buffer{}
	conf.OutWriter = writer
	conf.CKPDir = ".ckp_test"

	if err := setupFolder(conf); err != nil {
		fmt.Printf("Error: failed with %s", err)
	}

	command := cmd.NewListCommand(conf)
	//Set writer
	command.SetOutput(conf.OutWriter)

	//Set args
	command.SetArgs([]string{"--all", "--from-history"})

	err := command.Execute()
	if err != nil {
		fmt.Printf("Error: failed with %s", err)
	}

	if err := deleteFolder(conf); err != nil {
		fmt.Printf("Error: failed with %s", err)
	}
}
