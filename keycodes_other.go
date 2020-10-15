// +build !windows

package promptui

import "github.com/chzyer/readline"

var (
	// KeyBackspace is the key for deleting the previous char.
	KeyBackspace rune = readline.CharBackspace

	// KeyDelete is the key for deleting the next char.
	KeyDelete rune = readline.CharDelete
)
