package internal

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/fieldglass"
)

type detailLayout struct {
	*fyne.Container

	title          *widget.RichText
	closeBtn       *widget.Button
	segmentWrapper *fyne.Container

	remove       func(container *fyne.Container)
	typeSelected func(t fieldglass.Type)
}

func newDetailLayout(title string, subTitle *string, remove func(container *fyne.Container), typeSelected func(t fieldglass.Type)) *detailLayout {
	dl := &detailLayout{typeSelected: typeSelected}
	dl.title = widget.NewRichTextFromMarkdown(fmt.Sprintf("# %s", title))
	dl.closeBtn = widget.NewButtonWithIcon("Close", theme.ContentRemoveIcon(), func() { remove(dl.Container) })
	dl.segmentWrapper = container.NewMax()

	titleWrapper := container.NewVBox(container.NewHBox(widget.NewLabel("            "), dl.title, widget.NewLabel("            ")))
	sub := ""
	if subTitle != nil && *subTitle != "" {
		sub = *subTitle
	}
	titleWrapper.Add(widget.NewLabel(sub))
	titleWrapper.Add(widget.NewSeparator())

	child := container.NewScroll(dl.segmentWrapper)
	child.Direction = container.ScrollVerticalOnly
	dl.Container = container.NewPadded(container.NewBorder(titleWrapper,
		dl.closeBtn,
		nil,
		nil,
		container.NewPadded(child),
	),
	)

	return dl
}

func (dl *detailLayout) buildArgs(args []fieldglass.InputValue) *fyne.Container {
	ia := &inputAdapter{inputs: args}
	inputList := widget.NewList(ia.count, ia.createTemplate, ia.updateTemplate)
	inputBorder := container.NewBorder(widget.NewRichTextFromMarkdown("## Arguments"), nil, nil, nil, inputList)

	inputList.OnSelected = func(id widget.ListItemID) {
		i := ia.getItem(id)
		if i.Type.RootType() == fieldglass.TypeKindScalar {
			return
		}
		dl.typeSelected(*i.Type)
	}

	return inputBorder
}

func (dl *detailLayout) buildProperties(t *fieldglass.Type) *fyne.Container {
	adapter := &fullFieldAdapter{fields: t.Fields}
	list := widget.NewList(adapter.count, adapter.createTemplate, adapter.updateTemplate)
	propertiesBorder := container.NewBorder(widget.NewRichTextFromMarkdown("## Properties"), nil, nil, nil, list)
	list.OnSelected = func(id widget.ListItemID) {
		f := adapter.getItem(id)
		if f.Type.RootType() == fieldglass.TypeKindScalar {
			return
		}
		dl.typeSelected(*f.Type)
	}

	return propertiesBorder
}

func (dl *detailLayout) buildTypes(name string, ts []fieldglass.Type) *fyne.Container {
	typeSegment := container.NewVBox()
	for _, typ := range ts {
		func(t fieldglass.Type) {
			name := widget.NewHyperlink(*t.Name+":"+t.FormatName(), nil)
			name.OnTapped = func() {
				dl.typeSelected(t)
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
