package promptui

// Icons used for displaying prompts or status
var (
	IconInitial = Styler(FGBlue)("?")
	IconGood    = Styler(FGGreen)("v")
	IconWarn    = Styler(FGYellow)("!")
	IconBad     = Styler(FGRed)("x")
	IconSelect  = Styler(FGBold)(">")
)
