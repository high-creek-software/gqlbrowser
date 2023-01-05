package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/fieldglass"
)

var _ adapter[fieldglass.Field] = (*fieldAdapter)(nil)

type fieldAdapter struct {
	fields []fieldglass.Field
}

func (fa *fieldAdapter) count() int {
	return len(fa.fields)
}

func (fa *fieldAdapter) createTemplate() fyne.CanvasObject {
	return widget.NewLabel("template")
}

func (fa *fieldAdapter) updateTemplate(id widget.ListItemID, co fyne.CanvasObject) {
	f := fa.getItem(id)
	args := ""
	if f.Args != nil && len(f.Args) > 0 {
		args = "(...)"
	}
	co.(*widget.Label).SetText(f.Name + args + ":" + f.Type.FormatName())
}

func (fa *fieldAdapter) getItem(id widget.ListItemID) fieldglass.Field {
	return fa.fields[id]
}
