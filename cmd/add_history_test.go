package cmd_test

import (
	"bytes"
	"testing"

	"github.com/elhmn/ckp/cmd"
	"github.com/stretchr/testify/assert"
)

func TestAddHistoryCommand(t *testing.T) {
	t.Run("make sure that is runs successfully", func(t *testing.T) {
		conf, _ := createConfig()

		if err := setupFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		commandName := "history"
		command := cmd.NewAddCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		command.SetArgs([]string{commandName})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		got := writer.String()
		exp := "Your history was successfully added!\n"
		assert.Equal(t, exp, got)

		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})
}
