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
	dl.segmentWrapper = container.NewVBox()

	titleWrapper := container.NewVBox(container.NewHBox(widget.NewLabel("            "), dl.title, widget.NewLabel("            ")))
	sub := ""
	if subTitle != nil && *subTitle != "" {
		sub = *subTitle
	}
	titleWrapper.Add(widget.NewLabel(sub))
	titleWrapper.Add(widget.NewSeparator())

	//dl.segmentWrapper = container.NewGridWithColumns(1)
	//if len(rootType.Fields) > 0 && len(args) > 0 {
	//	child = container.NewGridWithRows(2, inputBorder, propertiesBorder)
	//} else if len(args) > 0 {
	//	child = inputBorder
	//} else if len(rootType.Fields) > 0 {
	//	child = propertiesBorder
	//}
	child := container.NewScroll(dl.segmentWrapper)
	child.Direction = container.ScrollVerticalOnly
	dl.Container = container.NewBorder(titleWrapper,
		dl.closeBtn,
		nil,
		nil,
		container.NewPadded(child),
	)

	return dl
}

func (dl *detailLayout) addArgs(args []fieldglass.InputValue) {
	//ia := &inputAdapter{inputs: args}
	//inputList := widget.NewList(ia.count, ia.createTemplate, ia.updateTemplate)
	//inputBorder := container.NewBorder(widget.NewRichTextFromMarkdown("## Arguments"), nil, nil, nil, inputList)
	//dl.segmentWrapper.Add(inputBorder)
	//dl.segmentWrapper.Refresh()

	argSegment := container.NewVBox()
	for _, arg := range args {
		func(a fieldglass.InputValue) {
			name := widget.NewHyperlink(arg.Name+":"+arg.Type.FormatName(), nil)
			name.OnTapped = func() {
				dl.typeSelected(*a.Type)
			}
			if arg.Description != nil && *arg.Description != "" {
				desc := widget.NewLabel(*arg.Description)
				desc.TextStyle = fyne.TextStyle{Italic: true}
				argSegment.Add(container.NewVBox(name, desc))
			} else {
				argSegment.Add(name)
			}
		}(arg)
	}

	dl.segmentWrapper.Add(container.NewBorder(widget.NewRichTextFromMarkdown("## Arguments"), nil, nil, nil, argSegment))
}

func (dl *detailLayout) addProperties(t *fieldglass.Type) {
	//adapter := &fieldAdapter{fields: t.Fields}
	//list := widget.NewList(adapter.count, adapter.createTemplate, adapter.updateTemplate)
	//propertiesBorder := container.NewBorder(widget.NewRichTextFromMarkdown("## Properties"), nil, nil, nil, list)
	//list.OnSelected = func(id widget.ListItemID) {
	//	f := adapter.getItem(id)
	//	if f.Type.RootType() == fieldglass.TypeKindScalar {
	//		return
	//	}
	//	dl.typeSelected(*f.Type)
	//}
	//dl.segmentWrapper.Add(propertiesBorder)

	fieldSegment := container.NewVBox()
	for _, fld := range t.Fields {
		func(f fieldglass.Field) {
			name := widget.NewHyperlink(f.Name+":"+f.Type.FormatName(), nil)
			name.OnTapped = func() {
				dl.typeSelected(*f.Type)
			}
			if f.Description != nil && *f.Description != "" {
				desc := widget.NewLabel(*f.Description)
				desc.TextStyle = fyne.TextStyle{Italic: true}
				fieldSegment.Add(container.NewVBox(name, desc))
			} else {
				fieldSegment.Add(name)
			}
		}(fld)
	}
	dl.segmentWrapper.Add(container.NewBorder(widget.NewRichTextFromMarkdown("## Fields"), nil, nil, nil, fieldSegment))
}

func (dl *detailLayout) addTypes(name string, ts []fieldglass.Type) {
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
	dl.segmentWrapper.Add(container.NewBorder(widget.NewRichTextFromMarkdown(fmt.Sprintf("## %s", name)), nil, nil, nil, typeSegment))
}

func (dl *detailLayout) addEnums(es []fieldglass.EnumValue) {
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

	dl.segmentWrapper.Add(container.NewBorder(widget.NewRichTextFromMarkdown("## Enum Values"), nil, nil, nil, enumSegment))
}
