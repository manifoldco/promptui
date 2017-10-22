package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func main() {
	prompt := promptui.SelectWithAdd{
		Label:    "What's your text editor",
		Items:    []string{"Vim", "Emacs", "Sublime", "VSCode", "Atom"},
		AddLabel: "Other",
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %s\n", result)
}
