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
	list   *widget.List
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
	dr.updateField(f)
	if fa.list != nil {
		fa.list.SetItemHeight(id, dr.MinSize().Height)
	}
}

func (fa *fullFieldAdapter) getItem(id widget.ListItemID) fieldglass.Field {
	return fa.fields[id]
}
