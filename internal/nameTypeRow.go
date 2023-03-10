package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type nameTypeRow struct {
	widget.BaseWidget

	name string
	typ  string
}

func (n *nameTypeRow) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (n *nameTypeRow) CreateRenderer() fyne.WidgetRenderer {
	return &nameTypeRowRenderer{
		row:     n,
		nameLbl: widget.NewLabel(n.name),
		typLbl:  canvas.NewText(n.typ, theme.PrimaryColor()),
	}
}

func newNameTypeRow(name, typ string) *nameTypeRow {
	r := &nameTypeRow{name: name, typ: typ}
	r.ExtendBaseWidget(r)

	return r
}

type nameTypeRowRenderer struct {
	row     *nameTypeRow
	nameLbl *widget.Label
	typLbl  *canvas.Text
}

func (n nameTypeRowRenderer) Destroy() {

}

func (n nameTypeRowRenderer) Layout(size fyne.Size) {
	nameSize := fyne.MeasureText(n.nameLbl.Text, theme.TextSize(), n.nameLbl.TextStyle)

	topLeft := fyne.NewPos(theme.Padding(), theme.Padding())
	n.nameLbl.Move(topLeft)
	topLeft = topLeft.Add(fyne.NewPos(nameSize.Width+theme.Padding()+12, 8))
	n.typLbl.Move(topLeft)
}

func (n nameTypeRowRenderer) MinSize() fyne.Size {
	nameSize := fyne.MeasureText(n.nameLbl.Text, theme.TextSize(), n.nameLbl.TextStyle)
	typSize := fyne.MeasureText(n.typLbl.Text, theme.TextSize(), n.typLbl.TextStyle)

	return fyne.NewSize(nameSize.Width+typSize.Width+3*theme.Padding()+30, fyne.Max(nameSize.Height, typSize.Height)+4*theme.Padding())
}

func (n nameTypeRowRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{n.nameLbl, n.typLbl}
}

func (n nameTypeRowRenderer) Refresh() {
	n.nameLbl.SetText(n.row.name)
	n.typLbl.Text = n.row.typ
	n.typLbl.Refresh()
}
