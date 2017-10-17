package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func main() {
	values := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	prompt := promptui.Select{
		Label: "Select Day",
		Items: values,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)
}
