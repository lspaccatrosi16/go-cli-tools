package input

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/manifoldco/promptui"
)

const LINE_UP = "\033[1A"
const LINE_CLEAR = "\x1b[2K"

func getInputTemplate() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Active:   "{{ .Name | cyan }}",
		Inactive: "{{ .Name }}",
		Selected: "{{ .Name }}",
	}
}

func makeSelector(label string, items []SelectOption) promptui.Select {
	template := getInputTemplate()

	prompt := promptui.Select{
		Label:     label,
		Templates: template,
		Items:     items,
	}

	return prompt
}

func GetSelection(label string, items []SelectOption) string {
	v, _ := GetSelectionIdx(label, items)
	return v
}

func GetSelectionIdx(label string, items []SelectOption) (string, int) {
	prompt := makeSelector(label, items)

	i, _, err := prompt.Run()

	if err != nil {
		log.Fatalln(err)
	}

	return items[i].Value, i
}

func GetConfirmSelection(label string) bool {
	items := []SelectOption{
		{Name: "Yes", Value: "y"},
		{Name: "No", Value: "n"},
	}

	val := GetSelection(label, items)

	return val == "y"
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
