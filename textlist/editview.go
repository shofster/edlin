package textlist

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

/*

  File:    editview.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
   Description: editView manages the view and functions for modifying TextList rows.
*/

// showEdit brings the edit view to the front
func (l *TextList) showEdit() {
	l.mode = modeEdit
	l.searchBox.Hide()
	l.edit.SetText("")
	l.editBox.Show()
	l.UnselectAll()
	l.startMark = -1
	l.endMark = -1
	l.results = nil
	l.Refresh()
	l.focus(l)
}

// editView creates a container with an editable Entry and optional control Buttons (callbacks)
func (l *TextList) editView(buttonBar *fyne.Container) {

	l.edit = widget.NewMultiLineEntry()
	l.edit.SetMinRowsVisible(4)
	l.edit.SetPlaceHolder("<empty>")
	confirm := widget.NewButtonWithIcon("", theme.ConfirmIcon(), nil)
	cancel := widget.NewButtonWithIcon("", theme.CancelIcon(), nil)

	confirm.Disable()
	cancel.Disable()

	var process = func(str string) {
		if str != l.editText {
			l.replaceRow(l.rowId, str)
		}
		confirm.Disable()
		cancel.Disable()
		l.showEdit()
	}

	l.edit.OnChanged = func(s string) {
		if s != l.editText {
			confirm.Enable()
			cancel.Enable()
		}
	}
	confirm.OnTapped = func() {
		process(l.edit.Text)
	}
	l.edit.OnSubmitted = func(str string) {
		process(str)
	}
	cancel.OnTapped = func() {
		confirm.Disable()
		cancel.Disable()
		l.showEdit()
	}

	if buttonBar != nil {
		sep := canvas.NewLine(l.Theme.Color("normalColor", 0))
		l.editBox = container.NewBorder(sep, buttonBar, nil,
			container.NewVBox(confirm, layout.NewSpacer(), cancel), l.edit)
	} else {
		l.editBox = container.NewBorder(nil, nil, nil,
			container.NewVBox(confirm, layout.NewSpacer(), cancel), l.edit)
	}

}
