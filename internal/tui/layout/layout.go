package layout

// UI defines fixed sections of the Luna TUI.
type UI struct {
	Width         int
	WelcomeWidth  int
	ComposerWidth int
}

// Compute returns the current UI layout based on terminal width.
func Compute(width int) UI {
	if width <= 0 {
		width = 80
	}
	return UI{
		Width:         width,
		WelcomeWidth:  width,
		ComposerWidth: width,
	}
}
