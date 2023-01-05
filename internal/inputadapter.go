package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/fieldglass"
)

type inputAdapter struct {
	inputs []fieldglass.InputValue
}

func (ia *inputAdapter) count() int {
	return len(ia.inputs)
}

func (ia *inputAdapter) createTemplate() fyne.CanvasObject {
	return widget.NewLabel("template")
}

func (ia *inputAdapter) updateTemplate(id widget.ListItemID, co fyne.CanvasObject) {
	iv := ia.getItem(id)
	co.(*widget.Label).SetText(iv.Name + ":" + iv.Type.FormatName())
}

func (ia *inputAdapter) getItem(id widget.ListItemID) fieldglass.InputValue {
	return ia.inputs[id]
}
