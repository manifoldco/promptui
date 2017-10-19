// +build !windows

package promptui

var (
	bold  = Styler(FGBold)
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

var red = Styler(FGBold, FGRed)
