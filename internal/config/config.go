package config

import (
	"fmt"
	"io"
	"os"

	"github.com/elhmn/ckp/internal/exec"
	"github.com/mitchellh/go-homedir"
)

const (
	StoreFileName     = "repo/store.yaml"
	StoreTempFileName = "repo/.temp_store.yaml"
)

//Config contains the entire cli dependencies
type Config struct {
	Exec             exec.IExec
	CKPDir           string
	CKPStorageFolder string

	//io Writers useful for testing
	OutWriter io.Writer
	ErrWriter io.Writer
}

//NewDefaultConfig creates a new default config
func NewDefaultConfig() Config {
	return Config{
		Exec:             exec.NewExec(),
		OutWriter:        os.Stdout,
		ErrWriter:        os.Stderr,
		CKPDir:           ".ckp",
		CKPStorageFolder: "repo",
	}
}

//GetStoreFilePath get the store file path from config
func GetStoreFilePath(conf Config) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("failed to read home directory: %s", err)
	}

	storepath := fmt.Sprintf("%s/%s/%s", home, conf.CKPDir, StoreFileName)
	return storepath, nil
}

//GetTempStoreFilePath get the temporary store file path from config
func GetTempStoreFilePath(conf Config) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("failed to read home directory: %s", err)
	}

	storepath := fmt.Sprintf("%s/%s/%s", home, conf.CKPDir, StoreTempFileName)
	return storepath, nil
}
