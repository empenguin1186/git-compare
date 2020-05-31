package main

import (
	"strconv"

	"github.com/nsf/termbox-go"
)

// Activity : activity
type Activity struct {
	selectRow   int
	itemPerPage int
	currentPage int
	allItems    []string
	pageItems   []string
	numOfPages  int
}

// NewActivity :
func NewActivity(itemPerPage int, allItems []string) *Activity {
	return &Activity{
		selectRow:   1,
		itemPerPage: itemPerPage,
		currentPage: 0,
		allItems:    allItems,
		pageItems:   []string{},
		numOfPages:  0,
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
				result = (a.itemPerPage * a.currentPage) + (a.selectRow - 1)
				break end
			case termbox.KeyArrowLeft:
				if a.currentPage != 0 {
					a.currentPage--
				}
				a.refreshPage(title)
			case termbox.KeyArrowRight:
				if a.numOfPages > a.currentPage+1 {
					a.currentPage++
				}
				a.refreshPage(title)
			case termbox.KeyArrowDown:
				if a.selectRow < len(a.pageItems) {
					a.selectRow++
				} else {
					a.selectRow = 1
				}
				a.setHighlight(a.selectRow)
			case termbox.KeyArrowUp:
				if a.selectRow == 1 {
					a.selectRow = len(a.pageItems)
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

	a.numOfPages = len(a.allItems) / a.itemPerPage
	if len(a.allItems)%a.itemPerPage != 0 {
		a.numOfPages++
	}

	start := a.currentPage * a.itemPerPage
	end := start + a.itemPerPage
	if len(a.allItems[start:]) < a.itemPerPage {
		end = start + len(a.allItems[start:])
	}

	a.pageItems = a.allItems[start:end]

	for index, item := range a.pageItems {
		a.writeLine(index+1, item)
	}
	a.resetHeader(header)
	a.resetFooter()
	a.setHighlight(1)
	termbox.Flush()
}

func (a *Activity) resetHeader(header string) {
	a.centerWrite(0, header, termbox.ColorMagenta, termbox.ColorCyan)
}

func (a *Activity) resetFooter() {
	_, bottom := termbox.Size()
	a.numOfPages = len(a.allItems) / a.itemPerPage
	if len(a.allItems)%a.itemPerPage != 0 {
		a.numOfPages++
	}
	prevText := "Prev <- "
	if a.currentPage == 0 {
		prevText = "        "
	}
	nextText := " -> Next"
	if a.currentPage+1 == a.numOfPages {
		nextText = "        "
	}
	footer := prevText + strconv.Itoa(a.currentPage+1) + "/" + strconv.Itoa(a.numOfPages) + nextText
	// termbox.co
	a.centerWrite(bottom-1, footer, termbox.ColorBlack, termbox.ColorCyan)
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

	width, _ := termbox.Size()
	for col := 0; col < width; col++ {
		char := ' '
		if col < len(str) {
			// rune は 文字列を byte 単位でなく文字単位で扱う場合に使用
			// https://text.baldanders.info/golang/string-and-rune/
			char = rune(str[col])
		}
		termbox.SetCell(col, row, char, termbox.ColorDefault, termbox.ColorDefault)
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
