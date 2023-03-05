package internal

import (
	"encoding/json"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"gitlab.com/high-creek-software/fieldglass"
)

func (g *GQLBrowser) pathSelected(path string) {
	/*var endpoint storage.Endpoint
	for _, e := range g.endpoints {
		if e.Path == path {
			endpoint = e
			break
		}
	}*/
	endpoint := g.endpoints[g.pathCombo.SelectedIndex()]
	g.saveSelectedEndpoint(path)
	g.displayContainer.RemoveAll()

	g.schema = fieldglass.Schema{}
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
		queryTab := container.NewTabItem("Query", newFilterableListLayout[fieldglass.Field](g.queryAdapter, g.querySelected).Container)
		g.tabs = append(g.tabs, queryTab)
		g.typeTabs.Append(queryTab)
	}

	mutation, _ := g.schema.GetMutation()
	if g.schema.HasMutations() {
		g.mutationAdapter = &fieldAdapter{fields: mutation.Fields}
		mutationTab := container.NewTabItem("Mutation", newFilterableListLayout[fieldglass.Field](g.mutationAdapter, g.mutationSelected).Container)
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
		tTab := container.NewTabItem("Types", newFilterableListLayout[fieldglass.Type](g.tAdapter, g.typeSelected).Container)
		g.tabs = append(g.tabs, tTab)
		g.typeTabs.Append(tTab)
	}
	if g.interfacesAdapter.count() > 0 {
		iTab := container.NewTabItem("Interfaces", newFilterableListLayout[fieldglass.Type](g.interfacesAdapter, g.interfaceSelected).Container)
		g.tabs = append(g.tabs, iTab)
		g.typeTabs.Append(iTab)
	}
	if g.uAdapter.count() > 0 {
		uTab := container.NewTabItem("Unions", newFilterableListLayout[fieldglass.Type](g.uAdapter, g.unionSelected).Container)
		g.tabs = append(g.tabs, uTab)
		g.typeTabs.Append(uTab)
	}
	if g.eAdapter.count() > 0 {
		eTab := container.NewTabItem("Enums", newFilterableListLayout[fieldglass.Type](g.eAdapter, g.enumSelected))
		g.tabs = append(g.tabs, eTab)
		g.typeTabs.Append(eTab)
	}
	if g.inputAdapter.count() > 0 {
		iTab := container.NewTabItem("Inputs", newFilterableListLayout[fieldglass.Type](g.inputAdapter, g.inputSelected).Container)
		g.tabs = append(g.tabs, iTab)
		g.typeTabs.Append(iTab)
	}

}
