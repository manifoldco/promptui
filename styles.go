// +build !windows

package promptui

import "fmt"

var (
	bold       = Styler(FGBold)
	faint      = Styler(FGFaint)
	underlined = Styler(FGUnderline)
)

// Icons used for displaying prompts or status
var (
	IconInitial = Styler(FGBlue)("?")
	IconGood    = Styler(FGGreen)("✔")
	IconWarn    = Styler(FGYellow)("⚠")
	IconBad     = Styler(FGRed)("✗")
)

var red = Styler(FGBold, FGRed)

func style(name string, value interface{}) string {
	s := fmt.Sprintf("%v", value)
	switch name {
	case "red":
		return red(s)
	case "bold":
		return bold(s)
	default:
		return s
	}
}
