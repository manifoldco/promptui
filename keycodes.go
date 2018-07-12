// +build !windows

package promptui

import "github.com/chzyer/readline"

// These runes are used to identity the commands entered by the user in the command prompt. They map
// to specific actions of prompt-ui and can be remapped if necessary,
var (
	// KeyEnter is the default key for submission/selection inside a command line prompt.
	KeyEnter rune = readline.CharEnter

	// KeyBackspace is the default key for deleting input text inside a command line prompt.
	KeyBackspace rune = readline.CharBackspace

	// KeyPrev is the default key to go up during selection inside a command line prompt.
	KeyPrev rune = readline.CharPrev

	// KeyNext is the default key to go down during selection inside a command line prompt.
	KeyNext rune = readline.CharNext

	// KeyBackward is the default key to page up during selection inside a command line prompt.
	KeyBackward rune = readline.CharBackward

	// KeyForward is the default key to page down during selection inside a command line prompt.
	KeyForward rune = readline.CharForward
)
