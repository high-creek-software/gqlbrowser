package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/gqlbrowser/internal/resources"
	"gitlab.com/high-creek-software/fieldglass"
	"slices"
	"strings"
)

type detailLayout struct {
	widget.BaseWidget

	titleLbl       *widget.Label
	subTitle       *string
	subTitleLbl    *widget.Label
	closeBtn       *widget.Button
	segmentWrapper *fyne.Container

	isDeprecated      bool
	cautionIcon       *widget.Icon
	deprecationReason *widget.Label

	remove       func(container *fyne.Container)
	typeSelected func(t fieldglass.Type, f *fieldglass.Field)
}

func (dl *detailLayout) CreateRenderer() fyne.WidgetRenderer {
	titleWrapper := container.NewBorder(
		nil,
		nil,
		nil,
		dl.closeBtn,
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewLabel("			"),
			dl.titleLbl,
			widget.NewLabel("			"),
			layout.NewSpacer(),
		),
	)

	vBox := container.NewVBox()
	if dl.subTitle != nil {
		vBox.Add(dl.subTitleLbl)
	}
	if dl.isDeprecated {
		vBox.Add(container.NewBorder(nil, nil, dl.cautionIcon, nil, dl.deprecationReason))
	}

	titleBorder := container.NewBorder(titleWrapper, widget.NewSeparator(), nil, nil, vBox)

	child := container.NewScroll(dl.segmentWrapper)
	child.Direction = container.ScrollVerticalOnly

	return widget.NewSimpleRenderer(container.NewPadded(container.NewBorder(titleBorder, nil, nil, nil, container.NewPadded(child))))
}

func newDetailLayout(title string, subTitle *string, isDeprecated bool, deprecationReason *string, remove func(container fyne.CanvasObject), typeSelected func(t fieldglass.Type, f *fieldglass.Field)) *detailLayout {
	dl := &detailLayout{typeSelected: typeSelected, isDeprecated: isDeprecated, subTitle: subTitle}
	dl.ExtendBaseWidget(dl)

	dl.titleLbl = widget.NewLabel(title)
	dl.titleLbl.SizeName = theme.SizeNameHeadingText
	sub := ""
	if subTitle != nil && *subTitle != "" {
		sub = *subTitle
	}

	dl.subTitleLbl = widget.NewLabel(sub)
	dl.subTitleLbl.SizeName = theme.SizeNameCaptionText
	dl.closeBtn = widget.NewButton("X", func() { remove(dl) })
	dl.closeBtn.Importance = widget.LowImportance
	dl.segmentWrapper = container.NewStack()

	dl.cautionIcon = widget.NewIcon(resources.CautionResource)
	dep := ""
	if deprecationReason != nil {
		dep = *deprecationReason
	}
	dl.deprecationReason = widget.NewLabel(dep)
	dl.deprecationReason.SizeName = theme.SizeNameCaptionText
	dl.deprecationReason.TextStyle = fyne.TextStyle{Italic: true}

	return dl
}

func (dl *detailLayout) buildArgs(args []fieldglass.InputValue) *fyne.Container {
	slices.SortFunc(args, func(a, b fieldglass.InputValue) int {
		return strings.Compare(a.Name, b.Name)
	})
	ia := &inputAdapter{inputs: args}
	fll := newFilterableListLayout[fieldglass.InputValue](ia, func(id widget.ListItemID) {
		i := ia.getItem(id)
		if i.Type.RootType() == fieldglass.TypeKindScalar {
			return
		}
		dl.typeSelected(*i.Type, nil)
	})
	argsLbl := widget.NewLabel("Arguments")
	argsLbl.SizeName = theme.SizeNameSubHeadingText
	inputBorder := container.NewBorder(argsLbl, nil, nil, nil, fll)

	return inputBorder
}

func (dl *detailLayout) buildProperties(t *fieldglass.Type) *fyne.Container {
	slices.SortFunc(t.Fields, func(a fieldglass.Field, b fieldglass.Field) int {
		return strings.Compare(a.Name, b.Name)
	})
	adapter := &fullFieldAdapter{fields: t.Fields}
	fll := newFilterableListLayout[fieldglass.Field](adapter, func(id widget.ListItemID) {
		f := adapter.getItem(id)
		if f.Type.RootType() == fieldglass.TypeKindScalar {
			return
		}
		dl.typeSelected(*f.Type, &f)
	})
	propsLbl := widget.NewLabel("Properties")
	propsLbl.SizeName = theme.SizeNameSubHeadingText
	propertiesBorder := container.NewBorder(propsLbl, nil, nil, nil, fll)

	return propertiesBorder
}

func (dl *detailLayout) buildTypes(name string, ts []fieldglass.Type) *fyne.Container {
	slices.SortFunc(ts, func(a, b fieldglass.Type) int {
		aName := ""
		if a.Name != nil {
			aName = *a.Name
		}
		bName := ""
		if b.Name != nil {
			bName = *b.Name
		}
		return strings.Compare(aName, bName)
	})
	typeSegment := container.NewVBox()
	for _, typ := range ts {
		func(t fieldglass.Type) {
			name := widget.NewHyperlink(*t.Name+":"+t.FormatName(), nil)
			name.OnTapped = func() {
				dl.typeSelected(t, nil)
			}
			if t.Description != nil && *t.Description != "" {
				desc := widget.NewLabel(*t.Description)
				desc.TextStyle = fyne.TextStyle{Italic: true}
				typeSegment.Add(container.NewVBox(name, desc))
			} else {
				typeSegment.Add(name)
			}
		}(typ)
	}

	nameLbl := widget.NewLabel(name)
	nameLbl.SizeName = theme.SizeNameSubHeadingText

	return container.NewBorder(nameLbl, nil, nil, nil, typeSegment)
}

func (dl *detailLayout) buildEnums(es []fieldglass.EnumValue) *fyne.Container {

	slices.SortFunc(es, func(a, b fieldglass.EnumValue) int {
		return strings.Compare(a.Name, b.Name)
	})

	enumSegment := container.NewVBox()
	for _, e := range es {
		func(en fieldglass.EnumValue) {
			name := widget.NewLabel(en.Name)
			if en.Description != nil && *en.Description != "" {
				desc := widget.NewLabel(*en.Description)
				desc.TextStyle = fyne.TextStyle{Italic: true}
				enumSegment.Add(container.NewVBox(name, desc))
			} else {
				enumSegment.Add(name)
			}
		}(e)
	}
	enumLbl := widget.NewLabel("Enum Values")
	enumLbl.SizeName = theme.SizeNameSubHeadingText
	return container.NewBorder(enumLbl, nil, nil, nil, enumSegment)
}
