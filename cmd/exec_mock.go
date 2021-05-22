package cmd

type MockedExec struct {
	RunErrorOutput error
	RunOutput      []byte

	DoGitCloneErrorOutput error
	DoGitCloneOutput      string

	DoGitPushErrorOutput error
	DoGitPushOutput      string

	DoGitErrorOutput error
	DoGitOutput      string

	CreateFolderIfDoesNotExistErrorOutput error
}

func (ex MockedExec) Run(dir string, command string, args ...string) ([]byte, error) {
	return ex.RunOutput, ex.RunErrorOutput
}

func (ex MockedExec) DoGit(dir string, args ...string) (string, error) {
	return ex.DoGitOutput, ex.DoGitErrorOutput
}

func (ex MockedExec) DoGitClone(dir string, args ...string) (string, error) {
	return ex.DoGitCloneOutput, ex.DoGitCloneErrorOutput
}

func (ex MockedExec) DoGitPush(dir string, args ...string) (string, error) {
	return ex.DoGitPushOutput, ex.DoGitPushErrorOutput
}

func (ex MockedExec) CreateFolderIfDoesNotExist(dir string) error {
	return ex.CreateFolderIfDoesNotExistErrorOutput
}
