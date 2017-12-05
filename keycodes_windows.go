package promptui

// source: https://msdn.microsoft.com/en-us/library/aa243025(v=vs.60).aspx

var (
	// KeyEnter is the default key for submission/selection
	KeyEnter rune = 13

	// KeyBackspace is the default key for deleting input text
	KeyBackspace rune = 8

	// FIXME: keys below are not triggered by readline, not working on Windows

	// KeyPrev is the default key to go up during selection
	KeyPrev rune = 38

	// KeyNext is the default key to go down during selection
	KeyNext rune = 40

	// KeyBackward is the default key to page up during selection
	KeyBackward rune = 37

	// KeyForward is the default key to page down during selection
	KeyForward rune = 39
)
