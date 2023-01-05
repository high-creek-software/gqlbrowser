package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type adapter[T any] interface {
	count() int
	createTemplate() fyne.CanvasObject
	updateTemplate(id widget.ListItemID, co fyne.CanvasObject)
	getItem(id widget.ListItemID) T
}
