package cmd_test

import (
	"bytes"
	"testing"

	"github.com/elhmn/ckp/cmd"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRunCommand(t *testing.T) {
	t.Run("make sure that it runs successfully", func(t *testing.T) {
		conf, mockedExec := createConfig(t)
		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		if err := setupFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		//Specify expectations
		gomock.InOrder(
			mockedExec.EXPECT().RunInteractive("bash", "-c", "echo \"mon code\"\n"),
		)

		command := cmd.NewRunCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		command.SetArgs([]string{"hash-of-file-content"})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		//function call assert
		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})

	t.Run("fail with solution id", func(t *testing.T) {
		conf, mockedExec := createConfig(t)
		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		if err := setupFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		//Specify expectations
		gomock.InOrder(
			mockedExec.EXPECT().RunInteractive("bash", "-c", "echo \"mon code\"\n").Times(0),
		)

		command := cmd.NewRunCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		command.SetArgs([]string{"hash-of-file-content-2"})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		got := writer.String()
		exp := "might not be a code entry, nothing to run"
		assert.Contains(t, got, exp)

		//function call assert
		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})
}
