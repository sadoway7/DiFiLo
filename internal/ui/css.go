package ui

// DifiCSS returns the complete application stylesheet, combining all
// CSS sub-sections into a single <style> block.
func DifiCSS() string {
	return "<style>\n" + cssBase + cssNav + cssPages + cssHome + cssComponents + "\n</style>"
}
