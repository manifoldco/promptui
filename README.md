[![Go Doc][godoc-image]][godoc-url]
[![Build Status][workflow-image]][workflow-url]

# promptui

This is a fork of [this](https://github.com/manifoldco/promptui) repo
for fixing issues (race conditions) and some housekeeping.

For documentation, please refer to the original repo [here](https://github.com/manifoldco/promptui).

## Overview

![promptui](https://media.giphy.com/media/xUNda0Ngb5qsogLsBi/giphy.gif)

## Quick Start

You can find more examples [here](./example).

### Prompt

```go
package main

import (
  "errors"
  "fmt"
  "strconv"

  "github.com/moorara/promptui"
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

### Select

```go
package main

import (
  "fmt"

  "github.com/moorara/promptui"
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


[godoc-url]: https://pkg.go.dev/github.com/moorara/promptui
[godoc-image]: https://pkg.go.dev/badge/github.com/moorara/promptui
[workflow-url]: https://github.com/moorara/promptui/actions
[workflow-image]: https://github.com/moorara/promptui/workflows/Go/badge.svg
