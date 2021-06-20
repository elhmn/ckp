package exec

import (
	"os"
	exe "os/exec"
)

const defaultEditor = "vim"

//OpenEditor opens the editor
func (ex Exec) OpenEditor(editor string, args ...string) error {
	if editor == "" {
		editor = defaultEditor
	}

	cmd := exe.Command(editor, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
