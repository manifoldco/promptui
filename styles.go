// +build !windows

package promptui

// Icons used for displaying prompts or status
var (
	IconInitial = Styler(FGBlue)("?")
	IconGood    = Styler(FGGreen)("✔")
	IconWarn    = Styler(FGYellow)("⚠")
	IconBad     = Styler(FGRed)("✗")
	IconSelect  = Styler(FGBold)("▸")
)
