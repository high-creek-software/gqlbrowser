package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/fieldglass"
)

var _ adapter[fieldglass.Field] = (*fullFieldAdapter)(nil)

type fullFieldAdapter struct {
	fields []fieldglass.Field
}

func (f *fullFieldAdapter) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (fa *fullFieldAdapter) count() int {
	return len(fa.fields)
}

func (fa *fullFieldAdapter) createTemplate() fyne.CanvasObject {
	return newDetailRow("template", "template", nil)
}

func (fa *fullFieldAdapter) updateTemplate(id widget.ListItemID, co fyne.CanvasObject) {
	f := fa.getItem(id)
	dr := co.(*detailRow)
	args := ""
	if len(f.Args) > 0 {
		args = "(...)"
	}
	dr.name = f.Name + args + ":"
	dr.typeName = f.Type.FormatName()
	dr.description = f.Description
	dr.isDeprecated = f.IsDeprecated
	dr.deprecationReason = f.DeprecationReason
	dr.Refresh()
}

func (fa *fullFieldAdapter) getItem(id widget.ListItemID) fieldglass.Field {
	return fa.fields[id]
}
