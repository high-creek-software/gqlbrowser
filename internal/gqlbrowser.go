package internal

import (
	"encoding/json"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/bento"
	"github.com/high-creek-software/gqlbrowser/internal/resources"
	"github.com/high-creek-software/gqlbrowser/internal/storage"
	"gitlab.com/high-creek-software/fieldglass"
	"golang.org/x/image/colornames"
)

type GQLBrowser struct {
	app         fyne.App
	mainWindow  fyne.Window
	setupWindow fyne.Window
	manager     storage.Manager

	stack             *fyne.Container
	settingsContainer *fyne.Container
	mainContainer     *fyne.Container

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

	setupBento *bento.Box

	schema fieldglass.Schema
}

func NewGQLBrowser() *GQLBrowser {
	// os.Setenv("FYNE_THEME", "light")
	gqlb := &GQLBrowser{app: app.NewWithID("github.com/high-creek-software/gqlbrowser")}
	gqlb.app.SetIcon(resources.IconResource)
	gqlb.mainWindow = gqlb.app.NewWindow("GQL Browser")
	gqlb.mainWindow.Resize(fyne.NewSize(1500, 800))
	gqlb.app.Lifecycle().SetOnStarted(gqlb.appStarted)
	gqlb.client = fieldglass.NewFieldGlass(true)
	gqlb.manager = storage.NewManager(gqlb.app.Storage().RootURI().Path(), gqlb.client)

	gqlb.setupBody()

	return gqlb
}

func (g *GQLBrowser) setupBody() {

	g.pathCombo = widget.NewSelect(nil, g.pathSelected)
	g.pathCombo.PlaceHolder = "Select path to inspect..."

	settingsBtn := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		g.mainContainer.Hide()
		g.settingsContainer.Show()
	})
	form := container.NewBorder(nil, nil, nil, settingsBtn, g.pathCombo)

	g.typeTabs = container.NewDocTabs()
	g.typeTabs.SetTabLocation(container.TabLocationLeading)

	g.displayContainer = container.NewGridWithRows(1)
	//g.displayContainer = container.NewHBox()
	displayScroll := container.NewScroll(g.displayContainer)
	displayScroll.Direction = container.ScrollHorizontalOnly

	// showButton := widget.NewButtonWithIcon("Show Raw Schema", theme.MenuExpandIcon(), g.showRawSchema)

	// container.NewBorder(nil, showButton, nil, nil, g.typeTabs)
	split := container.NewHSplit(g.typeTabs, displayScroll)
	split.SetOffset(0.24)

	g.mainContainer = container.NewBorder(container.NewPadded(form), nil, nil, nil, container.NewPadded(split))

	/* Setup endpoint management */
	g.setupBento = bento.NewBox()

	pathLbl := widget.NewLabel("URL")
	pathEntry := widget.NewEntry()

	g.endpointAdapter = newEndpointAdapter(g.refreshEndpoint, g.deleteEndpoint)
	g.endpointAdapter.resetAll(g.endpoints)
	g.endpointTable = widget.NewTableWithHeaders(g.endpointAdapter.count, g.endpointAdapter.createTemplate, g.endpointAdapter.updateTemplate)
	g.endpointTable.CreateHeader = g.endpointAdapter.createHeader
	g.endpointTable.UpdateHeader = g.endpointAdapter.updateHeader
	g.endpointTable.SetColumnWidth(0, 450)
	g.endpointTable.SetColumnWidth(1, 200)
	g.endpointTable.SetColumnWidth(2, 200)

	saveBtn := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
		path := pathEntry.Text
		if path == "" {
			return
		}

		g.addEndpoint(path)
	})

	backBtn := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		g.settingsContainer.Hide()
		g.mainContainer.Show()
	})

	g.settingsContainer = container.NewPadded(
		container.NewBorder(
			container.NewBorder(nil, nil, container.NewBorder(nil, nil, backBtn, nil, pathLbl), saveBtn, pathEntry),
			nil,
			nil,
			nil,
			container.NewPadded(g.endpointTable),
		),
	)

	g.settingsContainer.Hide()

	g.stack = container.New(layout.NewStackLayout(), g.mainContainer, g.settingsContainer)

	g.mainWindow.SetContent(g.stack)
}

func (g *GQLBrowser) showRawSchema() {

	schemaWindow := g.app.NewWindow("Schema")
	schemaWindow.Resize(fyne.NewSize(800, 450))

	entry := widget.NewEntry()

	schemaWindow.SetContent(entry)

	schemaWindow.Show()

	go func() {
		data, err := json.MarshalIndent(g.schema, "", "  ")

		if err != nil {
			log.Println("error generating json schema", err)
			entry.Text = "Error generating json schema: " + err.Error()
			return
		}

		entry.Text = string(data)
	}()
}

func (g *GQLBrowser) appStarted() {
	g.loadEndpoints()
	g.updatePathList()

	path := g.app.Preferences().String("selected-path")
	if path != "" {
		g.pathCombo.SetSelected(path)
	}
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
	detail := newDetailLayout(f.Name, f.Description, f.IsDeprecated, f.DeprecationReason, g.remove, g.showType)
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
	g.displayContainer.Add(detail)
}

func (g *GQLBrowser) typeSelected(id widget.ListItemID) {
	t := g.tAdapter.getItem(id)
	g.showType(t, nil)
}

func (g *GQLBrowser) interfaceSelected(id widget.ListItemID) {
	t := g.interfacesAdapter.getItem(id)
	g.showType(t, nil)
}

func (g *GQLBrowser) unionSelected(id widget.ListItemID) {
	t := g.uAdapter.getItem(id)
	g.showType(t, nil)
}

func (g *GQLBrowser) enumSelected(id widget.ListItemID) {
	t := g.eAdapter.getItem(id)
	g.showType(t, nil)
}

func (g *GQLBrowser) inputSelected(id widget.ListItemID) {
	t := g.inputAdapter.getItem(id)
	g.showType(t, nil)
}

func (g *GQLBrowser) remove(cont fyne.CanvasObject) {
	g.displayContainer.Remove(cont)
}

func (g *GQLBrowser) showType(t fieldglass.Type, f *fieldglass.Field) {

	rootType, _ := g.schema.FindType(t.RootName())
	detail := newDetailLayout(t.RootName(), rootType.Description, false, nil, g.remove, g.showType)
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
	if f != nil && len(f.Args) > 0 && argWrap == nil {
		argWrap = detail.buildArgs(f.Args)
	}
	var propWrap *fyne.Container
	if len(rootType.Fields) > 0 {
		propWrap = detail.buildProperties(rootType)
	}
	if argWrap != nil && propWrap != nil {
		split := container.NewVSplit(argWrap, propWrap)
		split.SetOffset(0.33)
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
	g.displayContainer.Add(detail)
}

func (g *GQLBrowser) refreshEndpoint(e storage.Endpoint) {
	err := g.manager.Update(e)
	if err != nil {
		if g.setupBento != nil {
			itm := bento.NewItemWithMessage(err.Error(), bento.LengthLong)
			itm.SetBackgroundColor(colornames.Red)
			g.setupBento.AddItem(itm)
		} else {
			dialog.ShowError(err, g.mainWindow)
		}
		return
	}
	g.loadEndpoints()
}

func (g *GQLBrowser) deleteEndpoint(e storage.Endpoint) {
	err := g.manager.Delete(e)
	if err != nil {
		if g.setupBento != nil {
			itm := bento.NewItemWithMessage(err.Error(), bento.LengthLong)
			itm.SetBackgroundColor(colornames.Red)
			g.setupBento.AddItem(itm)
		} else {
			dialog.ShowError(err, g.mainWindow)
		}
		return
	}
	g.loadEndpoints()
}

func (g *GQLBrowser) addEndpoint(path string) {
	endpoint, err := g.manager.Create(path)
	if err != nil {
		if g.setupBento != nil {
			itm := bento.NewItemWithMessage(err.Error(), bento.LengthLong)
			itm.SetBackgroundColor(colornames.Red)
			g.setupBento.AddItem(itm)
		} else {
			dialog.ShowError(err, g.mainWindow)
		}
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
	if g.endpointTable != nil {
		g.endpointTable.Refresh()
	}

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
	g.endpointTable = nil
	g.endpointAdapter = nil
}

func (g *GQLBrowser) Start() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in GQLBrowser.Start", r)
		}
	}()
	g.mainWindow.ShowAndRun()
}

func (g *GQLBrowser) saveSelectedEndpoint(path string) {
	g.app.Preferences().SetString("selected-path", path)
}
