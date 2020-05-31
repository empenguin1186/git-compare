package main

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

// Activity : activity
type Activity struct {
	selectRow   int
	items       []string
	itemPerPage int
}

// NewActivity :
func NewActivity(items []string) *Activity {
	return &Activity{
		selectRow:   1,
		items:       items,
		itemPerPage: len(items),
	}
}

// ChooseCommit :
func (a *Activity) ChooseCommit(title string) int {
	var result int

	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	defer termbox.Close()

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.HideCursor()

	a.refreshPage(title)
	a.setHighlight(1)
	termbox.Flush()

end:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventResize:
			a.refreshPage(title)
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				a.selectRow = 0
				result = -1
				break end
			case termbox.KeyEnter:
				result = a.selectRow - 1
				break end
			case termbox.KeyArrowDown:
				if a.selectRow < len(a.items) {
					a.selectRow++
				} else {
					a.selectRow = 1
				}
				a.setHighlight(a.selectRow)
			case termbox.KeyArrowUp:
				if a.selectRow == 1 {
					a.selectRow = len(a.items)
				} else {
					a.selectRow--
				}
				a.setHighlight(a.selectRow)
			default:
				break
			}
		default:
			break
		}
	}

	return result
}

func (a *Activity) refreshPage(header string) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	_, height := termbox.Size()

	a.itemPerPage = height - 2

	start := 0
	end := start + a.itemPerPage
	if len(a.items[start:]) < a.itemPerPage {
		end = start + len(a.items[start:])
	}

	a.items = a.items[start:end]

	for index, item := range a.items {
		a.writeLine(index+1, item)
	}
	a.resetHeader(header)
	a.setHighlight(1)
	termbox.Flush()
}

func (a *Activity) resetHeader(header string) {
	a.centerWrite(0, header, termbox.ColorMagenta, termbox.ColorCyan)
}

func (a *Activity) centerWrite(row int, text string, textColor termbox.Attribute, bgColor termbox.Attribute) {
	width, _ := termbox.Size()
	fixing := (width / 2) - (len(text) / 2)
	for col := 0; col < width; col++ {
		char := ' '
		if col >= fixing && col < fixing+(len(text)) {
			char = rune(text[col-fixing])
		}
		termbox.SetCell(col, row, char, textColor, bgColor)
	}
}

func (a *Activity) writeLine(row int, str string) {
	runes := []rune(str)
	x := 0
	for _, r := range runes {
		termbox.SetCell(x, row, r, termbox.ColorWhite, termbox.ColorDefault)
		x += runewidth.RuneWidth(r)
	}
}

func (a *Activity) setHighlight(targetRow int) {
	width, _ := termbox.Size()
	a.selectRow = targetRow
	for row := 1; row <= a.itemPerPage; row++ {
		bgColor := termbox.ColorDefault
		if row == targetRow {
			bgColor = termbox.ColorWhite
		}
		for col := 0; col < width; col++ {
			char := termbox.CellBuffer()[(width*row)+col].Ch
			termbox.SetCell(col, row, char, termbox.ColorDefault, bgColor)
		}
	}
	termbox.Flush()
}
