// Package promptui is a library providing a simple interface to create command-line prompts for go.
// It can be easily integrated into spf13/cobra, urfave/cli or any cli go application.
//
// promptui has two main input modes:
//
// Prompt provides a single line for user input. It supports optional live validation,
// confirmation and masking the input.
//
// Select provides a list of options to choose from. It supports pagination, search,
// detailed view and custom templates.
package promptui

import "errors"

// ErrEOF is the error returned from prompts when EOF is encountered.
var ErrEOF = errors.New("^D")

// ErrInterrupt is the error returned from prompts when an interrupt (ctrl-c) is
// encountered.
var ErrInterrupt = errors.New("^C")

// ErrInvalidInput is the error returned when confirm prompts are supplied input different from y/n
var ErrInvalidInput = errors.New("invalid input")

// ValidateFunc is a placeholder type for any validation functions that validates a given input. It should return
// a ValidationError if the input is not valid.
type ValidateFunc func(string) error
