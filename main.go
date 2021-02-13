package main

import (
	"os"

	"github.com/elhmn/ckp/cmd"
	"github.com/elhmn/ckp/internal/config"
	"github.com/elhmn/ckp/internal/exec"
)

func newConfig() config.Config {
	return config.Config{
		Exec:             exec.NewExec(),
		OutWriter:        os.Stdout,
		ErrWriter:        os.Stderr,
		CKPDir:           ".ckp",
		CKPStorageFolder: "repo",
	}
}

func main() {
	conf := newConfig()
	command := cmd.NewCKPCommand(conf)

	err := command.Execute()
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
