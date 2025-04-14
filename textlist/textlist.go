package textlist

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"iter"
	"strings"
	"time"
	//"time"
)

/*

  File:    textlist.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: TextList extends widget.List for a tabbed, text editing component.
*/

type TextList struct {
	widget.List

	Err    error
	Theme  MyTheme
	window fyne.Window

	rows               []listRow
	rowId              int
	startMark, endMark int

	controlBox *fyne.Container

	mode int

	edit     *widget.Entry
	editBox  *fyne.Container
	editText string

	search        *widget.Entry
	replace       *widget.Entry
	replaceAction *widget.Button
	searchBox     *fyne.Container
	searchText    string
	replaceText   string
	results       []result
	currentResult int
	resultCount   *widget.Label
	down          *widget.Button
	up            *widget.Button
	ignoreCase    bool

	style            *fyne.TextStyle
	spaces           string
	lineFormat       string
	searchLineFormat string
	charX, charY     float32
}

type result struct {
	rowId int
	col1  int
	col2  int
	text  string
}

// NewTextList creates a container with a TextList "widget"
func NewTextList(window fyne.Window, name string, buttonBar *fyne.Container,
	theme MyTheme) (*TextList, *fyne.Container) {

	l := &TextList{
		window:    window,
		Theme:     theme,
		style:     &theme.style,
		spaces:    strings.Repeat(" ", theme.tabSize),
		startMark: -1,
		endMark:   -1,
		mode:      modeEdit,
	}
	l.ExtendBaseWidget(l)
	l.HideSeparators = true

	// delegate List functions
	l.Length = func() int {
		return len(l.rows) + 1
	}
	l.CreateItem = func() fyne.CanvasObject {
		return l.createItem()
	}
	l.UpdateItem = func(id widget.ListItemID, item fyne.CanvasObject) {
		l.updateItem(id, item)
	}
	l.OnSelected = func(rowId int) {
		l.onSelected(rowId)
	}

	l.views(buttonBar)

	if len(name) > 0 {
		label := widget.NewLabel(name)
		label.Alignment = fyne.TextAlignTrailing
		sep := canvas.NewLine(l.Theme.Color("normalColor", 0))
		return l, container.NewBorder(container.NewVBox(label, sep), l.controlBox, nil, nil, l)
	}

	return l, container.NewBorder(nil, l.controlBox, nil, nil, l)
}

// SetContent adds all the strings (separated by \n)
func (l *TextList) SetContent(content string) {
	l.setContent(content)
}

// AddString adds the string to the end of the list
func (l *TextList) AddString(str string) {
	l.addString(str)
}

// Start starts the TextList in edit mode
func (l *TextList) Start() {
	l.showEdit()
}

// Count returns the current number of text rows
func (l *TextList) Count() int {
	return len(l.rows)
}

// GetContent returns a []string for all the rows
func (l *TextList) GetContent() []string {
	content := make([]string, len(l.rows))
	for ix := range l.rows {
		content[ix] = l.getRowString(ix)
	}
	return content
}

// Iterator provides a GO 1.23 range operator
func Iterator(l *TextList) iter.Seq[string] {
	return func(yield func(string) bool) {
		for i := 0; i < len(l.rows); i++ {
			str := l.getRowString(i)
			if !yield(str) { // return str, quit if scope has ended
				return
			}
		}
	}
}

// toast displays a message in lower right of the canvas
func (l *TextList) toast(txt string, c color.NRGBA, d time.Duration) {
	toast(l.window.Canvas(), fyne.Position{
		X: l.window.Canvas().Size().Width,
		Y: l.window.Canvas().Size().Height,
	}, txt, c, d)

}

// focus sets the Canvas.focus to a particular item
func (l *TextList) focus(f fyne.Focusable) {
	l.window.Canvas().Focus(f)
}

const (
	modeEdit   = 0
	modeSearch = 1
)

func (l *TextList) TypedShortcut(s fyne.Shortcut) {

	var collectRows = func(rows []string) string {
		str := ""
		for sx, s := range rows {
			str += s
			if sx < len(rows)-1 {
				str += "\n"
			}
		}
		return str
	}

	switch s.ShortcutName() {

	case "CustomDesktop:Control+Home":
		l.UnselectAll()
		l.ScrollToTop()
		l.rowId = 0
	case "CustomDesktop:Control+End":
		l.UnselectAll()
		l.ScrollToBottom() // doesn't go to the "end"
		l.ScrollToBottom() // do it twice?
		l.rowId = len(l.rows) - 1
	case "CustomDesktop:Control+Down", "CustomDesktop:Control+Next", "CustomDesktop:Alt+Down":
		l.pageDown(10)
	case "CustomDesktop:Control+Up", "CustomDesktop:Control+Prior", "CustomDesktop:Alt+Up":
		l.pageUp(10)

	case "CustomDesktop:Control+F", "CustomDesktop:Control+R":
		if l.mode == modeEdit {
			l.showSearch()
		}

	case "CustomDesktop:Control+M":
		if l.mode == modeEdit {
			l.markStartRow(l.rowId)
		}
	case "CustomDesktop:Control+E":
		if l.mode == modeEdit {
			l.markEndRow(l.rowId)
		}

	case "Cut": // Ctrl+X
		if l.mode == modeEdit {
			str := collectRows(l.deleteMarkedRows())
			cb := l.window.Clipboard()
			cb.SetContent(str)
		}
	case "Copy": // Ctrl+C
		if l.mode == modeEdit {
			str := collectRows(l.getMarkedRows())
			cb := l.window.Clipboard()
			cb.SetContent(str)
		}
	case "Paste": // Ctrl+V
		if l.mode == modeEdit {
			cb := l.window.Clipboard()
			str := cb.Content()
			l.insertRows(str)
		}

	}
}
