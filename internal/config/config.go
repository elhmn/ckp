package config

import (
	"fmt"
	"io"
	"os"

	"github.com/elhmn/ckp/internal/exec"
	"github.com/elhmn/ckp/internal/printers"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	StoreFileName            = "repo/store.yaml"
	HistoryFileName          = "repo/history_store.yaml"
	StoreTempFileName        = "repo/.temp_store.yaml"
	StoreHistoryTempFileName = "repo/.temp_history_store.yaml"

	MainBranch = "main"
)

//Config contains the entire cli dependencies
type Config struct {
	Version          string
	Viper            viper.Viper
	Exec             exec.IExec
	CKPDir           string
	CKPStorageFolder string
	Spin             printers.ISpinner
	Printers         printers.IPrinters

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

//Options config options
type Options struct {
	Version string
}

//NewDefaultConfig creates a new default config
func NewDefaultConfig(opt Options) Config {
	conf := Config{
		Exec:             exec.NewExec(),
		Spin:             printers.NewSpinner(),
		Printers:         printers.NewPrinters(),
		OutWriter:        os.Stdout,
		ErrWriter:        os.Stderr,
		CKPDir:           ".ckp",
		CKPStorageFolder: "repo",
		MainBranch:       MainBranch,
		WorkingBranch:    "working-" + MainBranch,
		Version:          opt.Version,
	}

	conf.Viper = setupViper(conf)
	return conf
}

func setupViper(conf Config) viper.Viper {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	dir, err := GetDirPath(conf)
	if err != nil {
		return viper.Viper{}
	}

	v.AddConfigPath(dir)
	err = v.ReadInConfig()
	if err != nil {
		return viper.Viper{}
	}

	return *v
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

//GetHistoryFilePath get the store file path from config
func GetHistoryFilePath(conf Config) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("failed to read home directory: %s", err)
	}

	storepath := fmt.Sprintf("%s/%s/%s", home, conf.CKPDir, HistoryFileName)
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

//GetTempHistoryStoreFilePath get the temporary store file path from config
func GetTempHistoryStoreFilePath(conf Config) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("failed to read home directory: %s", err)
	}

	storepath := fmt.Sprintf("%s/%s/%s", home, conf.CKPDir, StoreHistoryTempFileName)
	return storepath, nil
}

//GetStoreDirPath returns the path of the ckp store git repository
func GetStoreDirPath(conf Config) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("failed to read home directory: %s", err)
	}

	dir := fmt.Sprintf("%s/%s/%s", home, conf.CKPDir, conf.CKPStorageFolder)
	return dir, nil
}

//GetDirPath returns the path of the .ckp folder
func GetDirPath(conf Config) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("failed to read home directory: %s", err)
	}

	dir := fmt.Sprintf("%s/%s", home, conf.CKPDir)
	return dir, nil
}
