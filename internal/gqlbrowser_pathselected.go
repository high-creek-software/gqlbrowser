package internal

import (
	"encoding/json"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/gqlbrowser/internal/storage"
	"gitlab.com/high-creek-software/fieldglass"
)

func (g *GQLBrowser) pathSelected(path string) {
	var endpoint storage.Endpoint
	for _, e := range g.endpoints {
		if e.Path == path {
			endpoint = e
			break
		}
	}
	g.displayContainer.RemoveAll()

	err := json.Unmarshal([]byte(endpoint.Payload), &g.schema)
	if err != nil {
		dialog.ShowError(err, g.mainWindow)
		return
	}

	for _, t := range g.tabs {
		g.typeTabs.Remove(t)
	}
	g.tabs = nil
	query, err := g.schema.GetQuery()
	if err == nil {
		g.queryAdapter = &fieldAdapter{fields: query.Fields}
		list := widget.NewList(g.queryAdapter.count, g.queryAdapter.createTemplate, g.queryAdapter.updateTemplate)
		list.OnSelected = g.querySelected
		queryTab := container.NewTabItem("Query", list)
		g.tabs = append(g.tabs, queryTab)
		g.typeTabs.Append(queryTab)
	}

	mutation, _ := g.schema.GetMutation()
	if g.schema.HasMutations() {
		g.mutationAdapter = &fieldAdapter{fields: mutation.Fields}
		list := widget.NewList(g.mutationAdapter.count, g.mutationAdapter.createTemplate, g.mutationAdapter.updateTemplate)
		list.OnSelected = g.mutationSelected
		mutationTab := container.NewTabItem("Mutation", list)
		g.tabs = append(g.tabs, mutationTab)
		g.typeTabs.Append(mutationTab)
	}

	g.tAdapter = &topTypesAdapter{}
	g.interfacesAdapter = &topTypesAdapter{}
	g.uAdapter = &topTypesAdapter{}
	g.eAdapter = &topTypesAdapter{}
	g.inputAdapter = &topTypesAdapter{}

	for _, t := range g.schema.Types {
		if query != nil && *t.Name == *query.Name {
			continue
		}
		if mutation != nil && *t.Name == *mutation.Name {
			continue
		}
		if !t.IsBuiltin() {
			switch t.Kind {
			case fieldglass.TypeKindObject:
				g.tAdapter.types = append(g.tAdapter.types, t)
			case fieldglass.TypeKindInterface:
				g.interfacesAdapter.types = append(g.interfacesAdapter.types, t)
			case fieldglass.TypeKindUnion:
				g.uAdapter.types = append(g.uAdapter.types, t)
			case fieldglass.TypeKindEnum:
				g.eAdapter.types = append(g.eAdapter.types, t)
			case fieldglass.TypeKindInputObject:
				g.inputAdapter.types = append(g.inputAdapter.types, t)
			}
		}
	}

	if g.tAdapter.count() > 0 {
		list := widget.NewList(g.tAdapter.count, g.tAdapter.createTemplate, g.tAdapter.updateTemplate)
		list.OnSelected = g.typeSelected
		tTab := container.NewTabItem("Types", list)
		g.tabs = append(g.tabs, tTab)
		g.typeTabs.Append(tTab)
	}
	if g.interfacesAdapter.count() > 0 {
		list := widget.NewList(g.interfacesAdapter.count, g.interfacesAdapter.createTemplate, g.interfacesAdapter.updateTemplate)
		list.OnSelected = g.interfaceSelected
		iTab := container.NewTabItem("Interfaces", list)
		g.tabs = append(g.tabs, iTab)
		g.typeTabs.Append(iTab)
	}
	if g.uAdapter.count() > 0 {
		list := widget.NewList(g.uAdapter.count, g.uAdapter.createTemplate, g.uAdapter.updateTemplate)
		list.OnSelected = g.unionSelected
		uTab := container.NewTabItem("Unions", list)
		g.tabs = append(g.tabs, uTab)
		g.typeTabs.Append(uTab)
	}
	if g.eAdapter.count() > 0 {
		list := widget.NewList(g.eAdapter.count, g.eAdapter.createTemplate, g.eAdapter.updateTemplate)
		list.OnSelected = g.enumSelected
		eTab := container.NewTabItem("Enums", list)
		g.tabs = append(g.tabs, eTab)
		g.typeTabs.Append(eTab)
	}
	if g.inputAdapter.count() > 0 {
		list := widget.NewList(g.inputAdapter.count, g.inputAdapter.createTemplate, g.inputAdapter.updateTemplate)
		list.OnSelected = g.inputSelected
		iTab := container.NewTabItem("Inputs", list)
		g.tabs = append(g.tabs, iTab)
		g.typeTabs.Append(iTab)
	}

}
