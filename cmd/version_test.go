package cmd_test

import (
	"bytes"
	"testing"

	"github.com/elhmn/ckp/cmd"
	"github.com/stretchr/testify/assert"
)

//TestVersionCommand test the `ckp version` command
func TestVersionCommand(t *testing.T) {
	t.Run("showed version successfully", func(t *testing.T) {
		conf := createConfig(t)
		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		command := cmd.NewVersionCommand(conf)

		//Set writer
		command.SetOutput(conf.OutWriter)

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		got := writer.String()
		exp := "Version: 0.0.0+dev\nBuild by elhmn\nSupport osscameroon here https://opencollective.com/osscameroon\n"
		assert.Contains(t, got, exp)
	})
}
