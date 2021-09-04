package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/blang/semver"
	"github.com/elhmn/ckp/internal/config"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

//NewUpdateCommand will update your binary to the latest release
func NewUpdateCommand(conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "update",
		Short: "Update the binary to the latest release",
		Run: func(cmd *cobra.Command, args []string) {
			if err := updateCommand(conf); err != nil {
				fmt.Fprintf(conf.OutWriter, "Error: %s\n", err)
				return
			}
		},
	}

	return command
}

func updateCommand(conf config.Config) error {
	return confirmAndUpdate(conf)
}

func confirmAndUpdate(conf config.Config) error {
	latest, found, err := selfupdate.DetectLatest(conf.Repository)
	if err != nil {
		return fmt.Errorf("Error occurred while detecting version: %s", err)
	}

	v := semver.MustParse(conf.Version)
	if !found || latest.Version.LTE(v) {
		fmt.Fprintf(conf.OutWriter, "Current version is the latest\n")
		return nil
	}

	fmt.Fprintf(conf.OutWriter, "Do you want to update to %s ? (y/n): \n", latest.Version)
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil || (input != "y\n" && input != "n\n") {
		fmt.Fprintf(conf.OutWriter, "Invalid input\n")
		return nil
	}
	if input == "n\n" {
		return nil
	}

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("Could not locate executable path: %s", err)
	}
	if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
		return fmt.Errorf("Error occurred while updating binary: %s", err)
	}

	fmt.Fprintf(conf.OutWriter, "Successfully updated to version %s\n", latest.Version)
	return nil
}
