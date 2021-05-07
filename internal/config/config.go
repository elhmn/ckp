package config

import (
	"io"
	"os"

	"github.com/elhmn/ckp/internal/exec"
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
