package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

/*

  File:    editmenu.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:
*/

var typeShortcut = func(keyName fyne.KeyName) {
	cs := desktop.CustomShortcut{keyName, fyne.KeyModifierControl}
	tabs[tabix].editor.TypedShortcut(&cs)
}

func createEditMenu() *fyne.Menu {
	menu := fyne.NewMenu("Edit",
		fyne.NewMenuItem("Begin ^M", func() {
			typeShortcut("M")
		}),
		fyne.NewMenuItem("End   ^E", func() {
			typeShortcut("E")
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Cut   ^X", func() {
			typeShortcut("X")
		}),
		fyne.NewMenuItem("Copy  ^C", func() {
			typeShortcut("C")
		}),
		fyne.NewMenuItem("Paste ^V", func() {
			typeShortcut("V")
		}),
	)
	return menu
}
