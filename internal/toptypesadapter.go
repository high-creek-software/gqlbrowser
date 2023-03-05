package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"gitlab.com/high-creek-software/fieldglass"
	"sync"
)

var _ adapter[fieldglass.Type] = (*topTypesAdapter)(nil)

type topTypesAdapter struct {
	list  *widget.List
	types []fieldglass.Type

	filtered []fieldglass.Type
	locker   sync.Mutex
}

func (ta *topTypesAdapter) count() int {
	if len(ta.filtered) > 0 {
		return len(ta.filtered)
	}
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
	if len(ta.filtered) > 0 {
		return ta.filtered[id]
	}
	return ta.types[id]
}

func (ta *topTypesAdapter) setList(list *widget.List) {
	ta.list = list
}

func (ta *topTypesAdapter) filter(input string) {
	ta.locker.Lock()
	defer ta.locker.Unlock()

	ta.filtered = nil
	if input == "" || len(input) <= 3 {
		return
	}

	for _, t := range ta.types {
		if fuzzy.MatchNormalizedFold(input, t.RootName()) {
			ta.filtered = append(ta.filtered, t)
		}
	}
	ta.list.Refresh()
}

func (ta *topTypesAdapter) clear() {
	//ta.locker.Lock()
	//defer ta.locker.Lock()
	//
	//ta.filtered = nil
	//ta.list.Refresh()
}
