package main

import (
	"edlin/textlist"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"os"
	"time"
)

/*

  File:    edlin.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description:

fyne package --release --icon=typewriter.png --id=com.scsi.edlin

*/

type tab struct {
	editor    *textlist.TextList
	container *fyne.Container
	title     string
	path      string
}

func main() {
	a := app.NewWithID("com.scsi.edlin")
	w := a.NewWindow("EDLIN by Bob")
	w.SetIcon(resourceTypewriterPng)

	theme := textlist.NewTheme(a.Settings(), a.Preferences())
	tabItems = container.NewDocTabs()
	tabItems.SetTabLocation(container.TabLocationTop)
	tabItems.OnSelected = func(item *container.TabItem) {
		selectTabItem(item)
	}
	tabItems.OnClosed = func(item *container.TabItem) {
		removeTabItem(item)
	}

	w.Resize(fyne.NewSize(float32(theme.Width), float32(theme.Height)))

	// Set the main menu
	w.SetMainMenu(fyne.NewMainMenu(
		createFileMenu(w, theme),
		createEditMenu(),
		createHelpMenu(w)))

	setDefaultPaths(a.Preferences())

	image := canvas.NewImageFromResource(resourceTypewriterPng)
	image.FillMode = canvas.ImageFillContain
	box := container.NewStack(image)

	go func() {
		time.Sleep(time.Millisecond * 2500)
		fyne.Do(func() {
			box.Objects[0] = tabItems
			box.Refresh()
			image = nil
		})
	}()

	w.SetContent(box)
	w.ShowAndRun()

}

var tabItems *container.DocTabs
var tabix int
var tabs = make([]tab, 0)
var tabMap = make(map[string]int)
var nextSequence = 1

var buttonBar *fyne.Container
var openPath = binding.NewString()
var savePath = binding.NewString()

const lastFolderKey = "fyne:fileDialogLastFolder"

func setDefaultPaths(prefs fyne.Preferences) {
	open := prefs.String("open")
	if open == "" {
		if dir, err := os.UserHomeDir(); err == nil {
			_ = openPath.Set(dir)
		}
	} else {
		_ = openPath.Set(open)
	}
	openPath.AddListener(binding.NewDataListener(func() {
		path, _ := openPath.Get()
		if path != "" {
			prefs.SetString("open", path)
		}
	}))
	save := prefs.String("save")
	if save == "" {
		if dir, err := os.UserHomeDir(); err == nil {
			_ = savePath.Set(dir)
		}
	} else {
		_ = savePath.Set(save)
	}
	savePath.AddListener(binding.NewDataListener(func() {
		path, _ := savePath.Get()
		if path != "" {
			prefs.SetString("save", path)
		}
	}))
}

func selectTabItem(item *container.TabItem) {
	if tx, ok := tabMap[item.Text]; ok {
		tabix = tx
		//fmt.Println("select: ", tx, "tabix is", tabix)
	}
}

func removeTabItem(item *container.TabItem) {
	fmt.Println("removeTabItem: ", item.Text, len(tabs), len(tabMap))
	if tx, ok := tabMap[item.Text]; ok {
		fmt.Println("remove: ", tx, "tabix is", tabix, "tabs", len(tabs))
		if len(tabs) < 2 {
			tabix = 0
			tabs = nil
		} else if tx == len(tabs)-1 {
			tabsNew := tabs[:tx]
			tabs = tabsNew
			tabix = len(tabs) - 1
		} else {
			tabsNew := append(tabs[:tx], tabs[tx+1:]...)
			tabs = tabsNew
			tabix = tx
			if tabix > 0 {
				tabix--
			}
		}
	}
	delete(tabMap, item.Text)
}

func addTab(t tab) {
	// insure unique title
	if _, ok := tabMap[t.title]; ok {
		t.title = fmt.Sprintf("%s(%d)", t.title, nextSequence)
		nextSequence++
	}
	tabs = append(tabs, t)
	tabItem := container.NewTabItem(t.title, t.container)
	tabItems.Append(tabItem)
	tabItems.Select(tabItem)
	tabix = len(tabs) - 1
	tabMap[t.title] = tabix
}
