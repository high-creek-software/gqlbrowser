package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/fieldglass"
)

var _ adapter[fieldglass.Type] = (*topTypesAdapter)(nil)

type topTypesAdapter struct {
	types []fieldglass.Type
}

func (ta *topTypesAdapter) count() int {
	return len(ta.types)
}

func (ta *topTypesAdapter) createTemplate() fyne.CanvasObject {
	return widget.NewLabel("template")
}

func (ta *topTypesAdapter) updateTemplate(id widget.ListItemID, co fyne.CanvasObject) {
	t := ta.getItem(id)
	co.(*widget.Label).SetText(*t.Name)
}

func (ta *topTypesAdapter) getItem(id widget.ListItemID) fieldglass.Type {
	return ta.types[id]
}
