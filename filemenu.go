package main

import (
	"bufio"
	"edlin/textlist"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"io"
	"log"
	"path/filepath"
	"strings"
)

/*

  File:    filemenu.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: handle file menu options.
*/

func createFileMenu(w fyne.Window, theme *textlist.MyTheme) *fyne.Menu {

	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("New  ...", func() {
			var t tab
			t.title = "New"
			t.editor, t.container = fileNew(w, t.title, *theme)
			addTab(t)
		}),
		fyne.NewMenuItem("Open ...", func() {
			fileOpen(w, theme)
		}),
		fyne.NewMenuItem("Save ...", func() {
			fileSave(w, tabix)
		}),

		// fyne.NewMenuItem("Save All  ", func() { fmt.Println("Save   ") }),
		fyne.NewMenuItemSeparator(),
	)

	return fileMenu
}

func fileNew(w fyne.Window, _ string, theme textlist.MyTheme) (editor *textlist.TextList, content *fyne.Container) {
	editor, content = textlist.NewTextList(w, "", buttonBar, theme)
	editor.SetContent("hello\"こんにちは世界\"World")
	return
}

func fileOpen(w fyne.Window, theme *textlist.MyTheme) {
	nfo := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			log.Println("fileOpen:", err.Error())
			return
		} else if reader == nil {
			return
		}
		defer func(reader fyne.URIReadCloser) {
			_ = reader.Close()
		}(reader)

		var t tab
		t.path = reader.URI().Path()
		t.title = filepath.Base(t.path)
		t.editor, t.container = textlist.NewTextList(w, t.path, buttonBar, *theme)
		scanner := bufio.NewScanner(reader)
		var line string
		for scanner.Scan() {
			line = scanner.Text()
			t.editor.AddString(line)
		}

		addTab(t)
		_ = openPath.Set(filepath.Dir(t.path))

	}, w)

	path, _ := openPath.Get()
	uri := storage.NewFileURI(path)
	// override last folder
	fyne.CurrentApp().Preferences().SetString(lastFolderKey, path)
	if listable, err := storage.ListerForURI(uri); err == nil {
		nfo.SetLocation(listable)
	}
	nfo.SetView(dialog.ListView)
	nfo.Resize(w.Canvas().Size())
	nfo.Show()
}

func fileSave(w fyne.Window, tabix int) {
	nfs := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			log.Println("fileSave:", err.Error())
			return
		} else if writer == nil {
			return
		}
		defer func(writer fyne.URIWriteCloser) {
			_ = writer.Close()
		}(writer)

		// collect lines to string
		var builder strings.Builder
		iter := textlist.Iterator(tabs[tabix].editor)
		for line := range iter {
			bytes := []byte(line)
			builder.Write(bytes)
			builder.Write([]byte("\n"))
		}
		reader := strings.NewReader(builder.String())
		_, err = io.Copy(writer, reader)
		if err != nil {
			log.Println("fileSave: io.Copy:", err.Error())
			return
		}

		path := writer.URI().Path()
		_ = savePath.Set(filepath.Dir(path))
	}, w)

	path, _ := savePath.Get()
	uri := storage.NewFileURI(path)
	// override last folder
	fyne.CurrentApp().Preferences().SetString(lastFolderKey, path)
	if listable, err := storage.ListerForURI(uri); err == nil {
		nfs.SetLocation(listable)
	}
	nfs.SetView(dialog.ListView)
	nfs.Resize(w.Canvas().Size())
	nfs.Show()
}
