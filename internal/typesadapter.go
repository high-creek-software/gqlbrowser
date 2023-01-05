package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/fieldglass"
)

var _ adapter[fieldglass.Type] = (*typesAdapter)(nil)

type typesAdapter struct {
	types []fieldglass.Type
}

func (ta *typesAdapter) count() int {
	return len(ta.types)
}

func (ta *typesAdapter) createTemplate() fyne.CanvasObject {
	return widget.NewLabel("template")
}

func (ta *typesAdapter) updateTemplate(id widget.ListItemID, co fyne.CanvasObject) {
	t := ta.getItem(id)
	co.(*widget.Label).SetText(*t.Name)
}

func (ta *typesAdapter) getItem(id widget.ListItemID) fieldglass.Type {
	return ta.types[id]
}
