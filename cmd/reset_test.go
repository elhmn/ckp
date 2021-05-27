package cmd_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/elhmn/ckp/cmd"
	"github.com/stretchr/testify/assert"
)

func TestResetommand(t *testing.T) {
	t.Run("make sure that it runs successfully", func(t *testing.T) {
		conf, _ := createConfig()

		if err := setupFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		commandName := "reset"
		command := cmd.NewResetCommand(conf)

		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		command.SetArgs([]string{commandName})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		//Get remote storage folder
		folder, err := getTempStorageFolder(conf)
		if err != nil {
			t.Errorf("failed to get temporary storage folder: %s: %s", folder, err)
		}

		//Assert that the command deleted the remote storage folder
		if _, err := os.Stat(folder); !os.IsNotExist(err) {
			t.Errorf("Failed to remove %s folder : %s", folder, err)
		}

		got := writer.String()
		exp := "ckp was successfully reset\n"
		assert.Equal(t, exp, got)

		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})
}
