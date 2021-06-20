package cmd_test

import (
	"bytes"
	"testing"

	"github.com/elhmn/ckp/cmd"
	"github.com/elhmn/ckp/mocks"
	"github.com/golang/mock/gomock"
)

func TestFindComment(t *testing.T) {
	t.Run("make sure that is actually implemented", func(t *testing.T) {
		conf := createConfig(t)
		mockedPrinters := conf.Printers.(*mocks.MockIPrinters)

		//Specify expectations
		mockedPrinters.EXPECT().SelectScriptEntry(gomock.Any())

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
