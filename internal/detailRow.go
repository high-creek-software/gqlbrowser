package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/gqlbrowser/internal/resources"
	"github.com/rs/xid"
	"gitlab.com/high-creek-software/fieldglass"
	"sync"
)

type detailRow struct {
	widget.BaseWidget

	uid string

	nameLbl         *widget.Label
	typeLbl         *canvas.Text
	descriptionLbl  *widget.Label
	cautionIcon     *widget.Icon
	deprecationLbl  *widget.Label
	defaultTitleLbl *widget.Label
	defaultLbl      *widget.Label

	depBox     *fyne.Container
	defaultBox *fyne.Container

	locker sync.RWMutex
}

func (d *detailRow) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (d *detailRow) CreateRenderer() fyne.WidgetRenderer {
	d.nameLbl = widget.NewLabel("")
	d.typeLbl = canvas.NewText("", theme.Color(theme.ColorNamePrimary))
	d.descriptionLbl = widget.NewLabel("")
	d.descriptionLbl.Wrapping = fyne.TextWrapWord
	d.descriptionLbl.SizeName = theme.SizeNameCaptionText
	d.cautionIcon = widget.NewIcon(resources.CautionResource)
	d.deprecationLbl = widget.NewLabel("")
	d.deprecationLbl.Wrapping = fyne.TextWrapWord
	d.deprecationLbl.SizeName = theme.SizeNameCaptionText
	d.deprecationLbl.TextStyle = fyne.TextStyle{
		Italic: true,
	}
	d.defaultTitleLbl = widget.NewLabel("Default Value:")
	d.defaultLbl = widget.NewLabel("")

	nameBox := container.NewHBox(d.nameLbl, d.typeLbl)
	d.depBox = container.New(
		layout.NewCustomPaddedLayout(0, 0, 8, 0),
		container.NewBorder(nil, nil, d.cautionIcon, nil, d.deprecationLbl),
	)
	d.defaultBox = container.NewHBox(d.defaultTitleLbl, d.defaultLbl)

	return widget.NewSimpleRenderer(
		container.NewPadded(
			container.NewVBox(
				nameBox,
				d.descriptionLbl,
				d.depBox,
				d.defaultBox,
			),
		),
	)
}

func (d *detailRow) SetData(name, typeName string, description, deprecationReason, defaultVal *string, isDeprecated bool) {

	d.nameLbl.SetText(name)
	d.typeLbl.Text = typeName

	if description != nil {
		d.descriptionLbl.SetText(*description)
		d.descriptionLbl.Show()
	} else {
		d.descriptionLbl.Hide()
	}

	if isDeprecated {
		d.depBox.Show()
		d.cautionIcon.Show()
		if deprecationReason != nil {
			d.deprecationLbl.SetText(*deprecationReason)
			d.deprecationLbl.Show()
		}
	} else {
		d.cautionIcon.Hide()
		d.deprecationLbl.Hide()
		d.depBox.Hide()
	}

	if defaultVal != nil {
		d.defaultBox.Show()
		d.defaultLbl.SetText(*defaultVal)
		d.defaultLbl.Show()
		d.defaultTitleLbl.Show()
	} else {
		d.defaultTitleLbl.Hide()
		d.defaultLbl.Hide()
		d.defaultBox.Hide()
	}
}

func (dr *detailRow) updateField(f fieldglass.Field) {
	dr.locker.Lock()
	args := ""
	if len(f.Args) > 0 {
		args = "(...)"
	}
	name := f.Name + args + ":"
	typeName := f.Type.FormatName()
	description := f.Description
	isDeprecated := f.IsDeprecated
	deprecationReason := f.DeprecationReason
	dr.locker.Unlock()

	dr.SetData(name, typeName, description, deprecationReason, nil, isDeprecated)

	//dr.Refresh()
}

func (dr *detailRow) updateInput(input fieldglass.InputValue) {
	dr.locker.Lock()
	name := input.Name + ":"
	typeName := input.Type.FormatName()
	description := input.Description
	defaultValue := input.DefaultValue
	dr.locker.Unlock()

	dr.SetData(name, typeName, description, nil, defaultValue, false)

	//dr.Refresh()
}

func newDetailRow() *detailRow {
	dr := &detailRow{uid: xid.New().String()}
	dr.ExtendBaseWidget(dr)
	dr.Refresh()
	return dr
}
