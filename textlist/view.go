package textlist

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

/*

  File:    view.go
  Author:  Bob Shofner

  MIT License - https://opensource.org/license/mit/

  This permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*
  Description: setup view containers.
	A Stack contains the views for EDIT, SEARCH.
    Only 1 is visible based on user interactions.

	TEXT and TAB sizes are maintained in fyne preferences.
*/

func (l *TextList) views(buttonBar *fyne.Container) {

	l.editView(buttonBar)
	l.searchView(buttonBar)

	// edit and search containers. one is visible
	l.controlBox = container.NewStack(l.searchBox, l.editBox)
}

// TextList specific fyne.Theme
var _ fyne.Theme = (*MyTheme)(nil)

type MyTheme struct {
	Width, Height int
	textSize      float32
	tabSize       int
	separatorSize float32
	doubleClick   int
	style         fyne.TextStyle
	variant       fyne.ThemeVariant
}

// NewTheme creates a Theme from fyne preferences for this app
func NewTheme(settings fyne.Settings, prefs fyne.Preferences) *MyTheme {
	t := &MyTheme{}
	t.variant = settings.ThemeVariant()
	t.Width = prefs.IntWithFallback("width", 1000)
	t.Height = prefs.IntWithFallback("height", 800)
	t.textSize = float32(prefs.FloatWithFallback("sizeText", 20))
	t.tabSize = prefs.IntWithFallback("sizeTab", 4)
	t.separatorSize = float32(prefs.FloatWithFallback("sizeSeparator", 0))
	t.doubleClick = prefs.IntWithFallback("doubleClick", 500)
	t.style.Monospace = false
	t.style.TabWidth = t.tabSize

	prefs.SetInt("width", t.Width)
	prefs.SetInt("height", t.Height)
	prefs.SetFloat("sizeText", float64(t.textSize))
	prefs.SetInt("sizeTab", t.tabSize)
	prefs.SetFloat("sizeSeparator", float64(t.separatorSize))
	prefs.SetInt("doubleClick", t.doubleClick)

	settings.SetTheme(t)
	return t
}

// Color is called to get fyne and TextList colors
func (t MyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case "normalColor":
		switch t.variant {
		case theme.VariantLight:
			return Name2RGBA(theme.ColorNameForeground)
		default:
			return Name2RGBA(theme.ColorNameForeground)
		}
	case "markedColor":
		switch t.variant {
		case theme.VariantLight:
			return color.RGBA{R: 0x7f, A: 0xff}
		default:
			return color.RGBA{R: 0xff, G: 0xf0, B: 0x1f, A: 0xff}
		}
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (t MyTheme) Font(style fyne.TextStyle) fyne.Resource {
	style.Monospace = true
	style.TabWidth = t.tabSize
	return theme.DefaultTheme().Font(style)
}
func (t MyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}
func (t MyTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		//		log.Println("Text", theme.DefaultTheme().Size(name))
		return t.textSize
	}

	return theme.DefaultTheme().Size(name)
}

func Name2RGBA(name fyne.ThemeColorName) color.RGBA {
	c := theme.Color(name)
	return Color2RGBA(c)
}
func Color2RGBA(col color.Color) color.RGBA {
	r, g, b, alpha := col.RGBA()
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(alpha)}
}
