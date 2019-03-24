package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func main() {
	chosed := []int{}
	options := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday",
		"Saturday", "Sunday"}
	prompt := promptui.Select{
		Label:       "Select Day",
		Checkbox:    true,
		ChosedIcon:  promptui.IconGood,
		ChosenIndex: &chosed,
		Items:       options,
	}

	_, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %v\n", chosed)
}
