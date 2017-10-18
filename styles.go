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

func color(name string, value interface{}) string {
	s := fmt.Sprintf("%v", value)
	if name == "red" {
		return red(s)
	}
	return s
}
