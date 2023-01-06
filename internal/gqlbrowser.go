package internal

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/gqlbrowser/internal/storage"
	"gitlab.com/high-creek-software/fieldglass"
)

type GQLBrowser struct {
	app         fyne.App
	mainWindow  fyne.Window
	setupWindow fyne.Window
	manager     storage.Manager

	endpoints []storage.Endpoint

	client fieldglass.FieldGlass

	typeTabs          *container.DocTabs
	tabs              []*container.TabItem
	queryAdapter      *fieldAdapter
	mutationAdapter   *fieldAdapter
	tAdapter          *topTypesAdapter
	interfacesAdapter *topTypesAdapter
	uAdapter          *topTypesAdapter
	eAdapter          *topTypesAdapter
	inputAdapter      *topTypesAdapter

	pathCombo        *widget.Select
	displayContainer *fyne.Container

	schema fieldglass.Schema
}

func NewGQLBrowser() *GQLBrowser {
	gqlb := &GQLBrowser{manager: storage.NewManager(), app: app.NewWithID("github.com/high-creek-software/gqlbrowser")}
	gqlb.mainWindow = gqlb.app.NewWindow("GQL Browser")
	gqlb.mainWindow.Resize(fyne.NewSize(1200, 700))
	gqlb.app.Lifecycle().SetOnStarted(gqlb.appStarted)
	gqlb.client = fieldglass.NewFieldGlass()

	gqlb.setupBody()

	return gqlb
}

func (g *GQLBrowser) setupBody() {
	g.pathCombo = widget.NewSelect(nil, g.pathSelected)
	g.pathCombo.PlaceHolder = "Select path to inspect..."

	settingsBtn := widget.NewButtonWithIcon("", theme.SettingsIcon(), g.settingsTouched)

	form := container.New(layout.NewFormLayout(), settingsBtn, g.pathCombo)

	g.typeTabs = container.NewDocTabs()
	g.typeTabs.SetTabLocation(container.TabLocationLeading)
	g.displayContainer = container.NewGridWithRows(1)

	displayScroll := container.NewScroll(g.displayContainer)
	displayScroll.Direction = container.ScrollHorizontalOnly
	split := container.NewHSplit(g.typeTabs, displayScroll)
	split.SetOffset(0.33)

	border := container.NewBorder(form, nil, nil, nil, split)
	g.mainWindow.SetContent(border)
}

func (g *GQLBrowser) appStarted() {
	endpoints, err := g.manager.List()
	if err != nil {
		dialog.ShowError(err, g.mainWindow)
		return
	}
	g.endpoints = endpoints
	g.updatePathList()
}

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

	// TODO: Setup schema
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

func (g *GQLBrowser) querySelected(id widget.ListItemID) {
	f := g.queryAdapter.getItem(id)
	detail := newDetailLayout(f.Name, f.Description, g.remove, g.showType)
	if len(f.Args) > 0 {
		detail.addArgs(f.Args)
	}
	if rootType, err := g.schema.FindType(f.Type.RootName()); err == nil {
		detail.addProperties(rootType)
	}
	g.displayContainer.Add(detail.Container)
}

func (g *GQLBrowser) mutationSelected(id widget.ListItemID) {
	f := g.mutationAdapter.getItem(id)
	detail := newDetailLayout(f.Name, f.Description, g.remove, g.showType)
	if len(f.Args) > 0 {
		detail.addArgs(f.Args)
	}
	g.displayContainer.Add(detail.Container)
}

func (g *GQLBrowser) typeSelected(id widget.ListItemID) {
	t := g.tAdapter.getItem(id)
	g.showType(t)
}

func (g *GQLBrowser) interfaceSelected(id widget.ListItemID) {
	t := g.interfacesAdapter.getItem(id)
	g.showType(t)
}

func (g *GQLBrowser) unionSelected(id widget.ListItemID) {
	t := g.uAdapter.getItem(id)
	g.showType(t)
}

func (g *GQLBrowser) enumSelected(id widget.ListItemID) {
	t := g.eAdapter.getItem(id)
	g.showType(t)
}

func (g *GQLBrowser) inputSelected(id widget.ListItemID) {
	t := g.inputAdapter.getItem(id)
	g.showType(t)
}

func (g *GQLBrowser) remove(cont *fyne.Container) {
	g.displayContainer.Remove(cont)
}

func (g *GQLBrowser) showType(t fieldglass.Type) {

	rootType, _ := g.schema.FindType(t.RootName())
	detail := newDetailLayout(t.RootName(), nil, g.remove, g.showType)
	if len(rootType.Interfaces) > 0 {
		detail.addTypes("Implements", rootType.Interfaces)
	}
	if len(rootType.EnumValues) > 0 {
		detail.addEnums(rootType.EnumValues)
	}
	if len(rootType.InputFields) > 0 {
		detail.addArgs(rootType.InputFields)
	}
	if len(rootType.Fields) > 0 {
		detail.addProperties(rootType)
	}
	if len(rootType.PossibleTypes) > 0 {
		detail.addTypes("Union", rootType.PossibleTypes)
	}
	g.displayContainer.Add(detail.Container)
}

func (g *GQLBrowser) settingsTouched() {
	if g.setupWindow != nil {
		return
	}
	g.setupWindow = g.app.NewWindow("Setup")
	g.setupWindow.SetOnClosed(g.setupClosed)
	g.setupWindow.Resize(fyne.NewSize(600, 300))
	g.setupWindow.CenterOnScreen()

	pathEntry := widget.NewEntry()

	pathItem := widget.NewFormItem("Path:", pathEntry)
	inputForm := widget.NewForm(pathItem)

	saveBtn := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
		path := pathEntry.Text
		if path == "" {
			return
		}
		g.setupWindow.Close()
		go g.savePath(path)
	})

	g.setupWindow.SetContent(container.NewVBox(
		layout.NewSpacer(),
		container.NewVBox(inputForm, saveBtn),
		layout.NewSpacer(),
	))

	g.setupWindow.Show()
}

func (g *GQLBrowser) savePath(path string) {
	schema, err := g.client.Load(path)
	if err != nil {
		dialog.ShowError(err, g.mainWindow)
		return
	}

	payload, err := json.Marshal(schema)
	if err != nil {
		dialog.ShowError(err, g.mainWindow)
		return
	}

	endpoint, err := g.manager.Store(path, string(payload))
	if err != nil {
		dialog.ShowError(err, g.mainWindow)
		return
	}

	g.endpoints = append(g.endpoints, endpoint)
	g.updatePathList()
}

func (g *GQLBrowser) updatePathList() {
	var paths []string
	for _, e := range g.endpoints {
		paths = append(paths, e.Path)
	}

	g.pathCombo.Options = paths
	g.pathCombo.Refresh()
}

func (g *GQLBrowser) setupClosed() {
	g.setupWindow = nil
}

func (g *GQLBrowser) Start() {
	g.mainWindow.ShowAndRun()
}
