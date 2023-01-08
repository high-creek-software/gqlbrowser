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
	return newNameTypeRow("", "temp")
}

func (ta *topTypesAdapter) updateTemplate(id widget.ListItemID, co fyne.CanvasObject) {
	t := ta.getItem(id)
	row := co.(*nameTypeRow)
	row.typ = *t.Name
	row.Refresh()
}

func (ta *topTypesAdapter) getItem(id widget.ListItemID) fieldglass.Type {
	return ta.types[id]
}
