package main

import (
	"os"

	"github.com/elhmn/ckp/cmd"
	"github.com/elhmn/ckp/internal/config"
)

func main() {
	conf := config.NewDefaultConfig()
	command := cmd.NewCKPCommand(conf)

	err := command.Execute()
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
