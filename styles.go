// +build !windows

package promptui

var (
	faint = Styler(FGFaint)
)

// Icons used for displaying prompts or status
var (
	IconInitial = Styler(FGBlue)("?")
	IconGood    = Styler(FGGreen)("✔")
	IconWarn    = Styler(FGYellow)("⚠")
	IconBad     = Styler(FGRed)("✗")
	IconSelect  = Styler(FGBold)("▸")
)
