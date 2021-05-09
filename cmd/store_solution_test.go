package cmd_test

import (
	"bytes"
	"testing"

	"github.com/elhmn/ckp/cmd"
	"github.com/stretchr/testify/assert"
)

func TestStoreSolutionCommand(t *testing.T) {
	t.Run("make sure that is runs successfully", func(t *testing.T) {
		conf := createConfig()
		setupFolder(conf)
		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		commandName := "solution"
		command := cmd.NewStoreCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		command.SetArgs([]string{commandName,
			"our solution",
			"--path", "filepath",
			"--comment", "a_comment",
		})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		got := writer.String()
		exp := "Your code was successfully stored!\n"
		assert.Equal(t, exp, got)

		deleteFolder(conf)
	})
}
