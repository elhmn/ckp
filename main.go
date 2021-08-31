package main

import (
	"os"

	"github.com/elhmn/ckp/cmd"
	"github.com/elhmn/ckp/internal/config"
)

var version = "0.0.0.dev"

func main() {
	conf := config.NewDefaultConfig(config.Options{Version: version})
	command := cmd.NewCKPCommand(conf)

	err := command.Execute()
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
