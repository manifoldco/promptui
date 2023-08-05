package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
)

func main() {
	items := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday",
	"Saturday", "Sunday"}
	// Default count from 1!!
	prompt := promptui.Select{
		Label: "Select Day",
		Items: items,
		Default: 3,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)
}
