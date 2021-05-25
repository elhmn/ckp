package printers

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

//Confirm prompt a confirmation message
//
//Return true if the user entered Y/y and false if entered n/N
func Confirm(message string) bool {
	validate := func(input string) error {
		input = strings.ToLower(strings.TrimSpace(input))
		if input != "y" && input != "n" {
			return fmt.Errorf("wrong input %s, was expecting `y` or `n`", input)
		}

		return nil
	}

	msg := message + " Press (y/n)"
	prompt := promptui.Prompt{
		Label:    msg,
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return false
	}
	input := strings.ToLower(strings.TrimSpace(result))

	return input == "y"
}
