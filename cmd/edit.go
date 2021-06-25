package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"os"
	"time"

	"github.com/elhmn/ckp/internal/config"
	"github.com/elhmn/ckp/internal/store"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	EntryTemplateType         = "EntryTemplateType"
	CodeEntryTemplateType     = "CodeEntryTemplateType"
	SolutionEntryTemplateType = "SolutionEntryTemplateType"
)

const editorFileTemplate = `## You are editing the entry
## id:%s
##
##----------------------------------------------------------------------
## Add your comment
##----------------------------------------------------------------------

%s

##----------------------------------------------------------------------
## Set an alias
##----------------------------------------------------------------------

%s

##----------------------------------------------------------------------
## Here goes your code entry
##----------------------------------------------------------------------

%s

##----------------------------------------------------------------------
## Here goes your solution entry
##----------------------------------------------------------------------

%s

##----------------------------------------------------------------------
## Note that you can't set both
## a code and solution entry
## the code entry will take precedence
`

//NewEditCommand create new cobra command for the edit command
func NewEditCommand(conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:   "edit [code_id | solution_id]",
		Short: "edit code or solution entries from the store",
		Long: `edit code or solution entries from the store

		example: ckp edit
		Will prompt an interactive UI that will allow you to search and delete
		a code or solution entry

		example: ckp edit <entry_id>
		Will edit the entry corresponding the entry_id
`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := editCommand(cmd, args, conf); err != nil {
				fmt.Fprintf(conf.OutWriter, "Error: %s\n", err)
				return
			}
		},
	}

	command.PersistentFlags().Bool("from-history", false, `list code and solution records from history`)
	command.PersistentFlags().StringP("comment", "m", "", `ckp edit -m <comment>`)
	command.PersistentFlags().StringP("alias", "a", "", `ckp edit -a <alias>`)
	command.PersistentFlags().StringP("code", "c", "", `ckp edit -c <code>`)
	command.PersistentFlags().StringP("solution", "s", "", `ckp edit -s <solution>`)
	command.PersistentFlags().BoolP("interactive", "i", false, `open a text editor`)
	return command
}

func editCommand(cmd *cobra.Command, args []string, conf config.Config) error {
	var entryID string
	if len(args) >= 1 {
		entryID = args[0]
	}

	if err := cmd.Flags().Parse(args); err != nil {
		return err
	}
	flags := cmd.Flags()
	fromHistory, err := flags.GetBool("from-history")
	if err != nil {
		return fmt.Errorf("could not parse `fromHistory` flag: %s", err)
	}

	isInteractive, err := flags.GetBool("interactive")
	if err != nil {
		return fmt.Errorf("could not parse `interactive` flag: %s", err)
	}

	//Setup spinner
	conf.Spin.Start()
	defer conf.Spin.Stop()

	dir, err := config.GetStoreDirPath(conf)
	if err != nil {
		return fmt.Errorf("failed get repository path: %s", err)
	}

	//Get the store file path
	var storeFilePath string
	if !fromHistory {
		storeFilePath, err = config.GetStoreFilePath(conf)
		if err != nil {
			return fmt.Errorf("failed to get the store file path: %s", err)
		}
	} else {
		storeFilePath, err = config.GetHistoryFilePath(conf)
		if err != nil {
			return fmt.Errorf("failed to get the history store file path: %s", err)
		}
	}

	conf.Spin.Message(" pulling remote changes...")
	err = pullRemoteChanges(conf, dir, storeFilePath)
	if err != nil {
		return fmt.Errorf("failed to pull remote changes: %s", err)
	}
	conf.Spin.Message(" remote changes pulled")

	conf.Spin.Message(" removing changes")
	storeFile, storeData, storeBytes, err := loadStore(storeFilePath)
	if err != nil {
		return fmt.Errorf("failed to load the store: %s", err)
	}

	index, err := getScriptEntryIndex(conf, storeData.Scripts, entryID)
	if err != nil {
		return fmt.Errorf("failed to get script `%s` entry index: %s", entryID, err)
	}

	script, err := createNewEntry(flags, storeData.Scripts[index])
	if err != nil {
		return fmt.Errorf("failed to create new script entry: %s", err)
	}

	//if it is an interactive update
	if len(args) == 0 || isInteractive {
		conf.Spin.Stop()
		s, err := getNewEntryDataFromFile(conf, script, EntryTemplateType)
		if err != nil {
			return fmt.Errorf("failed to get new entry from the editor %s", err)
		}

		//Generate an ID for the newly added script
		s.ID, err = store.GenereateIdempotentID(s.Code.Content, s.Comment, s.Code.Alias, s.Solution.Content)
		if err != nil {
			return fmt.Errorf("failed to generate idem potent ID: %s", err)
		}

		script.ID = s.ID
		script.Code = s.Code
		script.Solution = s.Solution
		script.Comment = s.Comment
		conf.Spin.Start()
	}

	//Check if entry already exist
	if storeData.EntryAlreadyExist(script.ID) {
		return fmt.Errorf("An identical record was found in the storage, please try `ckp edit %s`", script.ID)
	}

	//Remove script entry
	storeData.Scripts, err = editScriptEntry(storeData.Scripts, script, index)
	if err != nil {
		return fmt.Errorf("failed to edit script entry: %s", err)
	}

	tempFile, err := createTempFile(conf, storeBytes)
	if err != nil {
		return fmt.Errorf("failed to create tempFile: %s", err)
	}

	//Save storeData in store
	if err := saveStore(storeData, storeBytes, storeFile, tempFile); err != nil {
		return fmt.Errorf("failed to save store in %s:  %s", storeFile, err)
	}

	//Delete the temporary file
	if err := os.RemoveAll(tempFile); err != nil {
		return fmt.Errorf("failed to delete file %s: %s", tempFile, err)
	}

	conf.Spin.Message(" pushing local changes...")
	err = pushLocalChanges(conf, dir, commitRemoveAction, storeFilePath)
	if err != nil {
		return fmt.Errorf("failed to push local changes: %s", err)
	}
	conf.Spin.Message(" local changes pushed")

	fmt.Fprintf(conf.OutWriter, "\nYour entry was successfully edited!\n")
	fmt.Fprintf(conf.OutWriter, "\n%s", storeData.Scripts[len(storeData.Scripts)-1])
	return nil
}

func editScriptEntry(scripts []store.Script, script store.Script, index int) ([]store.Script, error) {
	return append(scripts[:index], append(scripts[index+1:], script)...), nil
}

//createNewEntry return a new code Script entry
func createNewEntry(flags *pflag.FlagSet, script store.Script) (store.Script, error) {
	timeNow := time.Now()

	//Get alias
	alias := script.Code.Alias
	aliasTmp, err := flags.GetString("alias")
	if err != nil {
		return store.Script{}, fmt.Errorf("could not parse `alias` flag: %s", err)
	}
	if aliasTmp != "" {
		alias = aliasTmp
	}

	//Get comment
	comment := script.Comment
	commentTmp, err := flags.GetString("comment")
	if err != nil {
		return store.Script{}, fmt.Errorf("could not parse `comment` flag: %s", err)
	}
	if commentTmp != "" {
		comment = commentTmp
	}

	//Get code
	var code string
	if script.Code.Content != "" {
		code = script.Code.Content
	}
	codeTmp, err := flags.GetString("code")
	if err != nil {
		return store.Script{}, fmt.Errorf("could not parse `code` flag: %s", err)
	}
	if codeTmp != "" {
		code = codeTmp
	}

	//Get Solution
	var solution string
	if script.Solution.Content != "" {
		solution = script.Solution.Content
	}
	solutionTmp, err := flags.GetString("solution")
	if err != nil {
		return store.Script{}, fmt.Errorf("could not parse `solution` flag: %s", err)
	}
	if solutionTmp != "" {
		solution = solutionTmp
	}

	//Generate script entry unique id
	id, err := store.GenereateIdempotentID(code, comment, alias, solution)
	if err != nil {
		return store.Script{}, fmt.Errorf("failed to generate idem potent id: %s", err)
	}

	return store.Script{
		ID:           id,
		Comment:      comment,
		CreationTime: timeNow,
		UpdateTime:   timeNow,
		Code: store.Code{
			Content: code,
			Alias:   alias,
		},
		Solution: store.Solution{
			Content: solution,
		},
	}, nil
}

func getNewEntryDataFromFile(conf config.Config, origEntry store.Script, templateType string) (store.Script, error) {
	s := origEntry
	content := ""

	//Get template content
	switch templateType {
	case CodeEntryTemplateType:
		content = fmt.Sprintf(editorFileCodeTemplate, origEntry.Comment, origEntry.Code.Alias, origEntry.Code.Content)
	case SolutionEntryTemplateType:
		content = fmt.Sprintf(editorFileSolutionTemplate, origEntry.Comment, origEntry.Solution.Content)
	case EntryTemplateType:
		content = fmt.Sprintf(editorFileTemplate, origEntry.ID, origEntry.Comment, origEntry.Code.Alias, origEntry.Code.Content, origEntry.Solution.Content)
	}

	dir, err := config.GetDirPath(conf)
	if err != nil {
		return s, err
	}
	destination := fmt.Sprintf("%s/entry.%s.sh", dir, origEntry.ID)

	//Create the file with the original script data
	if err = ioutil.WriteFile(destination, []byte(content), 0666); err != nil {
		return s, fmt.Errorf("failed to write to file %s: %s", destination, err)
	}

	//Open and edit that file
	err = conf.Exec.OpenEditor("", destination)
	if err != nil {
		return s, err
	}

	s, err = parseDataFromEditorTemplateFile(destination, templateType)
	if err != nil {
		return s, fmt.Errorf("failed to parse data from template file file %s: %s", destination, err)
	}

	//Delete the temporary file
	if err := os.RemoveAll(destination); err != nil {
		return s, fmt.Errorf("failed to delete file %s: %s", destination, err)
	}

	return s, nil
}

func parseDataFromEditorTemplateFile(filepath string, templateType string) (store.Script, error) {
	//get store from template file
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return store.Script{}, err
	}

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return store.Script{}, fmt.Errorf("failed to read file: %s", err)
	}

	switch templateType {
	case CodeEntryTemplateType:
		return parseCodeDataFromEditorTemplateString(string(data)), nil
	case SolutionEntryTemplateType:
		return parseSolutionDataFromEditorTemplateString(string(data)), nil
	}

	return parseDataFromEditorTemplateString(string(data)), nil
}

func parseDataFromEditorTemplateString(data string) store.Script {
	lines := strings.Split(data, "\n")

	//get comment
	i := moveToNextEntry(lines, 0)
	comment, i := getEntry(lines, i)

	//get alias
	i = moveToNextEntry(lines, i)
	alias, i := getEntry(lines, i)

	//get code
	i = moveToNextEntry(lines, i)
	code, i := getEntry(lines, i)

	//get solution
	i = moveToNextEntry(lines, i)
	solution, _ := getEntry(lines, i)

	if code != "" {
		return store.Script{
			Comment: comment,
			Code: store.Code{
				Content: code,
				Alias:   alias,
			},
			Solution: store.Solution{},
		}
	}

	return store.Script{
		Comment: comment,
		Solution: store.Solution{
			Content: solution,
		},
		Code: store.Code{},
	}
}

func parseCodeDataFromEditorTemplateString(data string) store.Script {
	lines := strings.Split(data, "\n")

	//get comment
	comment, _ := getComment(lines)

	//get alias
	alias, _ := getAlias(lines)

	//get code
	code := getCode(lines)

	return store.Script{
		Comment: comment,
		Code: store.Code{
			Content: code,
			Alias:   alias,
		},
		Solution: store.Solution{},
	}
}

func parseSolutionDataFromEditorTemplateString(data string) store.Script {
	lines := strings.Split(data, "\n")

	//get comment
	comment, _ := getComment(lines)

	//get solution
	solution := getSolution(lines)

	return store.Script{
		Comment: comment,
		Solution: store.Solution{
			Content: solution,
		},
		Code: store.Code{},
	}
}

//moveToNextEntry skips comments and return the index to the next valid line
func moveToNextEntry(lines []string, i int) int {
	if i >= len(lines) {
		return i - 1
	}

	for i := i; i < len(lines); i++ {
		line := lines[i]
		//if "##" is not at the beginning of the line
		if strings.Index(line, "##") != 0 {
			return i
		}
	}

	return i
}

//getComment iterate over the file content and returns the first `## comment:` line
// it returns the comment it found or an error if no comment was founds
func getComment(lines []string) (string, error) {
	for _, line := range lines {
		if strings.Contains(line, "## comment") {
			s := strings.Split(line, ":")
			if len(s) >= 2 {
				return strings.Join(s[1:], ""), nil
			}
		}
	}

	return "", fmt.Errorf("comment not found")
}

//getAlias iterate over the file content and returns the first `## alias:` line
// it returns the alias it found or an error if no alias was founds
func getAlias(lines []string) (string, error) {
	for _, line := range lines {
		if strings.Contains(line, "## alias") {
			s := strings.Split(line, ":")
			if len(s) >= 2 {
				return strings.Join(s[1:], ""), nil
			}
		}
	}

	return "", fmt.Errorf("alias not found")
}

//getCode iterate over the file content and returns
//everything that does not start with "##" as a code
func getCode(lines []string) string {
	code := ""
	for _, line := range lines {
		//if "##" is at the beginning of the line
		if strings.Index(line, "##") != 0 {
			code += line + "\n"
		}
	}

	return strings.Trim(code, "\n")
}

//getSolution iterate over the file content and returns
//everything that does not start with "##" as a solution
func getSolution(lines []string) string {
	solution := ""
	for _, line := range lines {
		//if "##" is at the beginning of the line
		if strings.Index(line, "##") != 0 {
			solution += line + "\n"
		}
	}

	return strings.Trim(solution, "\n")
}

//getEntry returns the entry and returns an index to the next line
func getEntry(lines []string, i int) (string, int) {
	entry := ""
	if i >= len(lines) {
		return entry, i - 1
	}

	for i := i; i < len(lines); i++ {
		line := lines[i]
		//if "##" is at the beginning of the line
		if strings.Index(line, "##") == 0 {
			return strings.Trim(entry, "\n"), i
		}

		entry += line + "\n"
	}

	return strings.Trim(entry, "\n"), i
}
