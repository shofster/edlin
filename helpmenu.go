package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

/*

  File:    helpmenu.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/

func createHelpMenu(w fyne.Window) *fyne.Menu {

	fileMenuItem := fyne.NewMenuItem("FileMenu", func() {
		d := dialog.NewInformation("File", helpFile, w)
		d.Resize(w.Canvas().Size())
		d.Show()
	})
	editMenuItem := fyne.NewMenuItem("EditMenu", func() {
		d := dialog.NewInformation("Edit", helpEdit, w)
		d.Resize(w.Canvas().Size())
		d.Show()
	})
	shortcutItem := fyne.NewMenuItem("Shortcuts", func() {
		d := dialog.NewInformation("Shortcuts", helpShortcut, w)
		d.Resize(w.Canvas().Size())
		d.Show()
	})
	searchItem := fyne.NewMenuItem("Search", func() {
		d := dialog.NewInformation("Search / Replace", helpSearch, w)
		d.Resize(w.Canvas().Size())
		d.Show()
	})
	//	subMenu := fyne.NewMenu("HELP", fileMenuItem, editMenuItem, shortcutItem)
	helpMenu := fyne.NewMenuItem("Help", nil)
	helpMenu.ChildMenu = fyne.NewMenu("HELP", fileMenuItem, editMenuItem, shortcutItem, searchItem)

	menu := fyne.NewMenu("Help", helpMenu,
		fyne.NewMenuItem("About", func() {
			dialog.ShowInformation("EDLIN by Bob", "Version 0.9.1 - Apr 2025", w)
		}),
	)

	return menu
}

var helpFile = `EDLIN Help:

FileMenu:

Open a new empty tab:       New  ...
Open tab from a file:       Open ...
Save tab contets to a file: Save ...
`

var helpEdit = `EDLIN Help:

EditMenu: (Available only in Edit Mode.)

Mark Lines:
  Begin Line:   Begin ^M  or Keyboard Ctrl + M keys.
    End LIne:   End   ^E  or Keyboard Ctrl + E keys.

Remove Lines:   Cut   ^X  or keyboard Ctrl + X keys.
     If marked lines, remove all, else only the selected line.
Copy Lines:     Copy  ^C  or keyboard Ctrl + C keys.
     If marked lines, copy all, else only the selected line.

Paste Lines:    Paste ^V or keyboard Ctrl + V keys.

     (Deleted or Copied line are copied to the Clipboard.)
     (Lines from the Clipboard are inserted before the current line.)

Selecting (with mouse) a line allows editing of that line.
The replacement entry may be 1 or many new lines.
`

var helpShortcut = `EDLIN Help:

Shortcut Keys:

Ctrl + F or Ctrl + R enters Search/Replace Mode.
     ('x' in Search/Replace mode returns to Edit Mode.)

Ctrl + Home: Position the list at line 0.
Ctrl + End:  Position the list at the last line.
Ctrl + Down: Position the list 10 rows down. (PageDown)
Ctrl + Up:   Position the list 10 rows up.   (PageUp)
     Down:   Position the list 1 row down.   (Down Arrow)
       Up:   Position the list 1 row up.     (Up Arrow)
`

var helpSearch = `EDLIN Help:

Search / Replace:

Enter text to find in the Search Box.
Optionally choose Ignore Case.
Click on the search ICON.

All the matches are found, and the first is highlighted.
Press the UP or DOWN arrows to advance to a previous or next match.

If replacing : enter new text and press the confirm ICON. 
Press the UP or DOWN arrows to advance to a previous or next match.

`
