package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

type pepper struct {
	Name     string
	HeatUnit int
}

func main() {
	peppers := []pepper{
		pepper{Name: "Bell Pepper", HeatUnit: 0},
		pepper{Name: "Banana Pepper", HeatUnit: 100},
		pepper{Name: "Poblano", HeatUnit: 1000},
		pepper{Name: "Jalapeño", HeatUnit: 3500},
	}

	prompt := promptui.Select{
		Label:         "Spicy Level",
		Items:         peppers,
		ItemsTemplate: `{{.Name}} ({{color "red" .HeatUnit}})`,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %s\n", result)
}
