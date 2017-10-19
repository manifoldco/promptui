package promptui

var (
	bold  = Styler(FGBold)
	faint = Styler(FGFaint)
)

// Icons used for displaying prompts or status
var (
	IconInitial = Styler(FGBlue)("?")
	IconGood    = Styler(FGGreen)("v")
	IconWarn    = Styler(FGYellow)("!")
	IconBad     = Styler(FGRed)("x")
)

var red = Styler(FGBold, FGRed)
