package helm

import (
	"fmt"

	"github.com/c-bata/go-prompt"
)

const (
	PromptYes    = "Y"
	PromptNo     = "N"
	PromptCancel = "C"
)

func PromptYesNoCancel(label string) string {
	completer := func(d prompt.Document) []prompt.Suggest {
		s := []prompt.Suggest{
			{Text: PromptYes, Description: "Yes"},
			{Text: PromptNo, Description: "No"},
			{Text: PromptCancel, Description: "Cancel"},
		}
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	}
	fmt.Println(label, " (Yes/No)")
	t := prompt.Input("> ", completer)
	return t
}
