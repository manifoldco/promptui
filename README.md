# promptui

Interactive prompt for command-line applications.

[Code of Conduct](./CODE_OF_CONDUCT.md) |
[Contribution Guidelines](./.github/CONTRIBUTING.md)

[![GitHub release](https://img.shields.io/github/tag/manifoldco/promptui.svg?label=latest)](https://github.com/manifoldco/promptui/releases)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/manifoldco/promptui)
[![Travis](https://img.shields.io/travis/manifoldco/promptui/master.svg)](https://travis-ci.org/manifoldco/promptui)
[![Go Report Card](https://goreportcard.com/badge/github.com/manifoldco/promptui)](https://goreportcard.com/report/github.com/manifoldco/promptui)
[![License](https://img.shields.io/badge/license-BSD-blue.svg)](./LICENSE.md)

## Overview

![custom_select](/_examples/imgs/custom_select.gif)

Promptui is a library providing a simple interface to create command-line
prompts for go. It can be easily integrated into
[spf13/cobra](https://github.com/spf13/cobra),
[urfave/cli](https://github.com/urfave/cli) or any cli go application.

Promptui has two main input modes:

- `.Prompt` provides a single line for user input. Prompt supports
  optional live validation, confirmation and masking the input.

- `.Select` provides a list of options to choose from. Select supports
  pagination, search, detailed view and custom templates.

For a full list of options check [GoDoc](https://godoc.org/github.com/manifoldco/promptui).

## Basic Usage

### Prompt

[![asciicast](https://asciinema.org/a/148944.png)](https://asciinema.org/a/148944)

```go
package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
)

func main() {
	validate := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("Invalid number")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Number",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)
}
```

### Selection

[![asciicast](https://asciinema.org/a/148946.png)](https://asciinema.org/a/148946)

```go
package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func main() {
	prompt := promptui.Select{
		Label: "Select Day",
		Items: []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday",
			"Saturday", "Sunday"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)
}
```

### More Examples

See full list of examples at [examples](https://github.com/manifoldco/promptui/tree/master/_examples)
