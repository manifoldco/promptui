// +build !windows

package promptui

import "github.com/chzyer/readline"

var (
	KeyEnter     rune = readline.CharEnter
	KeyBackspace rune = readline.CharBackspace
	KeyPrev      rune = readline.CharPrev
	KeyNext      rune = readline.CharNext
	KeyBackward  rune = readline.CharBackward
	KeyForward   rune = readline.CharForward
)
