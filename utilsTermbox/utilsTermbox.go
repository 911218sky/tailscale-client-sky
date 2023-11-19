package utilsTermbox

import (
	"fmt"
	"strings"

	"github.com/nsf/termbox-go"
)

var (
	XTermbox     = 2
	YTermbox     = 0
	progressBarY = 0
)

type Option struct {
	NoNewLine bool
	NoFlush   bool
}

func DrawString(x, y int, str string) {
	for i, ch := range str {
		termbox.SetCell(x+i+2, y, ch, termbox.ColorDefault, termbox.ColorDefault)
	}
	termbox.Flush()
}

func PrintMessage(message string, options ...Option) {
	option := Option{}
	if len(options) > 0 {
		if options[0].NoFlush {
			option.NoFlush = true
		}
		if options[0].NoNewLine {
			option.NoNewLine = true
		}
	}

	lines := strings.Split(message, "\n")
	for _, line := range lines {
		for _, ch := range line {
			termbox.SetCell(XTermbox, YTermbox, ch, termbox.ColorDefault, termbox.ColorDefault)
			XTermbox++
		}
		YTermbox++
		XTermbox = 2
	}

	if option.NoNewLine {
		YTermbox--
	}

	if !option.NoFlush {
		termbox.Flush()
	}
}

func ClearMessage(options ...Option) {
	XTermbox = 2
	YTermbox = 0
	isFlush := true
	if len(options) > 0 {
		if options[0].NoFlush {
			isFlush = false
		}
	}
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	if isFlush {
		termbox.Flush()
	}
}

func ProgressBarInit() {
	progressBarY = YTermbox
	YTermbox++
}

func PrintProgressBar(percent int, options ...Option) {
	w, _ := termbox.Size()
	totalWidth := w - 2 - len(" 100/100")
	fillWidth := int(float64(totalWidth) * float64(percent) / 100)
	isFlush := true
	if len(options) > 0 {
		if options[0].NoFlush {
			isFlush = false
		}
	}

	// 绘制进度条左边界
	termbox.SetCell(0, progressBarY, '|', termbox.ColorWhite, termbox.ColorDefault)

	// 绘制已完成部分
	for i := 0; i < fillWidth; i++ {
		termbox.SetCell(i+1, progressBarY, '=', termbox.ColorGreen, termbox.ColorDefault)
	}

	// 绘制未完成部分
	for i := fillWidth; i < totalWidth; i++ {
		termbox.SetCell(i+1, progressBarY, '-', termbox.ColorDefault, termbox.ColorDefault)
	}

	// 绘制进度条右边界
	termbox.SetCell(totalWidth+1, progressBarY, '|', termbox.ColorWhite, termbox.ColorDefault)

	// 绘制进度百分比文本
	progressText := fmt.Sprintf(" %d/100", percent)
	for i, r := range progressText {
		termbox.SetCell(totalWidth+2+i, progressBarY, r, termbox.ColorWhite, termbox.ColorDefault)
	}

	if isFlush {
		termbox.Flush()
	}
}
