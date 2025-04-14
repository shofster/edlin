package textlist

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"time"
	"unicode"
)

/*

  File:    searchview.go
  Author:  Bob Shofner

  MIT License - https://opensourcl.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
   Description: searchView manages the view and functions for search / repalce TextList text.
*/

// showSearch brings the search / replace view to the front
func (l *TextList) showSearch() {
	l.mode = modeSearch
	l.editBox.Hide()
	l.searchBox.Show()
	l.search.Enable()
	l.search.SetText("")
	l.searchText = ""
	l.replace.SetText("")
	l.disableReplace()
	l.UnselectAll()
	l.focus(l.search)
}

// searchView creates a container to find & replace text and optional control Buttons (callbacks)
func (l *TextList) searchView(buttonBar *fyne.Container) {

	l.search = widget.NewEntry()
	l.search.PlaceHolder = "<empty>"
	l.replace = widget.NewEntry()
	l.replace.PlaceHolder = "<empty>"

	sr := container.NewVBox(l.search, l.replace)

	ignoreCase := widget.NewCheckWithData("IgnoreCase", binding.BindBool(&l.ignoreCase))
	l.resultCount = widget.NewLabel(countForm(0, 0))
	l.down = widget.NewButtonWithIcon("", theme.Icon(theme.IconNameMoveDown), nil)
	l.up = widget.NewButtonWithIcon("", theme.Icon(theme.IconNameMoveUp), nil)
	search := container.NewHBox(l.resultCount, l.up, l.down)
	cancel := widget.NewButtonWithIcon("", theme.CancelIcon(), nil)

	searching := container.NewVBox(ignoreCase, search, layout.NewSpacer(), cancel)

	l.replaceAction = widget.NewButtonWithIcon("", theme.ConfirmIcon(),
		func() {
			l.replaceCells(l.results[l.currentResult].rowId,
				l.results[l.currentResult].col1,
				l.results[l.currentResult].col2,
				[]rune(l.replace.Text), true)
			l.results[l.currentResult].rowId = -1
			l.Refresh()
		})
	l.replace.ActionItem = l.replaceAction

	changed := false

	queryAction := widget.NewButtonWithIcon("", theme.SearchIcon(),
		func() {
			if changed { // no different, ignore
				//				l.searchText = l.search.Text
				changed = false
				if l.query(l.search.Text) {
					l.moveToRow(l.results[0].rowId)
					moveToOffset(l, l.results[0].rowId, 0)
				}
			} else if !l.replaceAction.Disabled() {
				l.nextResult(1)
			}
		})
	l.search.ActionItem = queryAction

	l.search.OnSubmitted = func(str string) {
		queryAction.OnTapped()
	}

	l.down.OnTapped = func() {
		l.nextResult(1)
	}
	l.up.OnTapped = func() {
		l.nextResult(-1)
	}

	l.search.OnChanged = func(s string) {
		if s != "" && s != l.searchText {
			changed = true
		}
	}

	cancel.OnTapped = func() {
		l.showEdit()
	}

	if buttonBar != nil {
		sep := canvas.NewLine(l.Theme.Color("normalColor", 0))
		l.searchBox = container.NewBorder(sep, buttonBar, nil, searching, sr)
	} else {
		l.searchBox = container.NewBorder(nil, nil, nil, searching, sr)
	}

}

func moveToOffset(l *TextList, rowId int, next int) {
	rowId += 5 * next
	rowId = min(rowId, len(l.rows)-1)
	off := float32(rowId) * l.Theme.textSize
	l.ScrollToOffset(off)
	l.ScrollTo(l.rowId + next)
	l.moveToRow(l.rowId)
}

func (l *TextList) disableReplace() {
	l.down.Disable()
	l.up.Disable()
	l.replace.Disable()
	l.replaceAction.Disable()
	l.resultCount.SetText(countForm(0, 0))
}

func (l *TextList) nextResult(next int) (found bool) {

	for _, r := range l.results {
		if r.rowId > 0 {
			found = true
			break
		}
	}
	if !found {
		l.disableReplace()
		return
	}

	i := l.currentResult + next

loop:
	if i > len(l.results)-1 {
		i = 0
	} else if i < 0 {
		i = len(l.results) - 1
	}
	if l.results[i].rowId > -1 {
		if i == l.currentResult {
			return
		}
		l.currentResult = i
		l.resultCount.SetText(countForm(i+1, len(l.results)))
		l.UnselectAll()
		l.rowId = l.results[l.currentResult].rowId
		moveToOffset(l, l.rowId, next)
		return
	}
	i += next
	goto loop
}

func (l *TextList) query(str string) bool {
	find := str
	rowId := l.rowId
	ic := l.ignoreCase
	l.results = l.findListMatch(rowId, find, ic)
	if len(l.results) < 1 {
		l.resultCount.SetText(countForm(0, 0))
		m := fmt.Sprintf("<%s> Not Found", str)
		l.toast(m, infoColor, 500*time.Millisecond)
		return false
	}
	l.clearMarkedRows(true)
	for i := 0; i < len(l.results); i++ {
		l.markCells(l.results[i], true)
	}
	if len(l.results) > 0 {
		l.down.Enable()
		l.up.Enable()
	}
	l.currentResult = 0
	l.resultCount.SetText(countForm(1, len(l.results)))
	l.replace.Enable()
	l.replaceAction.Enable()
	l.Refresh()
	return true
}

// findListMatch finds the rows that have 1 match
func (l *TextList) findListMatch(startRow int, find string, ignoreCase bool) (fs []result) {

	match := []rune(find)
	if ignoreCase {
		for g := 0; g < len(find); g++ {
			match[g] = unicode.ToLower(match[g])
		}
	}
	nRows := len(l.rows)
	for i := 0; i < nRows; i++ {
		runes := l.getRowRunes(startRow, ignoreCase)
		fs = append(fs, findRowMatch(startRow, runes, match)...)
		startRow++
		if startRow >= len(l.rows)-1 {
			startRow = 0
		}
	}
	return
}

// findListMatch finds the cells that match
func findRowMatch(rowId int, runes []rune, match []rune) (fs []result) {
	f := findCellMatch(runes, match)
	if f == nil {
		return
	}
	f.rowId = rowId
	fs = append(fs, *f)
	return
}

func countForm(n, m int) string {
	return fmt.Sprintf("%3d/%-3d", n, m)
}

func findCellMatch(runes []rune, match []rune) *result {
	rx := 0
	mx := 0
	inMatch := false
	var f result
cellLoop:
	if rx >= len(runes) {
		return nil
	}
	r := runes[rx]
	m := match[mx]
	if r == m {
		if !inMatch {
			f.col1 = rx
			inMatch = true
		}
		rx++
		mx++
		if rx-f.col1 >= len(match) {
			f.col2 = rx - 1
			return &f
		}
		goto cellLoop
	}
	rx++
	mx = 0
	inMatch = false
	goto cellLoop
}
