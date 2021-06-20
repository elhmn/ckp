package cmd_test

import (
	"bytes"
	"testing"

	"github.com/elhmn/ckp/cmd"
	"github.com/elhmn/ckp/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRmCommand(t *testing.T) {
	t.Run("make sure that it runs successfully", func(t *testing.T) {
		conf := createConfig(t)
		mockedExec := conf.Exec.(*mocks.MockIExec)

		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		if err := setupFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		//Specify expectations
		gomock.InOrder(
			mockedExec.EXPECT().DoGit(gomock.Any(), "fetch", "origin", "main"),
			mockedExec.EXPECT().DoGit(gomock.Any(), "diff", "origin/main", "--", gomock.Any()),
			mockedExec.EXPECT().DoGit(gomock.Any(), "pull", "--rebase", "origin", "main"),
			mockedExec.EXPECT().DoGit(gomock.Any(), "fetch", "origin", "main"),
			mockedExec.EXPECT().DoGit(gomock.Any(), "diff", "origin/main", "--", gomock.Any()),
			mockedExec.EXPECT().DoGit(gomock.Any(), "add", gomock.Any()),
			mockedExec.EXPECT().DoGit(gomock.Any(), "commit", "-m", "ckp: add store"),
		)

		command := cmd.NewRmCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		command.SetArgs([]string{"hash-of-file-content"})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		got := writer.String()
		exp := "\nentry was removed successfully\n"
		assert.Equal(t, exp, got)

		//function call assert
		if err := deleteFolder(conf); err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})
}
