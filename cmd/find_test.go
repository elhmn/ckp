package cmd_test

import (
	"bytes"
	"testing"

	"github.com/elhmn/ckp/cmd"
)

func TestFindComment(t *testing.T) {
	t.Run("make sure that is actually implemented", func(t *testing.T) {
		conf, _ := createConfig()

		//setup temporary folder
		if err := setupFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		writer := &bytes.Buffer{}
		conf.OutWriter = writer
		command := cmd.NewFindCommand(conf)

		//Set writer
		command.SetOutput(conf.OutWriter)

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})
}
