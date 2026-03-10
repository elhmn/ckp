package cmd_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/elhmn/ckp/cmd"
	"github.com/elhmn/ckp/internal/config"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

func TestSyncAliasCommand(t *testing.T) {
	t.Run("make sure that it runs successfully", func(t *testing.T) {
		conf := createConfig(t)
		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		if err := setupFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		command := cmd.NewSyncCommand(conf)
		command.SetOutput(conf.OutWriter)
		command.SetArgs([]string{"alias"})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		got := writer.String()
		exp := "\nYour aliases were successfully synced!\n"
		assert.Contains(t, got, exp)

		// Verify the aliases file was created with the correct content
		aliasesFilePath, err := config.GetAliasesFilePath(conf)
		if err != nil {
			t.Errorf("Error: failed to get aliases file path: %s", err)
		}

		aliasesContent, err := ioutil.ReadFile(aliasesFilePath)
		if err != nil {
			t.Errorf("Error: failed to read aliases file: %s", err)
		}

		// The store fixture has one entry with alias "an alias" and content `echo "mon code"`
		assert.Contains(t, string(aliasesContent), "alias an alias='echo \"mon code\"'")

		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})

	t.Run("make sure rc file is updated with source line", func(t *testing.T) {
		conf := createConfig(t)
		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		if err := setupFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		// Create a temporary rc file for this test
		home, err := homedir.Dir()
		if err != nil {
			t.Errorf("Error: failed to read home directory: %s", err)
		}

		testRCFile := fmt.Sprintf("%s/.bashrc_ckp_test", home)
		if err := ioutil.WriteFile(testRCFile, []byte("# existing content\n"), 0666); err != nil {
			t.Errorf("Error: failed to create test rc file: %s", err)
		}
		defer os.Remove(testRCFile)

		command := cmd.NewSyncCommand(conf)
		command.SetOutput(conf.OutWriter)
		command.SetArgs([]string{"alias"})

		err = command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		got := writer.String()
		exp := "\nYour aliases were successfully synced!\n"
		assert.Contains(t, got, exp)

		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})

	t.Run("make sure source line is not duplicated", func(t *testing.T) {
		conf := createConfig(t)
		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		if err := setupFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		// Run sync alias twice and make sure source line is not duplicated
		command := cmd.NewSyncCommand(conf)
		command.SetOutput(conf.OutWriter)
		command.SetArgs([]string{"alias"})

		// First run
		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		// Second run
		command2 := cmd.NewSyncCommand(conf)
		command2.SetOutput(conf.OutWriter)
		command2.SetArgs([]string{"alias"})

		err = command2.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		// Verify the aliases file was regenerated
		aliasesFilePath, err := config.GetAliasesFilePath(conf)
		if err != nil {
			t.Errorf("Error: failed to get aliases file path: %s", err)
		}

		aliasesContent, err := ioutil.ReadFile(aliasesFilePath)
		if err != nil {
			t.Errorf("Error: failed to read aliases file: %s", err)
		}

		// The aliases file should contain the header and one alias entry
		assert.Equal(t, 1, strings.Count(string(aliasesContent), "alias an alias="))

		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})
}
