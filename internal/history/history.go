package history

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

const (
	BashHistoryFile = ".bash_history"
	ZshHistoryFile  = ".zsh_history"
)

//GetHistoryRecords returns a list of code records
//found in your history files
func GetHistoryRecords() ([]string, error) {
	records := []string{}

	// 	Get bash history records
	bashRecords, err := getBashHistoryRecords()
	if err != nil {
		return nil, err
	}
	records = append(records, bashRecords...)

	//Get zsh history records
	zshRecords, err := getZshHistoryRecords()
	if err != nil {
		return nil, err
	}

	records = append(records, zshRecords...)
	return records, nil
}

func getBashHistoryRecords() ([]string, error) {
	data, err := readFileData(BashHistoryFile)
	if err != nil {
		return nil, err
	}

	records := strings.Split(data, "\n")
	return records, nil
}

func getZshHistoryRecords() ([]string, error) {
	data, err := readFileData(ZshHistoryFile)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(data, "\n")
	records := []string{}
	for _, l := range lines {
		recordStartIndex := strings.Index(l, ";")
		if recordStartIndex < 0 {
			continue
		}

		line := l[recordStartIndex+1:]
		records = append(records, line)
	}
	return records, nil
}

func readFileData(filepath string) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("failed to read home directory: %s", err)
	}
	filepath = fmt.Sprintf("%s/%s", home, filepath)

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return "", nil
	}

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %s", filepath, err)
	}

	return string(data), nil
}
