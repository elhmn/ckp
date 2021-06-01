package cmd_test

import (
	"bytes"
	"testing"

	"github.com/elhmn/ckp/cmd"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//TestPullCommand test the `ckp pull` command
func TestPullCommand(t *testing.T) {
	t.Run("pull successfully", func(t *testing.T) {
		conf, mockedExec := createConfig(t)
		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		//Specify expectations
		gomock.InOrder(
			mockedExec.EXPECT().DoGit(gomock.Any(), "fetch", "origin", "main"),
			mockedExec.EXPECT().DoGit(gomock.Any(), "diff", "origin/main", "--", gomock.Any(), gomock.Any()),
			mockedExec.EXPECT().DoGit(gomock.Any(), "pull", "--rebase", "origin", "main"),
		)

		command := cmd.NewPullCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		exp := "\nckp store was pulled successfully\n"
		got := writer.String()
		if !assert.Equal(t, exp, got) {
			t.Errorf("expected failure with [%s], got [%s]", exp, got)
		}
	})
}
