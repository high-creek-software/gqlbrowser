package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/xid"
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
}

func (d *detailRow) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (d *detailRow) CreateRenderer() fyne.WidgetRenderer {
	nameLbl := widget.NewLabel("")
	typeText := canvas.NewText("", theme.PrimaryColor())
	descriptionLbl := widget.NewLabel("")
	descriptionLbl.Wrapping = fyne.TextWrapWord
	cautionLbl := widget.NewLabel("⚠️")
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

type detailRowRenderer struct {
	dr       *detailRow
	nameLbl  *widget.Label
	typeText *canvas.Text

	descriptionLbl *widget.Label

	cautionLbl     *widget.Label
	deprecationLbl *widget.Label
	defaultTitle   *widget.Label
	defaultLbl     *widget.Label
}

func (d *detailRowRenderer) Destroy() {

}

func (d *detailRowRenderer) Layout(size fyne.Size) {
	//log.Println("Layout:", size.Width, "X", size.Height, d.dr.uid, time.Now())
	prevSize := fyne.NewSize(0, 0)
	topLeft := fyne.NewPos(theme.Padding(), theme.Padding()+10)
	d.nameLbl.Move(topLeft)
	nameSize := d.nameSize()
	typTopLeft := topLeft.Add(fyne.NewPos(nameSize.Width+14, 8))
	d.typeText.Move(typTopLeft)
	prevSize = nameSize

	if d.descriptionLbl.Visible() {
		topLeft = topLeft.Add(fyne.NewPos(0, nameSize.Height+theme.Padding()))
		descSize := d.descriptionSize()
		d.descriptionLbl.Move(topLeft)
		d.descriptionLbl.Resize(fyne.NewSize(size.Width-2*theme.Padding(), size.Height))
		prevSize = descSize
	}

	if d.deprecationLbl.Visible() {
		topLeft = topLeft.Add(fyne.NewPos(0, prevSize.Height+15+theme.Padding()))
		d.cautionLbl.Move(topLeft)
		deprecTopLeft := topLeft.Add(fyne.NewSize(30, 0))
		d.deprecationLbl.Move(deprecTopLeft)
		d.deprecationLbl.Resize(fyne.NewSize(size.Width-2*theme.Padding(), size.Height))
		prevSize = d.deprecationSize()
	}

	if d.dr.defaultValue != nil {
		topLeft = topLeft.Add(fyne.NewPos(0, prevSize.Height+15+theme.Padding()))
		d.defaultTitle.Move(topLeft)
		titleSize := fyne.MeasureText(d.defaultTitle.Text, theme.TextSize(), d.defaultTitle.TextStyle)
		defTopleft := topLeft.Add(fyne.NewSize(titleSize.Width+10, 0))
		d.defaultLbl.Move(defTopleft)
	}
}

func (d *detailRowRenderer) MinSize() fyne.Size {

	nameSize := d.nameSize()
	typeSize := d.typeSize()
	descSize := d.descriptionSize()
	deprecSize := d.deprecationSize()
	defValSize := d.defaultValueSize()

	height := nameSize.Height + descSize.Height + deprecSize.Height + defValSize.Height + 30 + 4*theme.Padding()
	width := fyne.Max(nameSize.Width+typeSize.Width, descSize.Width) + 2*theme.Padding()
	width = fyne.Max(width, descSize.Width)
	return fyne.NewSize(width, height)
}

func (d *detailRowRenderer) nameSize() fyne.Size {
	if d.dr.name == "" {
		return fyne.MeasureText("Temp", theme.TextSize(), d.descriptionLbl.TextStyle)
	}
	return fyne.MeasureText(d.dr.name, theme.TextSize(), d.nameLbl.TextStyle)
}

func (d *detailRowRenderer) typeSize() fyne.Size {
	if d.dr.typeName == "" {
		return fyne.MeasureText("Temp", theme.TextSize(), d.descriptionLbl.TextStyle)
	}
	return fyne.MeasureText(d.dr.typeName, theme.TextSize(), d.typeText.TextStyle)
}

func (d *detailRowRenderer) descriptionSize() fyne.Size {
	if d.dr.description == nil || *d.dr.description == "" {
		return fyne.MeasureText("Temp", theme.TextSize(), d.descriptionLbl.TextStyle)
	}
	return fyne.MeasureText(*d.dr.description, theme.TextSize(), d.descriptionLbl.TextStyle)
}

func (d *detailRowRenderer) deprecationSize() fyne.Size {
	if !d.dr.isDeprecated {
		return fyne.MeasureText("Temp", theme.TextSize(), d.deprecationLbl.TextStyle)
	}
	return fyne.MeasureText(*d.dr.deprecationReason, theme.TextSize(), d.deprecationLbl.TextStyle)
}

func (d *detailRowRenderer) defaultValueSize() fyne.Size {
	if d.dr.defaultValue == nil || *d.dr.defaultValue == "" {
		return fyne.MeasureText("Temp", theme.TextSize(), d.defaultLbl.TextStyle)
	}
	return fyne.MeasureText(*d.dr.defaultValue, theme.TextSize(), d.defaultLbl.TextStyle)
}

func (d *detailRowRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{d.nameLbl, d.typeText, d.descriptionLbl, d.cautionLbl, d.deprecationLbl, d.defaultTitle, d.defaultLbl}
}

func (d *detailRowRenderer) Refresh() {
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
		if d.dr.deprecationReason != nil {
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
}
