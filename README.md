# EDLIN
GO/fyne text editor.

EDLIN is a basic text editor written in the GO language and using
the fyno.io graphical toolkit.

Features include file open and save, line editing, replacement, and deletion.
Case sensititive global search and replace is provided.
Strings encoded in UNICODE/UTF-8 are supported.

This GO version specifies 1.24, but only 1.23 is required
[for range Iterator functionality].
go.mod specifies fyne.io/fyne/v2 v2.6.0

EDLIN implements a container.DocTabs with each tab an extended widget.List to provide line and text editing.
The standard dialog.FileDialog is used for file [open and save] operations.

Window, text, and tab sizes are controllable via fyne.Preferences.
EDLIN is available at https://github.com/shofster/edlin.
The fyne toolkit is found at https://github.com/fyne-io/fyne.

The package "textlist" is the self contained editor and may be imbedded into other fyne.Containers.