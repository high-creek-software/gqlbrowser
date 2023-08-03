package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/gqlbrowser/internal/resources"
)

type nameTypeRow struct {
	widget.BaseWidget

	name              string
	typ               string
	isDeprecated      bool
	deprecationReason *string
}

func (n *nameTypeRow) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (n *nameTypeRow) CreateRenderer() fyne.WidgetRenderer {
	cautionIcon := widget.NewIcon(resources.CautionResource)
	return &nameTypeRowRenderer{
		row:         n,
		nameLbl:     widget.NewLabel(n.name),
		typLbl:      canvas.NewText(n.typ, theme.PrimaryColor()),
		cautionIcon: cautionIcon,
	}
}

func newNameTypeRow(name, typ string) *nameTypeRow {
	r := &nameTypeRow{name: name, typ: typ}
	r.ExtendBaseWidget(r)

	return r
}

type nameTypeRowRenderer struct {
	row         *nameTypeRow
	nameLbl     *widget.Label
	typLbl      *canvas.Text
	cautionIcon *widget.Icon
}

func (n nameTypeRowRenderer) Destroy() {

}

func (n nameTypeRowRenderer) Layout(size fyne.Size) {
	nameSize := fyne.MeasureText(n.nameLbl.Text, theme.TextSize(), n.nameLbl.TextStyle)
	typeSize := fyne.MeasureText(n.typLbl.Text, theme.TextSize(), n.typLbl.TextStyle)
	cautionSize := fyne.NewSize(0, 0)
	if n.row.isDeprecated {
		cautionSize = fyne.NewSize(32, 32)
	}

	useTwoLines := cautionSize.Width+nameSize.Width+typeSize.Width+5*theme.Padding() > size.Width

	topLeft := fyne.NewPos(theme.Padding(), theme.Padding())

	if n.row.isDeprecated {
		n.cautionIcon.Resize(cautionSize)
		n.cautionIcon.Move(topLeft)
		topLeft = topLeft.AddXY(32, 0)
	}

	n.nameLbl.Move(topLeft)
	if useTwoLines {
		topLeft = fyne.NewPos(topLeft.X+8, nameSize.Height+2*theme.Padding()+10)
	} else {
		topLeft = topLeft.Add(fyne.NewPos(nameSize.Width+theme.Padding()+12, 8))
	}

	n.typLbl.Move(topLeft)
}

func (n nameTypeRowRenderer) MinSize() fyne.Size {
	// nameSize := fyne.MeasureText(n.nameLbl.Text, theme.TextSize(), n.nameLbl.TextStyle)
	// typSize := fyne.MeasureText(n.typLbl.Text, theme.TextSize(), n.typLbl.TextStyle)
	nameSize := n.nameLbl.MinSize()
	typSize := n.typLbl.MinSize()
	cautionSize := fyne.NewSize(0, 0)
	if n.row.isDeprecated {
		cautionSize = fyne.NewSize(32, 32)
	}

	width := fyne.Max(cautionSize.Width+nameSize.Width, typSize.Width)
	height := fyne.Max(nameSize.Height, typSize.Height)
	height = fyne.Max(height, cautionSize.Height)

	if n.typLbl.Position().Y > n.nameLbl.Position().Y+8 {
		height = nameSize.Height + typSize.Height + theme.Padding()
	}

	return fyne.NewSize(width+2*theme.Padding(), height+2*theme.Padding())
}

func (n nameTypeRowRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{n.nameLbl, n.typLbl, n.cautionIcon}
}

func (n nameTypeRowRenderer) Refresh() {
	n.nameLbl.SetText(n.row.name)
	n.typLbl.Text = n.row.typ
	n.typLbl.Refresh()

	if n.row.isDeprecated {
		n.cautionIcon.Show()
	} else {
		n.cautionIcon.Hide()
	}
}
