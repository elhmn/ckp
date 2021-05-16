package cmd_test

import (
	"bytes"
	"testing"

	"github.com/elhmn/ckp/cmd"
)

func TestFindComment(t *testing.T) {
	t.Run("make sure that is actually implemented", func(t *testing.T) {
		conf := createConfig()

		//setup temporary folder
		setupFolder(conf)
		defer deleteFolder(conf)

		writer := &bytes.Buffer{}
		conf.OutWriter = writer
		command := cmd.NewFindCommand(conf)

		//Set writer
		command.SetOutput(conf.OutWriter)

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})
}
