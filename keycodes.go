// +build !windows

package promptui

import "github.com/chzyer/readline"

var (
	// KeyEnter is the default key for submission/selection
	KeyEnter rune = readline.CharEnter

	// KeyBackspace is the default key for deleting input text
	KeyBackspace rune = readline.CharBackspace

	// KeyPrev is the default key to go up during selection
	KeyPrev rune = readline.CharPrev

	// KeyNext is the default key to go down during selection
	KeyNext rune = readline.CharNext

	// KeyBackward is the default key to page up during selection
	KeyBackward rune = readline.CharBackward

	// KeyForward is the default key to page down during selection
	KeyForward rune = readline.CharForward
)
