package exec

//DoGit execute a `git command <args...>`
func (ex Exec) DoGit(dir string, args ...string) (string, error) {
	output, err := ex.Run(dir, "git", args...)
	return string(output), err
}

//DoGitClone execute a `git clone <args...>`
func (ex Exec) DoGitClone(dir string, args ...string) (string, error) {
	cmd := "clone"
	args = append([]string{cmd}, args...)
	output, err := ex.DoGit(dir, args...)
	return string(output), err
}

//DoGitPush execute a `git push <args...>`
func (ex Exec) DoGitPush(dir string, args ...string) (string, error) {
	cmd := "push"
	args = append([]string{cmd}, args...)
	output, err := ex.DoGit(dir, args...)
	return string(output), err
}
