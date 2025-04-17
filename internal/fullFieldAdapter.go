package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"gitlab.com/high-creek-software/fieldglass"
	"sync"
)

var _ adapter[fieldglass.Field] = (*fullFieldAdapter)(nil)

type fullFieldAdapter struct {
	fields []fieldglass.Field
	list   *widget.List

	filtered []fieldglass.Field
	locker   sync.Mutex
}

func (f *fullFieldAdapter) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (fa *fullFieldAdapter) count() int {
	if len(fa.filtered) > 0 {
		return len(fa.filtered)
	}
	return len(fa.fields)
}

func (fa *fullFieldAdapter) createTemplate() fyne.CanvasObject {
	return newDetailRow()
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
	if len(fa.filtered) > 0 {
		return fa.filtered[id]
	}
	return fa.fields[id]
}

func (fa *fullFieldAdapter) setList(list *widget.List) {
	fa.list = list
}

func (fa *fullFieldAdapter) filter(input string) {
	fa.locker.Lock()
	defer fa.locker.Unlock()

	fa.filtered = nil

	for _, f := range fa.fields {
		if fuzzy.MatchNormalizedFold(input, f.Name) || fuzzy.MatchNormalizedFold(input, f.Type.RootName()) {
			fa.filtered = append(fa.filtered, f)
		}
	}
	fa.list.Refresh()
}

func (fa *fullFieldAdapter) clear() {

}
