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
)

type nameTypeRow struct {
	widget.BaseWidget

	//name              string
	//typ               string
	//isDeprecated      bool
	//deprecationReason *string

	nameLbl         *widget.Label
	typeTxt         *canvas.Text
	deprecatedIcon  *widget.Icon
	deprecatedText  *widget.Label
	descriptionText *widget.Label
}

func (n *nameTypeRow) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (n *nameTypeRow) CreateRenderer() fyne.WidgetRenderer {
	n.nameLbl = widget.NewLabel("")
	n.typeTxt = canvas.NewText("", theme.Color(theme.ColorNamePrimary))
	n.deprecatedIcon = widget.NewIcon(resources.CautionResource)
	n.deprecatedText = widget.NewLabel("")
	n.deprecatedText.SizeName = theme.SizeNameCaptionText
	n.deprecatedText.Wrapping = fyne.TextWrapWord
	n.deprecatedText.TextStyle = fyne.TextStyle{Italic: true}
	n.descriptionText = widget.NewLabel("")
	n.descriptionText.SizeName = theme.SizeNameCaptionText
	n.descriptionText.Wrapping = fyne.TextWrapWord

	nameBox := container.NewHBox(n.nameLbl, n.typeTxt)
	//nameBox := container.NewGridWrap(fyne.Size{Height: 32, Width: 64}, n.nameLbl, n.typeTxt)
	depBox := container.New(
		layout.NewCustomPaddedLayout(0, 0, 8, 0),
		container.NewBorder(nil, nil, n.deprecatedIcon, nil, n.deprecatedText),
	)

	return widget.NewSimpleRenderer(
		container.NewVBox(
			nameBox,
			depBox,
			n.descriptionText,
		),
	)
}

func (n *nameTypeRow) SetData(name, typ, deprecationReason, description string, isDeprecated bool) {
	n.nameLbl.SetText(name)
	n.typeTxt.Text = typ
	n.deprecatedText.SetText(deprecationReason)
	n.descriptionText.SetText(description)
	if isDeprecated {
		n.deprecatedIcon.Show()
		if deprecationReason != "" {
			n.deprecatedText.Show()
		} else {
			n.deprecatedText.Hide()
		}
	} else {
		n.deprecatedIcon.Hide()
		n.deprecatedText.Hide()
	}

	if description == "" {
		n.descriptionText.Hide()
	} else {
		n.descriptionText.Show()
	}
}

func newNameTypeRow() *nameTypeRow {
	r := &nameTypeRow{}
	r.ExtendBaseWidget(r)

	return r
}
