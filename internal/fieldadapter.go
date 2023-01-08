package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/fieldglass"
)

var _ adapter[fieldglass.Field] = (*fieldAdapter)(nil)

type fieldAdapter struct {
	fields []fieldglass.Field
}

func (n *fieldAdapter) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (fa *fieldAdapter) count() int {
	return len(fa.fields)
}

func (fa *fieldAdapter) createTemplate() fyne.CanvasObject {
	return newNameTypeRow("temp", "temp")
}

func (fa *fieldAdapter) updateTemplate(id widget.ListItemID, co fyne.CanvasObject) {
	f := fa.getItem(id)
	args := ""
	if len(f.Args) > 0 {
		args = "(...)"
	}
	row := co.(*nameTypeRow)
	row.name = f.Name + args + ":"
	row.typ = f.Type.FormatName()
	row.Refresh()
}

func (fa *fieldAdapter) getItem(id widget.ListItemID) fieldglass.Field {
	return fa.fields[id]
}
