package api

import (
	"fyne.io/fyne/v2"
)

// Editor describes the editing capabilities
type Editor interface {
	BuildUI(fyne.Window)         // BuildUI Loads the main editor GUI
	LoadFile(fyne.URIReadCloser) // LoadFile specifies a data stream to load from
	Reload()                     // Reload will reset the image to its original state
	Save()                       // Save writes the image back to its source location
	SaveAs(fyne.URIWriteCloser)  // SaveAs specifies a data stream to save to

	SetHexCoordinates(x int, y int)
}
