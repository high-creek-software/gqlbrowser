package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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

	endpointAdapter *endpointAdapter
	endpointTable   *widget.Table

	schema fieldglass.Schema
}

func NewGQLBrowser() *GQLBrowser {
	gqlb := &GQLBrowser{app: app.NewWithID("github.com/high-creek-software/gqlbrowser")}
	gqlb.mainWindow = gqlb.app.NewWindow("GQL Browser")
	gqlb.mainWindow.Resize(fyne.NewSize(1500, 800))
	gqlb.app.Lifecycle().SetOnStarted(gqlb.appStarted)
	gqlb.client = fieldglass.NewFieldGlass()
	gqlb.manager = storage.NewManager(gqlb.client)

	gqlb.setupBody()

	return gqlb
}

func (g *GQLBrowser) setupBody() {
	g.pathCombo = widget.NewSelect(nil, g.pathSelected)
	g.pathCombo.PlaceHolder = "Select path to inspect..."

	settingsBtn := widget.NewButtonWithIcon("", theme.SettingsIcon(), g.settingsTouched)
	form := container.NewBorder(nil, nil, nil, settingsBtn, g.pathCombo)

	g.typeTabs = container.NewDocTabs()
	g.typeTabs.SetTabLocation(container.TabLocationLeading)

	g.displayContainer = container.NewGridWithRows(1)
	//g.displayContainer = container.NewHBox()
	displayScroll := container.NewScroll(g.displayContainer)
	displayScroll.Direction = container.ScrollHorizontalOnly

	split := container.NewHSplit(g.typeTabs, displayScroll)
	split.SetOffset(0.24)

	border := container.NewBorder(container.NewPadded(form), nil, nil, nil, container.NewPadded(split))
	g.mainWindow.SetContent(border)
}

func (g *GQLBrowser) appStarted() {
	g.loadEndpoints()
	g.updatePathList()
}

func (g *GQLBrowser) querySelected(id widget.ListItemID) {
	f := g.queryAdapter.getItem(id)
	g.showQueryMut(f)
}

func (g *GQLBrowser) mutationSelected(id widget.ListItemID) {
	f := g.mutationAdapter.getItem(id)
	g.showQueryMut(f)
}

func (g *GQLBrowser) showQueryMut(f fieldglass.Field) {
	detail := newDetailLayout(f.Name, f.Description, g.remove, g.showType)
	var argContainer *fyne.Container
	if len(f.Args) > 0 {
		argContainer = detail.buildArgs(f.Args)
	}
	var propContainer *fyne.Container
	if rootType, err := g.schema.FindType(f.Type.RootName()); err == nil {
		propContainer = detail.buildProperties(rootType)
	}
	if argContainer != nil && propContainer != nil {
		split := container.NewVSplit(argContainer, propContainer)
		split.SetOffset(0.33)
		detail.segmentWrapper.Add(split)
	} else if argContainer != nil {
		detail.segmentWrapper.Add(argContainer)
	} else if propContainer != nil {
		detail.segmentWrapper.Add(propContainer)
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
		wrap := detail.buildTypes("Implements", rootType.Interfaces)
		detail.segmentWrapper.Add(wrap)
	}
	if len(rootType.EnumValues) > 0 {
		wrap := detail.buildEnums(rootType.EnumValues)
		detail.segmentWrapper.Add(wrap)
	}
	var argWrap *fyne.Container
	if len(rootType.InputFields) > 0 {
		argWrap = detail.buildArgs(rootType.InputFields)
	}
	var propWrap *fyne.Container
	if len(rootType.Fields) > 0 {
		propWrap = detail.buildProperties(rootType)
	}
	if argWrap != nil && propWrap != nil {
		split := container.NewVSplit(argWrap, propWrap)
		detail.segmentWrapper.Add(split)
	} else if propWrap != nil {
		detail.segmentWrapper.Add(propWrap)
	} else if argWrap != nil {
		detail.segmentWrapper.Add(argWrap)
	}
	if len(rootType.PossibleTypes) > 0 {
		wrap := detail.buildTypes("Union", rootType.PossibleTypes)
		detail.segmentWrapper.Add(wrap)
	}
	g.displayContainer.Add(detail.Container)
}

func (g *GQLBrowser) settingsTouched() {
	if g.setupWindow != nil {
		return
	}
	g.setupWindow = g.app.NewWindow("Setup")
	g.setupWindow.SetOnClosed(g.setupClosed)
	g.setupWindow.Resize(fyne.NewSize(950, 600))

	pathEntry := widget.NewEntry()

	pathItem := widget.NewFormItem("Path:", pathEntry)
	inputForm := widget.NewForm(pathItem)

	g.endpointAdapter = newEndpointAdapter(g.refreshEndpoint, g.deleteEndpoint)
	g.endpointAdapter.resetAll(g.endpoints)
	g.endpointTable = widget.NewTable(g.endpointAdapter.count, g.endpointAdapter.createTemplate, g.endpointAdapter.updateTemplate)
	g.endpointTable.SetColumnWidth(0, 350)
	g.endpointTable.SetColumnWidth(1, 200)
	g.endpointTable.SetColumnWidth(2, 200)

	saveBtn := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
		path := pathEntry.Text
		if path == "" {
			return
		}

		g.addEndpoint(path)
	})

	g.setupWindow.SetContent(container.NewBorder(
		container.NewVBox(inputForm, saveBtn),
		nil,
		nil,
		nil,
		container.NewPadded(g.endpointTable),
	))

	g.setupWindow.Show()
}

func (g *GQLBrowser) refreshEndpoint(e storage.Endpoint) {
	err := g.manager.Update(e)
	if err != nil {
		dialog.ShowError(err, g.mainWindow)
		return
	}
	g.loadEndpoints()
}

func (g *GQLBrowser) deleteEndpoint(e storage.Endpoint) {
	err := g.manager.Delete(e)
	if err != nil {
		dialog.ShowError(err, g.mainWindow)
		return
	}
	g.loadEndpoints()
}

func (g *GQLBrowser) addEndpoint(path string) {
	endpoint, err := g.manager.Create(path)
	if err != nil {
		dialog.ShowError(err, g.mainWindow)
		return
	}
	g.endpointAdapter.addEndpoint(endpoint)
	g.endpoints = append(g.endpoints, endpoint)
	g.updatePathList()
}

func (g *GQLBrowser) loadEndpoints() {
	endpoints, err := g.manager.List()
	if err != nil {
		dialog.ShowError(err, g.mainWindow)
		return
	}
	g.endpoints = endpoints
	if g.endpointAdapter != nil {
		g.endpointAdapter.resetAll(g.endpoints)
	}
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
	g.endpointTable = nil
	g.endpointAdapter = nil
}

func (g *GQLBrowser) Start() {
	g.mainWindow.ShowAndRun()
}
