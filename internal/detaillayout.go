package internal

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/gqlbrowser/internal/resources"
	"gitlab.com/high-creek-software/fieldglass"
)

type detailLayout struct {
	widget.BaseWidget

	titleRT        *widget.RichText
	subTitle       *string
	subTitleRT     *widget.RichText
	closeBtn       *widget.Button
	segmentWrapper *fyne.Container

	isDeprecated      bool
	cautionIcon       *widget.Icon
	deprecationReason *widget.Label

	remove       func(container *fyne.Container)
	typeSelected func(t fieldglass.Type, f *fieldglass.Field)
}

func (dl *detailLayout) CreateRenderer() fyne.WidgetRenderer {
	titleWrapper := container.NewBorder(nil, nil, nil, dl.closeBtn, container.NewHBox(layout.NewSpacer(), widget.NewLabel("			"), dl.titleRT, widget.NewLabel("			"), layout.NewSpacer()))

	vBox := container.NewVBox()
	if dl.subTitle != nil {
		vBox.Add(dl.subTitleRT)
	}
	if dl.isDeprecated {
		vBox.Add(container.NewBorder(nil, nil, dl.cautionIcon, nil, dl.deprecationReason))
	}

	titleBorder := container.NewBorder(titleWrapper, widget.NewSeparator(), nil, nil, vBox)

	//if dl.isDeprecated {
	//	titleBorder = container.NewBorder(titleBorder, container.NewBorder(nil, nil, dl.cautionIcon, nil, dl.deprecationReason), nil, nil)
	//}

	child := container.NewScroll(dl.segmentWrapper)
	child.Direction = container.ScrollVerticalOnly

	// dl.Resize(fyne.NewSize(450, dl.Container.MinSize().Height))
	return widget.NewSimpleRenderer(container.NewPadded(container.NewBorder(titleBorder, nil, nil, nil, container.NewPadded(child))))
}

func newDetailLayout(title string, subTitle *string, isDeprecated bool, deprecationReason *string, remove func(container fyne.CanvasObject), typeSelected func(t fieldglass.Type, f *fieldglass.Field)) *detailLayout {
	dl := &detailLayout{typeSelected: typeSelected, isDeprecated: isDeprecated, subTitle: subTitle}
	dl.ExtendBaseWidget(dl)

	dl.titleRT = widget.NewRichTextFromMarkdown(fmt.Sprintf("# %s", title))
	sub := ""
	if subTitle != nil && *subTitle != "" {
		sub = *subTitle
	}
	dl.subTitleRT = widget.NewRichTextFromMarkdown(fmt.Sprintf("### %s", sub))
	dl.closeBtn = widget.NewButton("X", func() { remove(dl) })
	dl.closeBtn.Importance = widget.LowImportance
	dl.segmentWrapper = container.NewMax()

	dl.cautionIcon = widget.NewIcon(resources.CautionResource)
	dep := ""
	if deprecationReason != nil {
		dep = *deprecationReason
	}
	dl.deprecationReason = widget.NewLabel(dep)

	return dl
}

func (dl *detailLayout) buildArgs(args []fieldglass.InputValue) *fyne.Container {
	ia := &inputAdapter{inputs: args}
	fll := newFilterableListLayout[fieldglass.InputValue](ia, func(id widget.ListItemID) {
		i := ia.getItem(id)
		if i.Type.RootType() == fieldglass.TypeKindScalar {
			return
		}
		dl.typeSelected(*i.Type, nil)
	})
	inputBorder := container.NewBorder(widget.NewRichTextFromMarkdown("## Arguments"), nil, nil, nil, fll)

	return inputBorder
}

func (dl *detailLayout) buildProperties(t *fieldglass.Type) *fyne.Container {
	adapter := &fullFieldAdapter{fields: t.Fields}
	fll := newFilterableListLayout[fieldglass.Field](adapter, func(id widget.ListItemID) {
		f := adapter.getItem(id)
		if f.Type.RootType() == fieldglass.TypeKindScalar {
			return
		}
		dl.typeSelected(*f.Type, &f)
	})
	propertiesBorder := container.NewBorder(widget.NewRichTextFromMarkdown("## Properties"), nil, nil, nil, fll)

	return propertiesBorder
}

func (dl *detailLayout) buildTypes(name string, ts []fieldglass.Type) *fyne.Container {
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
	return container.NewBorder(widget.NewRichTextFromMarkdown(fmt.Sprintf("## %s", name)), nil, nil, nil, typeSegment)
}

func (dl *detailLayout) buildEnums(es []fieldglass.EnumValue) *fyne.Container {
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

	return container.NewBorder(widget.NewRichTextFromMarkdown("## Enum Values"), nil, nil, nil, enumSegment)
}
