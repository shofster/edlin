package textlist

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"math"
	"strings"
	"unicode"
)

/*

  File:    rowlist.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: rowlist implements widget.List functions for managing a TextList
*/

type listCell struct {
	r      rune
	style  *fyne.TextStyle
	marked bool
}

type listRow struct {
	cells  []listCell
	style  *fyne.TextStyle
	marked bool
}

func (l *TextList) setContent(content string) {
	s := strings.Split(content, "\n")
	for i := 0; i < len(s); i++ {
		l.addString(s[i])
	}
}

func (l *TextList) addString(str string) {
	row := l.createRow(str)
	l.rows = append(l.rows, row)
	lb10 := int(math.Log10(float64(len(l.rows)))) + 1
	l.lineFormat = fmt.Sprintf(" %%%dd  ", lb10)
	l.searchLineFormat = fmt.Sprintf(" %%%dd%s ", lb10, "\u2192")
	l.Refresh()
}

// createItem called by List to create empty CanvasObject
func (l *TextList) createItem() fyne.CanvasObject {
	return container.New(layout.NewCustomPaddedHBoxLayout(0))
}

// updateItem called by List to generate the current CanvasObject
func (l *TextList) updateItem(id widget.ListItemID, item fyne.CanvasObject) {

	item.(*fyne.Container).Objects = nil
	item.(*fyne.Container).Objects = append(item.(*fyne.Container).Objects,
		l.lineNo(id))

	runes := l.getRowRunes(id, false)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		cell := l.rows[id].cells[i]

		var text *canvas.Text
		var style = l.rows[id].cells[i].style

		if style == nil {
			style = l.rows[id].style
		}
		if style == nil {
			style = l.style
		}

		if cell.marked || l.rows[id].marked {
			text = canvas.NewText(string(r), l.Theme.Color("markedColor", 0))
		} else {
			text = canvas.NewText(string(r), l.Theme.Color("normalColor", 0))
		}

		text.Alignment = fyne.TextAlignCenter
		text.TextSize = l.Theme.textSize
		text.TextStyle = *style

		item.(*fyne.Container).Objects = append(item.(*fyne.Container).Objects, text)
	}
	l.SetItemHeight(id, l.Theme.textSize)

}

func (l *TextList) onSelected(rowId widget.ListItemID) {
	l.moveToRow(rowId)
	l.focus(l)
}

func (l *TextList) moveToRow(rowId int) {
	l.UnselectAll()
	l.ScrollTo(rowId)
	l.rowId = rowId
	l.editText = l.getRowString(rowId)
	l.edit.SetText(l.editText)
	l.Refresh()
}

func (l *TextList) pageDown(n int) {
	rowId := l.rowId + n
	rowId = min(rowId, len(l.rows)-1)
	off := float32(rowId) * l.Theme.textSize
	l.ScrollToOffset(off)
	l.rowId = rowId
}

func (l *TextList) pageUp(n int) {
	rowId := l.rowId - n
	rowId = max(rowId, 0)
	l.rowId = rowId
	off := float32(rowId) * l.Theme.textSize
	l.ScrollToOffset(off)
}

func (l *TextList) clearMarkedRows(doCells bool) {
	for rowId := 0; rowId < len(l.rows); rowId++ {
		l.rows[rowId].marked = false
		if doCells {
			f := result{
				rowId: rowId,
				col1:  0,
				col2:  len(l.rows[rowId].cells) - 1,
			}
			l.markCells(f, false)
		}
	}
	l.startMark = -1
	l.endMark = -1
}

// markStartRow clears marks and marks the beginning of a new start row
func (l *TextList) markStartRow(rowId int) {

	l.clearMarkedRows(true)
	if rowId >= len(l.rows)-1 {
		return
	}

	l.startMark = rowId
	l.rows[rowId].marked = true
	l.endMark = -1
	l.Refresh()
}

// markEndRow marks all rows between (inclusive) start and rowId
func (l *TextList) markEndRow(rowId int) {

	// may not be in mark mode, or invalid rowId
	if l.startMark == -1 || rowId >= len(l.rows) {
		return
	}
	if l.endMark != -1 { // moving endMark
		sm := l.startMark // clear marks and restart mark
		l.markStartRow(sm)
	}
	l.endMark = rowId
	inc := 1
	if rowId > l.startMark {
		inc = -1
	}
mark:
	if rowId == l.startMark {
		l.Refresh()
		return
	}
	l.rows[rowId].marked = true
	rowId += inc
	goto mark
}

func (l *TextList) getMarkedRows() (rows []string) {
	if l.endMark == -1 {
		rows = append(rows, l.getRowString(l.rowId))
		return
	}
	start := l.startMark
	end := l.endMark
	if end < start {
		start = l.endMark
		end = l.startMark
	}
	for rowId := start; rowId <= end; rowId++ {
		rows = append(rows, l.getRowString(rowId))
	}
	return
}

func (l *TextList) deleteMarkedRows() (rows []string) {
	if l.endMark == -1 {
		l.deleteRow(l.rowId)
	} else {
		start := l.startMark
		end := l.endMark
		if end < start {
			start = l.endMark
			end = l.startMark
		}
		for rowId := start; rowId <= end; rowId++ {
			rows = append(rows, l.getRowString(rowId))
		}
		n := end - start + 1
		for i := 0; i < n; i++ {
			l.deleteRow(start)
		}
	}
	l.startMark = -1
	l.endMark = -1
	return
}

func (l *TextList) insertRows(content string) {
	rowId := l.rowId
	s := strings.Split(content, "\n")
	replace := make([]listRow, len(s))
	for i := 0; i < len(s); i++ {
		replace[i] = l.createRow(s[i])
	}
	rows := append(l.rows[:rowId], append(replace, l.rows[rowId:]...)...)
	l.rows = rows
	l.Refresh()
}

func (l *TextList) deleteRow(rowId int) {
	switch {
	case rowId >= len(l.rows): // last row
	default:
		rows := append(l.rows[:rowId], l.rows[rowId+1:]...)
		l.rows = rows
		l.moveToRow(rowId)
		l.UnselectAll()
	}
}

func (l *TextList) replaceRow(rowId int, content string) {
	lastRow := len(l.rows) - 1
	s := strings.Split(content, "\n")
	replace := make([]listRow, len(s))
	rows := make([]listRow, 0)
	for i := 0; i < len(s); i++ {
		replace[i] = l.createRow(s[i])
	}
	if rowId >= lastRow {
		rows = append(l.rows[:rowId], replace...)
	} else {
		rows = append(l.rows[:rowId], append(replace, l.rows[rowId+1:]...)...)
	}
	l.rows = rows
	if rowId == lastRow {
		l.AddString("")
	}
}

func (l *TextList) replaceCells(rowId, col1, col2 int, runes []rune, marked bool) {
	row := l.rows[rowId]
	replace := make([]listCell, len(runes))
	for i := 0; i < len(runes); i++ {
		replace[i] = listCell{
			r:      runes[i],
			style:  row.style,
			marked: marked,
		}
	}
	cells := append(row.cells[:col1], append(replace, row.cells[col2+1:]...)...)
	l.rows[rowId].cells = cells
}

func (l *TextList) markCells(f result, mark bool) {
	for col := f.col1; col <= f.col2; col++ {
		l.rows[f.rowId].cells[col].marked = mark
	}
}

func (l *TextList) createRow(str string) listRow {
	runes := []rune(str)
	row := listRow{
		cells: make([]listCell, len(runes)),
		style: nil,
	}
	for i, r := range runes {
		row.cells[i] = listCell{
			r:     r,
			style: nil,
		}
	}
	return row
}

func (l *TextList) getRowRunes(rowId int, ignoreCase bool) []rune {
	if rowId >= len(l.rows) {
		return nil
	}
	runes := make([]rune, len(l.rows[rowId].cells))
	for i := 0; i < len(l.rows[rowId].cells); i++ {
		r := l.rows[rowId].cells[i].r
		if ignoreCase {
			r = unicode.ToLower(r)
		}
		runes[i] = r
	}
	return runes
}

func (l *TextList) getRowString(rowId int) string {
	runes := l.getRowRunes(rowId, false)
	return string(runes)
}

func (l *TextList) lineNo(rowId int) *canvas.Text {
	ln := canvas.NewText(fmt.Sprintf(l.lineFormat, rowId+1), theme.Color(theme.ColorNameForeground))
	ln.Alignment = fyne.TextAlignCenter
	ln.TextSize = l.Theme.textSize
	ln.TextStyle = l.Theme.style
	if len(l.results) > 0 && rowId == l.results[l.currentResult].rowId {
		ln.Text = fmt.Sprintf(l.searchLineFormat, rowId+1)
		ln.Color = Name2RGBA(theme.ColorNameHyperlink)
		ln.TextStyle.Bold = true
	}
	return ln
}
