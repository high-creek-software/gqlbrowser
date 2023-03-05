package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"gitlab.com/high-creek-software/fieldglass"
	"sync"
)

type inputAdapter struct {
	inputs []fieldglass.InputValue
	list   *widget.List

	filtered []fieldglass.InputValue
	locker   sync.Mutex
}

func (ia *inputAdapter) count() int {
	if len(ia.filtered) > 0 {
		return len(ia.filtered)
	}
	return len(ia.inputs)
}

func (ia *inputAdapter) createTemplate() fyne.CanvasObject {
	return newDetailRow("temp", "temp", nil)
}

func (ia *inputAdapter) updateTemplate(id widget.ListItemID, co fyne.CanvasObject) {
	iv := ia.getItem(id)
	dr := co.(*detailRow)
	dr.updateInput(iv)
	if ia.list != nil {
		ia.list.SetItemHeight(id, dr.MinSize().Height)
	}
}

func (ia *inputAdapter) getItem(id widget.ListItemID) fieldglass.InputValue {
	if len(ia.filtered) > 0 {
		return ia.filtered[id]
	}
	return ia.inputs[id]
}

func (ia *inputAdapter) setList(list *widget.List) {
	ia.list = list
}

func (ia *inputAdapter) filter(input string) {
	ia.locker.Lock()
	defer ia.locker.Unlock()

	ia.filtered = nil

	for _, i := range ia.inputs {
		if fuzzy.MatchNormalizedFold(input, i.Name) || fuzzy.MatchNormalizedFold(input, i.Type.RootName()) {
			ia.filtered = append(ia.filtered, i)
		}
	}
	ia.list.Refresh()
}

func (ia *inputAdapter) clear() {

}
