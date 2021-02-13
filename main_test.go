package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/elhmn/ckp/cmd"
)

type MockedExec struct {
	RunErrorOutput error
	RunOutput      []byte

	DoGitCloneErrorOutput error
	DoGitCloneOutput      string

	CreateFolderIfDoesNotExistErrorOutput error
}

func (ex MockedExec) Run(dir string, command string, args ...string) ([]byte, error) {
	return ex.RunOutput, ex.RunErrorOutput
}

func (ex MockedExec) DoGitClone(dir string, args ...string) (string, error) {
	return ex.DoGitCloneOutput, ex.DoGitCloneErrorOutput
}

func (ex MockedExec) CreateFolderIfDoesNotExist(dir string) error {
	return ex.CreateFolderIfDoesNotExistErrorOutput
}

//TestInitCommand test the `ckp init` command
func TestInitCommand(t *testing.T) {
	fakeRemoteFolder := "https://github.com/elhmn/fakefolder"
	initCommand := "init"

	t.Run("initialised successfully", func(t *testing.T) {
		conf := newConfig()
		conf.Exec = &MockedExec{}
		writer := &bytes.Buffer{}
		conf.OutWriter = writer

		command := cmd.NewCKPCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		command.SetArgs([]string{initCommand, fakeRemoteFolder})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}
	})

	t.Run("failed to create folder", func(t *testing.T) {
		conf := newConfig()
		writer := &bytes.Buffer{}
		conf.OutWriter = writer
		exp := "failed to create folder"

		//Setup for failure
		conf.Exec = &MockedExec{
			CreateFolderIfDoesNotExistErrorOutput: fmt.Errorf(exp),
		}

		command := cmd.NewCKPCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		command.SetArgs([]string{initCommand, fakeRemoteFolder})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		got := writer.String()
		if !strings.Contains(got, exp) {
			t.Errorf("expected failure with [%s], got [%s]", exp, got)
		}
	})

	t.Run("failed to clone remote repository", func(t *testing.T) {
		conf := newConfig()
		writer := &bytes.Buffer{}
		conf.OutWriter = writer
		exp := "failed to clone remote repository"

		//Setup for failure
		conf.Exec = &MockedExec{
			DoGitCloneErrorOutput: fmt.Errorf(exp),
		}

		command := cmd.NewCKPCommand(conf)
		//Set writer
		command.SetOutput(conf.OutWriter)

		//Set args
		command.SetArgs([]string{initCommand, fakeRemoteFolder})

		err := command.Execute()
		if err != nil {
			t.Errorf("Error: failed with %s", err)
		}

		got := writer.String()
		if !strings.Contains(got, exp) {
			t.Errorf("expected failure with [%s], got [%s]", exp, got)
		}
	})
}
