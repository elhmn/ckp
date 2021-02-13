package exec

import (
	"os"
)

//IExec defines exec global interface, useful for testing
type IExec interface {
	Run(dir string, command string, args ...string) ([]byte, error)
	DoGitClone(dir string, args ...string) (string, error)
	CreateFolderIfDoesNotExist(dir string) error
}

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
