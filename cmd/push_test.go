package cmd_test

import (
	"bytes"
	"testing"

	"github.com/elhmn/ckp/cmd"
	"github.com/elhmn/ckp/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//TestPushCommand test the `ckp push` command
func TestPushCommand(t *testing.T) {
	t.Run("push successfully", func(t *testing.T) {
		conf := createConfig(t)
		mockedExec := conf.Exec.(*mocks.MockIExec)
		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		//Specify expectations
		gomock.InOrder(
			mockedExec.EXPECT().DoGit(gomock.Any(), "fetch", "origin", "main"),
			mockedExec.EXPECT().DoGit(gomock.Any(), "diff", "origin/main", "--", gomock.Any(), gomock.Any()),
			mockedExec.EXPECT().DoGit(gomock.Any(), "pull", "--rebase", "origin", "main"),
			mockedExec.EXPECT().DoGit(gomock.Any(), "fetch", "origin", "main"),
			mockedExec.EXPECT().DoGit(gomock.Any(), "diff", "origin/main", "--", gomock.Any(), gomock.Any()),
			mockedExec.EXPECT().DoGit(gomock.Any(), "add", gomock.Any(), gomock.Any()),
			mockedExec.EXPECT().DoGit(gomock.Any(), "commit", "-m", "ckp: add store"),
		)

		command := cmd.NewPushCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		exp := "\nckp store was pushed successfully\n"
		got := writer.String()
		if !assert.Equal(t, exp, got) {
			t.Errorf("expected failure with [%s], got [%s]", exp, got)
		}
	})
}
