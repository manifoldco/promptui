# promptui

Interactive prompt for command-line applications.

[Code of Conduct](./CODE_OF_CONDUCT.md) |
[Contribution Guidelines](./.github/CONTRIBUTING.md)

[![GitHub release](https://img.shields.io/github/tag/manifoldco/promptui.svg?label=latest)](https://github.com/manifoldco/promptui/releases)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/manifoldco/promptui)
[![Travis](https://img.shields.io/travis/manifoldco/promptui/master.svg)](https://travis-ci.org/manifoldco/promptui)
[![Go Report Card](https://goreportcard.com/badge/github.com/manifoldco/promptui)](https://goreportcard.com/report/github.com/manifoldco/promptui)
[![License](https://img.shields.io/badge/license-BSD-blue.svg)](./LICENSE.md)

## Usage

### Selection

```go
package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func main() {
	values := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}

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
```
