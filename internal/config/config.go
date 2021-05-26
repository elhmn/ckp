package config

import (
	"fmt"
	"io"
	"os"

	"github.com/elhmn/ckp/internal/exec"
	"github.com/elhmn/ckp/internal/printers"
	"github.com/mitchellh/go-homedir"
)

const (
	StoreFileName     = "repo/store.yaml"
	StoreTempFileName = "repo/.temp_store.yaml"

	MainBranch = "master"
)

//Config contains the entire cli dependencies
type Config struct {
	Exec             exec.IExec
	CKPDir           string
	CKPStorageFolder string
	Spin             printers.ISpinner

	//MainBranch is a your remote repository main branch
	MainBranch string

	//WorkingBranch is `ckp` local working branch
	//instead of using your main branch `ckp` uses a separate
	//branch locally to facilite diff checks between your local
	//and remote changes
	WorkingBranch string

	//io Writers useful for testing
	OutWriter io.Writer
	ErrWriter io.Writer
}

//NewDefaultConfig creates a new default config
func NewDefaultConfig() Config {
	return Config{
		Exec:             exec.NewExec(),
		Spin:             printers.NewSpinner(),
		OutWriter:        os.Stdout,
		ErrWriter:        os.Stderr,
		CKPDir:           ".ckp",
		CKPStorageFolder: "repo",
		MainBranch:       MainBranch,
		WorkingBranch:    "working-" + MainBranch,
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

//GetStoreDirPath returns the path of the ckp store git repository
func GetStoreDirPath(conf Config) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("failed to read home directory: %s", err)
	}

	//Create ckp folder if it does not exist
	dir := fmt.Sprintf("%s/%s/%s", home, conf.CKPDir, conf.CKPStorageFolder)
	return dir, nil
}
