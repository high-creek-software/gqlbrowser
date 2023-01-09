package resources

import (
	_ "embed"
	"fyne.io/fyne/v2"
)

//go:embed caution.svg
var cautionBytes []byte

var CautionResource = fyne.NewStaticResource("caution", cautionBytes)

//go:embed icon.svg
var iconBytes []byte
var IconResource = fyne.NewStaticResource("icon", iconBytes)
