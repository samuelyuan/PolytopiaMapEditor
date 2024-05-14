package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type palette struct {
	edit *editor

	zoom *widget.Label
}

func (p *palette) updateZoom(val int) {
	if val < 1 {
		val = 1
	} else if val > 16 {
		val = 16
	}
	p.edit.setZoom(val)

	p.zoom.SetText(fmt.Sprintf("%d%%", p.edit.zoom*100))
}

func newPalette(edit *editor) fyne.CanvasObject {
	p := &palette{edit: edit, zoom: widget.NewLabel("100%")}

	zoom := container.NewHBox(
		widget.NewButtonWithIcon("", theme.ZoomOutIcon(), func() {
			p.updateZoom(p.edit.zoom / 2)
		}),
		p.zoom,
		widget.NewButtonWithIcon("", theme.ZoomInIcon(), func() {
			p.updateZoom(p.edit.zoom * 2)
		}))

	return container.NewVBox(append([]fyne.CanvasObject{
		container.NewGridWithColumns(1),
		zoom,
		edit.tileCoordinatesProperties,
		container.NewHBox(widget.NewLabel("Terrain"), edit.terrainSelect),
		edit.mapTileProperties,
	})...)
}
