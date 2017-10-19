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
		{Name: "Aleppo", HeatUnit: 10000, Peppers: 4},
		{Name: "Tabasco", HeatUnit: 30000, Peppers: 5},
		{Name: "Malagueta", HeatUnit: 50000, Peppers: 6},
		{Name: "Habanero", HeatUnit: 100000, Peppers: 7},
		{Name: "Red Savina Habanero", HeatUnit: 350000, Peppers: 8},
		{Name: "Dragon’s Breath", HeatUnit: 855000, Peppers: 9},
	}

	fns := promptui.FuncMap
	fns["rangeLoop"] = count

	tpl := "{{ .Name | bold }} ({{ .HeatUnit | red | italic }}) {{ range rangeLoop .Peppers }}\U0001F525{{ end }}"

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

func count(n int) []struct{} {
	return make([]struct{}, n)
}
