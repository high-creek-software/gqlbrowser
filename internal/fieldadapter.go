package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"gitlab.com/high-creek-software/fieldglass"
	"sync"
)

var _ adapter[fieldglass.Field] = (*fieldAdapter)(nil)

type fieldAdapter struct {
	fields []fieldglass.Field
	list   *widget.List

	filtered []fieldglass.Field
	locker   sync.Mutex
}

func (n *fieldAdapter) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (fa *fieldAdapter) count() int {
	if len(fa.filtered) > 0 {
		return len(fa.filtered)
	}
	return len(fa.fields)
}

func (fa *fieldAdapter) createTemplate() fyne.CanvasObject {
	return newNameTypeRow()
}

func (fa *fieldAdapter) updateTemplate(id widget.ListItemID, co fyne.CanvasObject) {
	f := fa.getItem(id)
	args := ""
	if len(f.Args) > 0 {
		args = "(...)"
	}
	row := co.(*nameTypeRow)
	depReason := ""
	if f.DeprecationReason != nil {
		depReason = *f.DeprecationReason
	}
	description := ""
	if f.Description != nil {
		description = *f.Description
	}
	row.SetData(f.Name+args+":", f.Type.FormatName(), depReason, description, f.IsDeprecated)
	if fa.list != nil {
		fa.list.SetItemHeight(id, row.MinSize().Height)
	}
}

func (fa *fieldAdapter) getItem(id widget.ListItemID) fieldglass.Field {
	if len(fa.filtered) > 0 {
		return fa.filtered[id]
	}
	return fa.fields[id]
}

func (fa *fieldAdapter) setList(list *widget.List) {
	fa.list = list
}

func (fa *fieldAdapter) filter(input string) {
	fa.locker.Lock()
	defer fa.locker.Unlock()

	fa.filtered = nil
	if input == "" || len(input) <= 3 {
		return
	}

	for _, f := range fa.fields {
		if fuzzy.MatchNormalizedFold(input, f.Name) || fuzzy.MatchNormalizedFold(input, f.Type.RootName()) {
			fa.filtered = append(fa.filtered, f)
		}
	}
	fa.list.Refresh()
}

func (fa *fieldAdapter) clear() {

}
