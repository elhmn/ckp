package exec

import (
	"os"
	exe "os/exec"
)

//IExec defines exec global interface, useful for testing
type IExec interface {
	Run(dir string, command string, args ...string) ([]byte, error)
	RunInteractive(command string, args ...string) error
	DoGitClone(dir string, args ...string) (string, error)
	DoGitPush(dir string, args ...string) (string, error)
	DoGit(dir string, args ...string) (string, error)
	CreateFolderIfDoesNotExist(dir string) error
	OpenEditor(editor string, args ...string) error
}

//Exec struct
type Exec struct{}

//NewExec returns a new Exec
func NewExec() Exec {
	return Exec{}
}

//CreateFolderIfDoesNotExist checks, will check that a folder exist and create the folder if it does not exist
func (ex Exec) CreateFolderIfDoesNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

//Run run command and return output
func (ex Exec) Run(dir string, command string, args ...string) ([]byte, error) {
	cmd := exe.Command(command, args...)
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

//RunInteractive run the command in interactive mode
func (ex Exec) RunInteractive(command string, args ...string) error {
	cmd := exe.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
