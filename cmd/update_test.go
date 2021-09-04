package cmd_test

import (
	"bytes"
	"testing"

	"github.com/elhmn/ckp/cmd"
)

//TestUpdateCommand test the `ckp update` command
func TestUpdateCommand(t *testing.T) {
	t.Run("Run the update successfully", func(t *testing.T) {
		conf := createConfig(t)
		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		command := cmd.NewUpdateCommand(conf)

		//Set writer
		command.SetOutput(conf.OutWriter)

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})
}
