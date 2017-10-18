package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

type pepper struct {
	Name     string
	HeatUnit int
	Peppers  int
}

func main() {
	peppers := []pepper{
		{Name: "Bell Pepper", HeatUnit: 0, Peppers: 0},
		{Name: "Banana Pepper", HeatUnit: 100, Peppers: 1},
		{Name: "Poblano", HeatUnit: 1000, Peppers: 2},
		{Name: "Jalapeño", HeatUnit: 3500, Peppers: 3},
	}

	fns := promptui.FuncMap
	fns["rangeLoop"] = rangeLoop

	tpl := "{{bold .Name}} ({{red .HeatUnit}}) {{ range rangeLoop .Peppers }}\U0001F525{{ end }}"

	prompt := promptui.Select{
		Label:         "Spicy Level",
		Items:         peppers,
		ItemsTemplate: tpl,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %s\n", result)
}

func rangeLoop(n int) []struct{} {
	return make([]struct{}, n)
}
