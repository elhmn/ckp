package files

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
)

//CopyFileToHomeDirectory copy the file at `filepath` to the home directory
func CopyFileToHomeDirectory(filepath, contentPath string) error {
	content, err := ioutil.ReadFile(contentPath)
	if err != nil {
		return fmt.Errorf("failed to read file %s data: %s", contentPath, err)
	}

	home, err := homedir.Dir()
	if err != nil {
		return fmt.Errorf("failed to read home directory: %s", err)
	}
	destination := fmt.Sprintf("%s/%s", home, filepath)

	//Copy the store file to a temporary destination
	if err := ioutil.WriteFile(destination, content, 0666); err != nil {
		return fmt.Errorf("failed to write to file %s: %s", filepath, err)
	}

	return nil
}

//DeleteFileFromHomeDirectory delete the file at `filepath` from the home directory
func DeleteFileFromHomeDirectory(filepath string) error {
	home, err := homedir.Dir()
	if err != nil {
		return fmt.Errorf("failed to read home directory: %s", err)
	}

	return os.Remove(fmt.Sprintf("%s/%s", home, filepath))
}
