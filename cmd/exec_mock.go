package cmd

type MockedExec struct {
	RunErrorOutput error
	RunOutput      []byte

	DoGitCloneErrorOutput error
	DoGitCloneOutput      string

	CreateFolderIfDoesNotExistErrorOutput error
}

func (ex MockedExec) Run(dir string, command string, args ...string) ([]byte, error) {
	return ex.RunOutput, ex.RunErrorOutput
}

func (ex MockedExec) DoGitClone(dir string, args ...string) (string, error) {
	return ex.DoGitCloneOutput, ex.DoGitCloneErrorOutput
}

func (ex MockedExec) CreateFolderIfDoesNotExist(dir string) error {
	return ex.CreateFolderIfDoesNotExistErrorOutput
}
