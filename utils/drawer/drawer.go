package drawer

import (
	"fmt"
	"strings"

	"github.com/nsf/termbox-go"
)

// Drawer represents the Termbox-based UI drawer.
type Drawer struct {
	X int // X is the current horizontal position
	Y int // Y is the current vertical position
}

type Option struct {
	NewLine bool              // NewLine is whether to start a new line
	Flush   bool              // Flush is whether to flush the buffer
	Fg      termbox.Attribute // Fg is the foreground color
	Bg      termbox.Attribute // Bg is the background color
}

var (
	DefaultOption = Option{
		NewLine: true,
		Flush:   true,
		Fg:      termbox.ColorDefault,
		Bg:      termbox.ColorDefault,
	}

	DefaultOptionNoFlush = Option{
		NewLine: true,
		Flush:   false,
		Fg:      termbox.ColorDefault,
		Bg:      termbox.ColorDefault,
	}
)

// Global instance of Drawer
var instance *Drawer

// Init initializes the Termbox environment and the drawer.
func Init() error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	instance = &Drawer{X: 0, Y: 0}
	return nil
}

// Close cleans up the Termbox environment.
func Close() {
	if instance != nil {
		termbox.Close()
		instance = nil
	}
}

// Flush flushes the buffer.
func Flush() {
	termbox.Flush()
}

// Render draws a string at the specified Y coordinate.
func Render(y int, x int, str string) {
	for i, ch := range str {
		termbox.SetCell(x+i, y, ch, termbox.ColorDefault, termbox.ColorDefault)
	}
	termbox.Flush()
}

// Print displays a message at the current position.
func Print(message string, option Option) {
	lines := strings.Split(message, "\n")
	for _, line := range lines {
		for _, ch := range line {
			termbox.SetCell(instance.X, instance.Y, ch, option.Fg, option.Bg)
			instance.X++
		}
		instance.Y++
		instance.X = 0
	}

	if !option.NewLine {
		instance.Y--
	}

	if option.Flush {
		termbox.Flush()
	}
}

// Clear clears the message area.
func Clear(option Option) {
	instance.X = 0
	instance.Y = 0

	termbox.Clear(option.Bg, option.Fg)
	if option.Flush {
		termbox.Flush()
	}
}

// DrawProgressBar draws a progress bar at the specified Y coordinate.
func DrawProgressBar(y int, percent int, option Option) {
	width, _ := termbox.Size()
	totalWidth := width - len("  100/100") - 4
	fillWidth := int(float64(totalWidth) * float64(percent) / 100)

	// Draw the left boundary of the progress bar
	termbox.SetCell(0, y, '|', option.Fg, option.Bg)

	// Draw the completed part
	for i := 1; i < fillWidth; i++ {
		termbox.SetCell(i, y, '=', option.Fg, option.Bg)
	}

	// Draw the remaining part
	for i := fillWidth; i < totalWidth; i++ {
		termbox.SetCell(i, y, '-', option.Fg, option.Bg)
	}

	// Draw the right boundary of the progress bar
	termbox.SetCell(totalWidth, y, '|', option.Fg, option.Bg)

	// Draw the progress percentage text
	progressText := fmt.Sprintf(" %d/100", percent)
	for i, r := range progressText {
		termbox.SetCell(totalWidth+1+i, y, r, option.Fg, option.Bg)
	}

	// Flush if necessary
	if option.Flush {
		termbox.Flush()
	}
}

// GetY returns the current vertical position.
func GetY() int {
	return instance.Y
}

// GetX returns the current horizontal position.
func GetX() int {
	return instance.X
}

// NextLine moves the cursor to the next line.
func NextLine() {
	instance.Y++
	instance.X = 0
}
