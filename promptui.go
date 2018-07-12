// Package Promptui is a library providing a simple interface to create command-line prompts for go.
// It can be easily integrated into spf13/cobra, urfave/cli or any cli go application.
//
// Promptui has two main input modes:
//
// 		Prompt provides a single line for user input. Prompt supports optional live validation, confirmation and masking the input.
//
//		Select provides a list of options to choose from. Select supports pagination, search, detailed view and custom templates.
//
// Basic Usage
//
// Prompt
// 		package main
//
//		import (
//			"errors"
//			"fmt"
//			"strconv"
//
//			"github.com/manifoldco/promptui"
//		)
//
//		func main() {
//			validate := func(input string) error {
//				_, err := strconv.ParseFloat(input, 64)
//				if err != nil {
//					return errors.New("Invalid number")
//				}
//				return nil
//			}
//
//			prompt := promptui.Prompt{
//				Label:    "Number",
//				Validate: validate,
//			}
//
//			result, err := prompt.Run()
//
//			if err != nil {
//				fmt.Printf("Prompt failed %v\n", err)
//				return
//			}
//
//			fmt.Printf("You choose %q\n", result)
//		}
//
// Select
//		package main
//
//		import (
//			"fmt"
//
//			"github.com/manifoldco/promptui"
//		)
//
//		func main() {
//			prompt := promptui.Select{
//			Label: "Select Day",
//			Items: []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday",
//				"Saturday", "Sunday"},
//		}
//
//		_, result, err := prompt.Run()
//
//		if err != nil {
//			fmt.Printf("Prompt failed %v\n", err)
//			return
//		}
//
//		fmt.Printf("You choose %q\n", result)
//	}
package promptui

import "errors"

// ErrEOF is the error returned from prompts when EOF is encountered. It will be triggered when using the select mode
// and starting a search if no element is found.
var ErrEOF = errors.New("^D")

// ErrInterrupt is the error returned from prompts when an interrupt (ctrl-c) is
// encountered.
var ErrInterrupt = errors.New("^C")

// ErrAbort is the error returned when confirm prompts are supplied "n"
var ErrAbort = errors.New("")

// ValidateFunc is a placeholder type for any validation functions that validates a given input. It should return
// a ValidationError if the input is not valid.
type ValidateFunc func(string) error
