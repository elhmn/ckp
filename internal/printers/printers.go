package printers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/elhmn/ckp/internal/store"
	"github.com/manifoldco/promptui"
)

const selectItemsSize = 10

var defaultPrinters = Printers{}

type IPrinters interface {
	Confirm(message string) bool
	SelectScriptEntry(scripts []store.Script, entryType string) (int, string, error)
}

type Printers struct{}

//NewPrinters returns new printers struct
func NewPrinters() *Printers {
	return &Printers{}
}

func (p Printers) Confirm(message string) bool {
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

//Confirm prompt a confirmation message
//
//Return true if the user entered Y/y and false if entered n/N
func Confirm(message string) bool {
	return defaultPrinters.Confirm(message)
}

func SelectScriptEntry(scripts []store.Script, entryType string) (int, string, error) {
	return defaultPrinters.SelectScriptEntry(scripts, entryType)
}

//SelectScriptEntry prompt a search
//returns the selected entry index
func (p Printers) SelectScriptEntry(scripts []store.Script, entryType string) (int, string, error) {
	searchScript := func(input string, index int) bool {
		s := scripts[index]
		if entryType == store.EntryTypeCode {
			return s.Code.Content != "" && DoesScriptContain(s, input)
		} else if entryType == store.EntryTypeSolution {
			return s.Solution.Content != "" && DoesScriptContain(s, input)
		}

		return DoesScriptContain(s, input)
	}

	prompt := promptui.Select{
		Label:             "Enter your search text",
		Items:             scripts,
		Size:              selectItemsSize,
		StartInSearchMode: true,
		Searcher:          searchScript,
		Templates:         getTemplates(),
	}

	i, result, err := prompt.Run()
	if err != nil {
		return i, "", fmt.Errorf("prompt failed %v", err)
	}

	return i, result, nil
}

func trimText(s string) string {
	if len(s) > 50 {
		return s[:50] + "..."
	}
	return s
}

func getTemplates() *promptui.SelectTemplates {
	funcMap := promptui.FuncMap
	funcMap["inline"] = func(s string) string {
		return strings.ReplaceAll(trimText(s), "\n", " ")
	}

	//if you find a hard time understand it check out golang templating format documentation
	//here https://golang.org/pkg/text/template
	return &promptui.SelectTemplates{
		Label: "{{ if .Code.Content -}} {{`code:` | bold | green}} " +
			"{{ inline .Code.Content}} {{- else -}} {{ inline .Solution.Content }} {{ end }}",
		Active: "* {{ if .Code.Content -}} {{`code:` | bold | green}} {{ inline .Code.Content | bold}} {{ else }} " +
			"{{`solution:` | bold | yellow }} {{ inline .Solution.Content | bold }} {{ end }}",
		Inactive: "{{ if .Code.Content -}} {{`code:` | green }} {{ inline .Code.Content }} " +
			"{{- else -}} {{`solution:` | yellow}} {{ inline .Solution.Content }} {{ end }}",
		Selected: " {{ `✓` | green }} {{if .Code.Content -}} {{ inline .Code.Content | bold }} {{- else -}} {{ inline .Solution.Content | bold }} {{ end }}",
		Details: "Type: {{- if .Code.Content }} code {{ else }} solution {{- end }}" +
			"{{ if .Code.Alias }} | Alias: {{ .Code.Alias }} {{- end }}" +
			"{{ if .Comment }} | Comment: {{ .Comment }} {{- end }}",
		FuncMap: funcMap,
	}
}

func extractScriptStringContent(script store.Script) string {
	code := strings.Replace(strings.ToLower(script.Code.Content), " ", "", -1)
	solution := strings.Replace(strings.ToLower(script.Solution.Content), " ", "", -1)
	comment := strings.Replace(strings.ToLower(script.Comment), " ", "", -1)
	alias := strings.Replace(strings.ToLower(script.Code.Alias), " ", "", -1)
	content := fmt.Sprintf("%s %s %s %s", code, solution, comment, alias)
	return content
}

//DoesScriptContain return true if the script contains the input value
func DoesScriptContain(script store.Script, input string) bool {
	input = strings.TrimSpace(strings.ToLower(input))

	//Build pattern
	pattern := ".*" + strings.Join(strings.Split(input, " "), ".*")

	matched, err := regexp.Match(pattern, []byte(extractScriptStringContent(script)))
	if err != nil {
		return false
	}

	return matched
}
