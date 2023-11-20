package utilsTermbox

import (
	"fmt"
	"strings"

	"github.com/nsf/termbox-go"
)

// TermboxDrawer represents the drawer for Termbox-based UI.
var Td TermboxDrawer

// InIt initializes the TermboxDrawer.
func InIt() {
	Td = *NewTermboxDrawer()
}

// NewTermboxDrawer creates and returns a new TermboxDrawer.
func NewTermboxDrawer() *TermboxDrawer {
	return &TermboxDrawer{xTermbox: 2, yTermbox: 0}
}

// DrawStringAtY draws a string at the specified Y coordinate.
func (td *TermboxDrawer) DrawStringAtY() func(x int, str string) {
	y := td.yTermbox
	td.yTermbox++
	return func(x int, str string) {
		for i, ch := range str {
			termbox.SetCell(x+i, y, ch, termbox.ColorDefault, termbox.ColorDefault)
		}
		termbox.Flush()
	}
}

// PrintMessage prints a message at the current position.
func (td *TermboxDrawer) PrintMessage(message string, options ...MessageOption) {
	option := MessageOption{
		NewLine: true,
		Flush:   true,
	}
	if len(options) > 0 {
		option = options[0]
	}
	lines := strings.Split(message, "\n")
	for _, line := range lines {
		for _, ch := range line {
			termbox.SetCell(td.xTermbox+2, td.yTermbox, ch, termbox.ColorDefault, termbox.ColorDefault)
			td.xTermbox++
		}
		td.yTermbox++
		td.xTermbox = 2
	}

	if !option.NewLine {
		td.yTermbox--
	}

	if option.Flush {
		termbox.Flush()
	}
}

// ClearMessage clears the message area.
func (td *TermboxDrawer) ClearMessage(options ...Option) {
	td.xTermbox = 2
	td.yTermbox = 0
	option := Option{
		Flush: true,
	}
	if len(options) > 0 {
		option = options[0]
	}
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	if option.Flush {
		termbox.Flush()
	}
}

// ProgressBarAtY draws a progress bar at the specified Y coordinate.
func (td *TermboxDrawer) ProgressBarAtY() func(percent int, options ...Option) {
	y := td.yTermbox
	td.yTermbox++
	return func(percent int, options ...Option) {
		w, _ := termbox.Size()
		totalWidth := w - 2 - len(" 100/100")
		fillWidth := int(float64(totalWidth) * float64(percent) / 100)

		option := Option{
			Flush: true,
		}

		if len(options) > 0 {
			option = options[0]
		}

		// Draw the left boundary of the progress bar
		termbox.SetCell(0, y, '|', termbox.ColorWhite, termbox.ColorDefault)

		// Draw the completed part
		for i := 0; i < fillWidth; i++ {
			termbox.SetCell(i+1, y, '=', termbox.ColorGreen, termbox.ColorDefault)
		}

		// Draw the remaining part
		for i := fillWidth; i < totalWidth; i++ {
			termbox.SetCell(i+1, y, '-', termbox.ColorDefault, termbox.ColorDefault)
		}

		// Draw the right boundary of the progress bar
		termbox.SetCell(totalWidth+1, y, '|', termbox.ColorWhite, termbox.ColorDefault)

		// Draw the progress percentage text
		progressText := fmt.Sprintf(" %d/100", percent)
		for i, r := range progressText {
			termbox.SetCell(totalWidth+2+i, y, r, termbox.ColorWhite, termbox.ColorDefault)
		}

		// Flush if necessary
		if option.Flush {
			termbox.Flush()
		}
	}
}

// GetX returns the X coordinate of the TermboxDrawer.
func (td *TermboxDrawer) GetX() int {
	return td.xTermbox
}

// GetY returns the Y coordinate of the TermboxDrawer.
func (td *TermboxDrawer) GetY() int {
	return td.yTermbox
}
