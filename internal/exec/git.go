package exec

import "os/exec"

//Run run command and return output
func (ex Exec) Run(dir string, command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

//DoGitClone execute a `git clone <args...>`
func (ex Exec) DoGitClone(dir string, args ...string) (string, error) {
	cmd := "clone"
	args = append([]string{cmd}, args...)
	output, err := ex.Run(dir, "git", args...)
	return string(output), err
}
