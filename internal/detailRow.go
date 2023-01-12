package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/gqlbrowser/internal/resources"
	"github.com/rs/xid"
	"gitlab.com/high-creek-software/fieldglass"
	"sync"
)

type detailRow struct {
	widget.BaseWidget

	uid         string
	name        string
	typeName    string
	description *string

	isDeprecated      bool
	deprecationReason *string

	defaultValue *string

	locker sync.RWMutex
}

func (d *detailRow) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (d *detailRow) CreateRenderer() fyne.WidgetRenderer {
	nameLbl := widget.NewLabel("")
	typeText := canvas.NewText("", theme.PrimaryColor())
	descriptionLbl := widget.NewLabel("")
	descriptionLbl.Wrapping = fyne.TextWrapWord
	cautionLbl := widget.NewIcon(resources.CautionResource)
	deprecationLbl := widget.NewLabel("")
	deprecationLbl.Wrapping = fyne.TextWrapWord
	defaultTitle := widget.NewLabel("Default Value:")
	defaultLbl := widget.NewLabel("")

	return &detailRowRenderer{
		dr:             d,
		nameLbl:        nameLbl,
		typeText:       typeText,
		descriptionLbl: descriptionLbl,
		cautionLbl:     cautionLbl,
		deprecationLbl: deprecationLbl,
		defaultTitle:   defaultTitle,
		defaultLbl:     defaultLbl,
	}
}

func newDetailRow(name, typeName string, description *string) *detailRow {
	dr := &detailRow{uid: xid.New().String(), name: name, typeName: typeName, description: description}
	dr.ExtendBaseWidget(dr)
	dr.Refresh()
	return dr
}

func newDetailRowFull(name, typeName string, description, deprecationReason, defaultValue *string, isDeprecated bool) *detailRow {
	dr := &detailRow{
		uid:               xid.New().String(),
		name:              name,
		typeName:          typeName,
		description:       description,
		deprecationReason: deprecationReason,
		defaultValue:      defaultValue,
		isDeprecated:      isDeprecated,
	}
	dr.ExtendBaseWidget(dr)
	dr.Refresh()

	return dr
}

func (dr *detailRow) updateField(f fieldglass.Field) {
	dr.locker.Lock()
	args := ""
	if len(f.Args) > 0 {
		args = "(...)"
	}
	dr.name = f.Name + args + ":"
	dr.typeName = f.Type.FormatName()
	dr.description = f.Description
	dr.isDeprecated = f.IsDeprecated
	dr.deprecationReason = f.DeprecationReason
	dr.locker.Unlock()
	dr.Refresh()
}

func (dr *detailRow) updateInput(input fieldglass.InputValue) {
	dr.locker.Lock()
	dr.name = input.Name + ":"
	dr.typeName = input.Type.FormatName()
	dr.defaultValue = input.DefaultValue
	dr.locker.Unlock()
	dr.Refresh()
}

type detailRowRenderer struct {
	dr       *detailRow
	nameLbl  *widget.Label
	typeText *canvas.Text

	descriptionLbl *widget.Label

	cautionLbl     *widget.Icon
	deprecationLbl *widget.Label
	defaultTitle   *widget.Label
	defaultLbl     *widget.Label
}

func (d *detailRowRenderer) Destroy() {

}

func (d *detailRowRenderer) Layout(size fyne.Size) {
	prevSize := fyne.NewSize(0, 0)
	topLeft := fyne.NewPos(theme.Padding(), theme.Padding())
	d.nameLbl.Move(topLeft)
	nameSize := d.nameLbl.MinSize()
	typTopLeft := topLeft.Add(fyne.NewPos(nameSize.Width+theme.Padding(), 8))
	d.typeText.Move(typTopLeft)
	prevSize = nameSize

	if d.descriptionLbl.Visible() {
		topLeft = topLeft.Add(fyne.NewPos(0, prevSize.Height+theme.Padding()))
		d.descriptionLbl.Move(topLeft)
		newDescSize := fyne.NewSize(size.Width-theme.Padding(), d.descriptionLbl.MinSize().Height)
		d.descriptionLbl.Resize(newDescSize)
		prevSize = newDescSize
	}

	if d.dr.isDeprecated {
		topLeft = topLeft.Add(fyne.NewPos(0, prevSize.Height+theme.Padding()))
		d.cautionLbl.Resize(fyne.NewSize(32, 32))
		d.cautionLbl.Move(topLeft.Add(fyne.NewPos(8, 0)))
		deprecTopLeft := topLeft.Add(fyne.NewPos(38, 0))
		d.deprecationLbl.Move(deprecTopLeft)
		d.deprecationLbl.Resize(fyne.NewSize(size.Width-theme.Padding(), d.deprecationLbl.MinSize().Height))
		prevSize = d.deprecationLbl.MinSize()
	}

	if d.defaultLbl.Visible() {
		topLeft = topLeft.Add(fyne.NewPos(0, prevSize.Height+theme.Padding()))
		d.defaultTitle.Move(topLeft)
		titleSize := d.defaultTitle.MinSize()
		defTopleft := topLeft.Add(fyne.NewSize(titleSize.Width+10, 0))
		d.defaultLbl.Move(defTopleft)
	}
}

func (d *detailRowRenderer) MinSize() fyne.Size {

	maxWidth := float32(0.0)
	runningHeight := float32(0.0)

	nameSize := d.nameLbl.MinSize()
	typeSize := d.typeText.MinSize()
	maxWidth = theme.Padding() + nameSize.Width + theme.Padding() + typeSize.Width
	runningHeight += nameSize.Height + (2 * theme.Padding())

	if d.dr.description != nil {
		descSize := d.descriptionLbl.MinSize()
		maxWidth = fyne.Max(maxWidth, descSize.Width)
		runningHeight += descSize.Height + (2 * theme.Padding())
	}

	if d.dr.deprecationReason != nil {
		deprecSize := d.deprecationLbl.MinSize()
		cautionSize := d.cautionLbl.MinSize()
		maxWidth = fyne.Max(maxWidth, theme.Padding()+cautionSize.Width+theme.Padding()+deprecSize.Width)
		runningHeight += fyne.Max(deprecSize.Height, cautionSize.Height) + (2 * theme.Padding())
	}

	if d.dr.defaultValue != nil {
		defTitleSize := d.defaultTitle.MinSize()
		defValSize := d.defaultLbl.MinSize()
		maxWidth = fyne.Max(maxWidth, theme.Padding()+defTitleSize.Width+theme.Padding()+defValSize.Width)
		runningHeight += fyne.Max(defTitleSize.Height, defValSize.Height) + (2 * theme.Padding())
	}

	return fyne.NewSize(maxWidth, runningHeight)
}

func (d *detailRowRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{d.nameLbl, d.typeText, d.descriptionLbl, d.cautionLbl, d.deprecationLbl, d.defaultTitle, d.defaultLbl}
}

func (d *detailRowRenderer) Refresh() {
	d.dr.locker.RLock()
	d.nameLbl.SetText(d.dr.name)
	d.typeText.Text = d.dr.typeName
	d.typeText.Refresh()
	if d.dr.description != nil && *d.dr.description != "" {
		d.descriptionLbl.Show()
		d.descriptionLbl.SetText(*d.dr.description)
	} else {
		d.descriptionLbl.Hide()
	}

	if d.dr.isDeprecated {
		d.cautionLbl.Show()
		if d.dr.deprecationReason != nil && *d.dr.deprecationReason != "" {
			d.deprecationLbl.Show()
			d.deprecationLbl.SetText(*d.dr.deprecationReason)
		}
	} else {
		d.cautionLbl.Hide()
		d.deprecationLbl.Hide()
	}

	if d.dr.defaultValue != nil {
		d.defaultTitle.Show()
		d.defaultLbl.Show()
		d.defaultLbl.SetText(*d.dr.defaultValue)
	} else {
		d.defaultTitle.Hide()
		d.defaultLbl.Hide()
	}
	d.dr.locker.RUnlock()
}
