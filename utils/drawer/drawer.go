package drawer

import (
	"fmt"
	"strings"

	"github.com/nsf/termbox-go"
)

// Drawer represents the Termbox-based UI drawer.
// It maintains the current cursor position for drawing operations.
type Drawer struct {
	x int // x is the current horizontal position (column)
	y int // y is the current vertical position (row)
}

// DrawerOption defines drawing options for terminal output operations.
// It controls line breaks, buffer flushing, and text colors.
type DrawerOption struct {
	newLine bool              // determines if cursor moves to next line after drawing
	flush   bool              // determines if buffer should be flushed after drawing
	fg      termbox.Attribute // specifies the foreground color attribute
	bg      termbox.Attribute // specifies the background color attribute
}

// NewDrawerOption creates a new DrawerOption with the given parameters
func NewDrawerOption(newLine, flush bool, fg, bg termbox.Attribute) *DrawerOption {
	return &DrawerOption{
		newLine: newLine,
		flush:   flush,
		fg:      fg,
		bg:      bg,
	}
}

// NewDefaultDrawerOption creates a new DrawerOption with default values
func NewDefaultDrawerOption() *DrawerOption {
	return &DrawerOption{
		newLine: true,
		flush:   true,
		fg:      termbox.ColorDefault,
		bg:      termbox.ColorDefault,
	}
}

// NewDefaultDrawerOptionNoFlush creates a new DrawerOption with default values and no flush
func NewDefaultDrawerOptionNoFlush() *DrawerOption {
	return &DrawerOption{
		newLine: true,
		flush:   false,
		fg:      termbox.ColorDefault,
		bg:      termbox.ColorDefault,
	}
}

// WithNewLine sets the newLine option for the DrawerOption
func (opt *DrawerOption) WithNewLine(newLine bool) *DrawerOption {
	opt.newLine = newLine
	return opt
}

// WithFlush sets the flush option for the DrawerOption
func (opt *DrawerOption) WithFlush(flush bool) *DrawerOption {
	opt.flush = flush
	return opt
}

// WithFg sets the foreground color attribute for the DrawerOption
func (opt *DrawerOption) WithFg(fg termbox.Attribute) *DrawerOption {
	opt.fg = fg
	return opt
}

// WithBg sets the background color attribute for the DrawerOption
func (opt *DrawerOption) WithBg(bg termbox.Attribute) *DrawerOption {
	opt.bg = bg
	return opt
}

// Default drawing options
var (
	// DefaultOption provides standard drawing options with buffer flush
	DefaultOption = NewDefaultDrawerOption()

	// DefaultOptionNoFlush provides standard drawing options without buffer flush
	DefaultOptionNoFlush = NewDefaultDrawerOptionNoFlush()
)

// Global instance of Drawer
var instance *Drawer

// Init initializes the Termbox environment and creates a new drawer instance.
// Returns an error if termbox initialization fails.
func Init() error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	instance = &Drawer{x: 0, y: 0}
	return nil
}

// Close cleans up the Termbox environment and releases resources.
// Should be called when the drawer is no longer needed.
func Close() {
	if instance != nil {
		termbox.Close()
		instance = nil
	}
}

// Flush forces the terminal to display all pending drawing operations.
func Flush() {
	termbox.Flush()
}

// Render draws a string at the specified coordinates (x, y).
// The string is drawn with default colors and the buffer is flushed immediately.
func Render(y int, x int, str string) {
	for i, ch := range str {
		termbox.SetCell(x+i, y, ch, termbox.ColorDefault, termbox.ColorDefault)
	}
	termbox.Flush()
}

// Print displays a message at the current cursor position with specified options.
// Handles multi-line strings and updates cursor position accordingly.
func Print(message string, opt *DrawerOption) {
	lines := strings.Split(message, "\n")
	for _, line := range lines {
		for _, ch := range line {
			termbox.SetCell(instance.x, instance.y, ch, opt.fg, opt.bg)
			instance.x++
		}
		instance.y++
		instance.x = 0
	}

	if !opt.newLine {
		instance.y--
	}

	if opt.flush {
		termbox.Flush()
	}
}

// Clear clears the entire terminal screen and resets cursor position.
// Uses specified background and foreground colors from the option.
func Clear(opt *DrawerOption) {
	instance.x = 0
	instance.y = 0

	termbox.Clear(opt.bg, opt.fg)
	if opt.flush {
		termbox.Flush()
	}
}

// DrawProgressBar draws a progress bar at specified Y coordinate with given percentage.
// The progress bar includes a percentage indicator and uses the full terminal width.
func DrawProgressBar(y int, percent int, opt *DrawerOption) {
	width, _ := termbox.Size()
	totalWidth := width - len("  100/100") - 4
	fillWidth := int(float64(totalWidth) * float64(percent) / 100)

	// Draw the left boundary of the progress bar
	termbox.SetCell(0, y, '|', opt.fg, opt.bg)

	// Draw the completed part
	for i := 1; i < fillWidth; i++ {
		termbox.SetCell(i, y, '=', opt.fg, opt.bg)
	}

	// Draw the remaining part
	for i := fillWidth; i < totalWidth; i++ {
		termbox.SetCell(i, y, '-', opt.fg, opt.bg)
	}

	// Draw the right boundary of the progress bar
	termbox.SetCell(totalWidth, y, '|', opt.fg, opt.bg)

	// Draw the progress percentage text
	progressText := fmt.Sprintf(" %d/100", percent)
	for i, r := range progressText {
		termbox.SetCell(totalWidth+1+i, y, r, opt.fg, opt.bg)
	}

	// Flush if necessary
	if opt.flush {
		termbox.Flush()
	}
}

// GetY returns the current vertical cursor position.
func GetY() int {
	return instance.y
}

// GetX returns the current horizontal cursor position.
func GetX() int {
	return instance.x
}

// NextLine moves the cursor to the beginning of the next line.
func NextLine() {
	instance.y++
	instance.x = 0
}
