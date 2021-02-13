package config

import (
	"io"

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
