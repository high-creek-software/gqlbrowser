package internal

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/gqlbrowser/internal/storage"
)

type endpointAdapter struct {
	endpoints []storage.Endpoint
	refresh   func(storage.Endpoint)
	delete    func(storage.Endpoint)
}

func newEndpointAdapter(refresh func(endpoint storage.Endpoint), delete func(endpoint storage.Endpoint)) *endpointAdapter {
	return &endpointAdapter{refresh: refresh, delete: delete}
}

func (e *endpointAdapter) resetAll(endpoints []storage.Endpoint) {
	e.endpoints = endpoints
}

func (e *endpointAdapter) addEndpoint(ep storage.Endpoint) {
	e.endpoints = append(e.endpoints, ep)
}

func (e *endpointAdapter) count() (int, int) {
	return len(e.endpoints), 5
}

func (e *endpointAdapter) createTemplate() fyne.CanvasObject {
	lbl := widget.NewLabel("template")
	btn := widget.NewButton("template", nil)

	return container.NewMax(lbl, container.NewPadded(btn))
}

func (e *endpointAdapter) createHeader() fyne.CanvasObject {
	return widget.NewLabel("-- --")
}

func (e *endpointAdapter) updateHeader(id widget.TableCellID, co fyne.CanvasObject) {
	lbl := co.(*widget.Label)
	if id.Row == -1 {
		switch id.Col {
		case 0:
			lbl.SetText("Path")
		case 1:
			lbl.SetText("Created")
		case 2:
			lbl.SetText("Updated")
		case 3:
			lbl.SetText("Refresh")
		case 4:
			lbl.SetText("Delete")
		}
	} else if id.Col == -1 {
		lbl.SetText(fmt.Sprintf("%d", id.Row+1))
	}
}

func (e *endpointAdapter) updateTemplate(i widget.TableCellID, co fyne.CanvasObject) {
	lbl := co.(*fyne.Container).Objects[0].(*widget.Label)
	btn := co.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*widget.Button)
	lbl.Hide()
	btn.Hide()

	endpoint := e.item(i.Row)
	switch i.Col {
	case 0:
		lbl.Show()
		lbl.SetText(endpoint.Path)
		lbl.Wrapping = fyne.TextWrapWord
	case 1:
		lbl.Show()
		lbl.SetText(endpoint.CreatedAt.Format(time.DateTime))
	case 2:
		lbl.Show()
		txt := "Not refreshed"
		if endpoint.UpdatedAt != nil {
			txt = endpoint.UpdatedAt.Format(time.DateTime)
		}
		lbl.SetText(txt)
	case 3:
		btn.Show()
		btn.SetText("Refresh")
		btn.OnTapped = func() {
			e.refresh(endpoint)
		}
	case 4:
		btn.Show()
		btn.SetText("Delete")
		btn.OnTapped = func() {
			e.delete(endpoint)
		}
	}
}

func (e *endpointAdapter) item(idx int) storage.Endpoint {
	return e.endpoints[idx]
}
