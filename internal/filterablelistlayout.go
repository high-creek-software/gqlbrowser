package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"time"
)

type filterableListLayout[T any] struct {
	*fyne.Container

	adapter     filterableAdapter[T]
	filterEntry *widget.Entry
	clearBtn    *widget.Button
	list        *widget.List

	filterTimer *time.Timer
}

func newFilterableListLayout[T any](adapter filterableAdapter[T], onSelected func(id widget.ListItemID)) *filterableListLayout[T] {
	f := &filterableListLayout[T]{adapter: adapter}
	f.filterEntry = widget.NewEntry()
	f.filterEntry.SetPlaceHolder("Filter...")
	f.filterEntry.OnChanged = func(change string) {
		if f.filterTimer != nil {
			f.filterTimer.Stop()
			f.filterTimer = nil
		}
		f.filterTimer = time.NewTimer(500 * time.Millisecond)
		go func() {
			<-f.filterTimer.C
			f.adapter.filter(change)
		}()
	}
	f.clearBtn = widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() {
		f.filterEntry.SetText("")
		f.adapter.clear()
	})

	f.list = widget.NewList(f.adapter.count, f.adapter.createTemplate, f.adapter.updateTemplate)
	f.list.OnSelected = onSelected
	f.adapter.setList(f.list)

	f.Container = container.NewBorder(container.NewPadded(container.NewBorder(nil, nil, nil, container.NewPadded(f.clearBtn), container.NewPadded(f.filterEntry))),
		nil,
		nil,
		nil,
		container.NewPadded(f.list),
	)

	return f
}
