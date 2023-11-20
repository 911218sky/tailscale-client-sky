package utilsTermbox

// Drawer defines an interface with drawing methods.
type Drawer interface {
	InIt()
	DrawString(x, y int, str string)
	PrintMessage(message string, options ...MessageOption)
	ClearMessage(options ...Option)
	PrintProgressBar(percent int, options ...Option)
}

// Option represents options for drawing operations.
type Option struct {
	Flush bool // Flush indicates whether to immediately update the display.
}

// MessageOption represents options for message printing.
type MessageOption struct {
	NewLine bool // NewLine specifies whether to move to a new line after printing.
	Flush   bool // Flush indicates whether to immediately update the display.
}

// TermboxDrawer implements the Drawer interface for drawing text and progress bars in termbox.
type TermboxDrawer struct {
	xTermbox int // xTermbox represents the X coordinate in termbox.
	yTermbox int // yTermbox represents the Y coordinate in termbox.
}
