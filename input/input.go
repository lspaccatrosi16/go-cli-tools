package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/lspaccatrosi16/go-cli-tools/pkgError"
	"github.com/manifoldco/promptui"
)

const LINE_UP = "\033[1A"
const LINE_CLEAR = "\x1b[2K"

var wrap = pkgError.WrapErrorFactory("input")

func getInputTemplate() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Active:   "{{ .Name | green }}",
		Inactive: "{{ .Name }}",
		Selected: "{{ .Name }}",
	}
}

func makeSelector(label string, items []SelectOption) promptui.Select {
	template := getInputTemplate()

	prompt := promptui.Select{
		Label:        label,
		Templates:    template,
		Items:        items,
		HideSelected: true,
	}

	return prompt
}

func makeSearchableSelector(label string, items []SelectOption) promptui.Select {
	template := getInputTemplate()

	searcher := func(input string, index int) bool {
		item := items[index]
		inputted := strings.ToLower(input)
		name := strings.ToLower(item.Name)
		return strings.Contains(name, inputted)
	}

	prompt := promptui.Select{
		Label:        label,
		Items:        items,
		Templates:    template,
		Searcher:     searcher,
		IsVimMode:    false,
		HideSelected: true,
	}

	return prompt

}

func GetSelection(label string, items []SelectOption) (string, error) {
	v, _, err := GetSelectionIdx(label, items)

	if err != nil {
		return "", wrap(err)
	}

	return v, nil
}

func GetSelectionIdx(label string, items []SelectOption) (string, int, error) {
	prompt := makeSelector(label, items)

	i, _, err := prompt.Run()

	if err != nil {
		return "", -1, wrap(err)
	}

	return items[i].Value, i, nil
}

func GetSearchableSelection(label string, items []SelectOption) (string, error) {
	v, _, err := GetSearchableSelectionIdx(label, items)

	if err != nil {
		return "", wrap(err)
	}

	return v, nil
}

func GetSearchableSelectionIdx(label string, items []SelectOption) (string, int, error) {
	prompt := makeSearchableSelector(label, items)

	i, _, err := prompt.Run()

	if err != nil {
		return "", -1, wrap(err)
	}

	return items[i].Value, i, nil
}

func GetConfirmSelection(label string) (bool, error) {
	items := []SelectOption{
		{Name: "Yes", Value: "y"},
		{Name: "No", Value: "n"},
	}

	val, err := GetSelection(label, items)

	if err != nil {
		return false, wrap(err)
	}

	return val == "y", nil
}

func GetInput(label string) string {
	stubValidator := func(in string) error {
		return nil
	}

	return GetValidatedInput(label, stubValidator)
}

func GetValidatedInput(label string, validator func(str string) error) string {
	var result string

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s? ", label)
	linesUsed := 0

	for {
		scanner.Scan()
		linesUsed++

		result = scanner.Text()

		validationError := validator(result)

		if validationError != nil {
			fmt.Printf("ERROR: %s\n", validationError.Error())
			linesUsed++
			continue
		}
		break
	}

	for i := 0; i < linesUsed; i++ {
		fmt.Print(LINE_UP)
		fmt.Print(LINE_CLEAR)
	}

	fmt.Printf("%s: %s\n", label, result)

	return result
}
