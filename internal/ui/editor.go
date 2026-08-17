package ui

import (
	"fmt"
	"image"
	"image/png"
	"path/filepath"

	"golang.org/x/image/draw"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/samuelyuan/PolytopiaMapEditor/internal/api"
	"github.com/samuelyuan/PolytopiaMapEditor/internal/mapdraw"
	polytopiamapmodel "github.com/samuelyuan/polytopiamapmodelgo"
)

type editor struct {
	drawSurface                      *interactiveRaster
	status                           *widget.Label
	zoomLabel                        *widget.Label
	showResourcesImprovementCheckbox *widget.Check
	leftPropertiesBar                fyne.CanvasObject
	tileProperties                   TileProperties
	playerPropertiesBar              fyne.CanvasObject

	// map properties
	uri       string
	img       *image.RGBA
	mapData   *polytopiamapmodel.PolytopiaSaveOutput
	mapHeight int
	mapWidth  int
	zoom      int

	tileX           int
	tileY           int
	graphicsOptions mapdraw.GraphicsOptions

	win        fyne.Window
	recentMenu *fyne.Menu
}

func (e *editor) Clicked(x, y int, edit api.Editor) {
	edit.SetTileCoordinates(x, y)
}

func (e *editor) SetTileCoordinates(x int, y int) {
	tileX, tileY := mapdraw.GetTileCoordinates(x, y, e.mapHeight, e.mapWidth)
	e.tileX = tileX
	e.tileY = tileY
	e.tileProperties.UpdateTileProperties(tileX, tileY, e.mapData.TileData[tileY][tileX])
}

func (edit *editor) buildUI() fyne.CanvasObject {
	return container.NewScroll(edit.drawSurface)
}

func (edit *editor) setZoom(zoom int) {
	edit.zoom = zoom
	edit.updateSizes()
	edit.drawSurface.Refresh()
}

// updateSizes pushes edit.img to the GPU-backed image widget and resizes it for the current zoom.
func (e *editor) updateSizes() {
	if e.img == nil {
		return
	}
	physWidth := e.img.Bounds().Dx() * e.zoom
	physHeight := e.img.Bounds().Dy() * e.zoom

	c := fyne.CurrentApp().Driver().CanvasForObject(e.status)
	scale := float32(1.0)
	if c != nil {
		scale = c.Scale()
	}

	e.drawSurface.img.Image = e.img
	e.drawSurface.img.Refresh()
	e.drawSurface.SetMinSize(fyne.NewSize(
		float32(physWidth)/scale,
		float32(physHeight)/scale))

	e.drawSurface.Refresh()
}

func fixEncoding(img image.Image) *image.RGBA {
	if rgba, ok := img.(*image.RGBA); ok {
		return rgba
	}

	newImg := image.NewRGBA(img.Bounds())
	draw.Draw(newImg, newImg.Bounds(), img, img.Bounds().Min, draw.Over)
	return newImg
}

func (edit *editor) LoadFile(read fyne.URIReadCloser) {
	defer read.Close()

	inputFilename := read.URI().String()[7:]
	fmt.Println("Input filename: ", inputFilename)

	polytopiamapmodel.DecompressFile(inputFilename)
	saveFileData, err := polytopiamapmodel.ReadPolytopiaDecompressedFile(inputFilename + ".decomp")
	if err != nil {
		fyne.LogError("Failed to read input file: ", err)
		edit.status.SetText(err.Error())
		return
	}

	img := getNewImage(edit, saveFileData)

	edit.addRecent(read.URI())
	edit.uri = read.URI().String()
	edit.img = fixEncoding(img)
	edit.mapData = saveFileData
	edit.mapHeight = len(saveFileData.TileData)
	edit.mapWidth = len(saveFileData.TileData[0])
	edit.status.SetText(fmt.Sprintf("File: %s | Map Rows: %d | Map Cols: %d | Game Version: %d",
		filepath.Base(read.URI().String()), edit.mapHeight, edit.mapWidth, edit.mapData.GameVersion))
	edit.tileProperties.UpdateOwnerOptions(edit)

	content := edit.createContainer()
	edit.win.SetContent(content)

	edit.updateSizes()
}

func (e *editor) Reload() {
	if e.uri == "" {
		return
	}

	u, _ := storage.ParseURI(e.uri)
	read, err := storage.Reader(u)
	if err != nil {
		fyne.LogError("Unable to open file \""+e.uri+"\"", err)
		return
	}
	e.LoadFile(read)
}

func (edit *editor) WriteMapState() {
	originalFilname := edit.uri[7:]
	decompressedFilename := originalFilname + ".decomp"
	outputFilename := originalFilname

	fileInfo := polytopiamapmodel.FileInfo{
		InputFilename: decompressedFilename,
		GameVersion:   int(edit.mapData.GameVersion),
	}
	polytopiamapmodel.WriteMapToFile(fileInfo, edit.mapData.TileData)
	polytopiamapmodel.WritePlayersToFile(decompressedFilename, edit.mapData.PlayerData, fileInfo.GameVersion)
	polytopiamapmodel.WriteMapHeaderToFile(decompressedFilename, edit.mapData.MapHeaderOutput)
	fmt.Println("Exporting map to", outputFilename)
	polytopiamapmodel.CompressFile(decompressedFilename, outputFilename)
}

func (e *editor) Save() {
	if e.uri == "" {
		return
	}

	uri, _ := storage.ParseURI(e.uri)
	if !e.isSupported(uri.Extension()) {
		fyne.LogError("Save only supports PNG", nil)
		return
	}
	write, err := storage.Writer(uri)
	if err != nil {
		fyne.LogError("Error opening file to replace", err)
		return
	}

	e.saveToWriter(write)
}

func (e *editor) saveToWriter(write fyne.URIWriteCloser) {
	defer write.Close()
	if e.isPNG(write.URI().Extension()) {
		err := png.Encode(write, e.img)

		if err != nil {
			fyne.LogError("Could not encode image", err)
		}
	}
}

func (e *editor) SaveAs(writer fyne.URIWriteCloser) {
	e.saveToWriter(writer)
}

func (e *editor) isSupported(path string) bool {
	return e.isPNG(path)
}

func (e *editor) isPNG(path string) bool {
	return isPNGPath(path)
}

func (edit *editor) updateZoom(val int) {
	edit.setZoom(clampZoom(val))
	edit.zoomLabel.SetText(fmt.Sprintf("%d%%", edit.zoom*100))
}

// Default split offsets for the property panel dividers.
const (
	leftPanelDefaultOffset  = 0.2
	rightPanelDefaultOffset = 0.8
)

func (edit *editor) createContainer() fyne.CanvasObject {
	toolbar := buildToolbar(edit)
	edit.leftPropertiesBar = newLeftPropertiesBar(edit)
	edit.playerPropertiesBar = newPlayerPropertiesBar(edit)

	mapAndRightSplit := container.NewHSplit(edit.buildUI(), edit.playerPropertiesBar)
	mapAndRightSplit.SetOffset(rightPanelDefaultOffset)

	leftAndRestSplit := container.NewHSplit(edit.leftPropertiesBar, mapAndRightSplit)
	leftAndRestSplit.SetOffset(leftPanelDefaultOffset)

	return fyne.NewContainerWithLayout(
		layout.NewBorderLayout(toolbar, edit.status, nil, nil),
		toolbar, edit.status, leftAndRestSplit,
	)
}

func newLeftPropertiesBar(edit *editor) fyne.CanvasObject {
	zoom := container.NewHBox(
		widget.NewButtonWithIcon("", theme.ZoomOutIcon(), func() {
			edit.updateZoom(edit.zoom / 2)
		}),
		edit.zoomLabel,
		widget.NewButtonWithIcon("", theme.ZoomInIcon(), func() {
			edit.updateZoom(edit.zoom * 2)
		}))

	graphicsOptionsResourcesImprovements := container.NewVBox(
		widget.NewCheck("Show resources and improvements", func(value bool) {
			edit.graphicsOptions.ShowResourcesImprovements = value
			if edit.mapData != nil {
				edit.refreshMapImage()
			}
		}),
		widget.NewCheck("Show roads", func(value bool) {
			edit.graphicsOptions.ShowRoads = value
			if edit.mapData != nil {
				edit.refreshMapImage()
			}
		}),
		widget.NewCheck("Show units", func(value bool) {
			edit.graphicsOptions.ShowUnits = value
			if edit.mapData != nil {
				edit.refreshMapImage()
			}
		}),
	)

	options := []fyne.CanvasObject{
		container.NewGridWithColumns(1),
		zoom,
		graphicsOptionsResourcesImprovements,
	}
	options = append(options, edit.tileProperties.GetOptions()...)

	// NewScroll (not NewVScroll) so the panel can be narrowed via the split divider.
	return container.NewScroll(container.NewVBox(options...))
}

func (edit *editor) refreshMapImage() {
	img := getUpdatedImage(edit)
	edit.img = fixEncoding(img)
	edit.updateSizes()
}

func getNewImage(edit *editor, newMapData *polytopiamapmodel.PolytopiaSaveOutput) image.Image {
	return mapdraw.DrawMap(
		newMapData,
		edit.tileX,
		edit.tileY,
		edit.graphicsOptions,
	)
}

func getUpdatedImage(edit *editor) image.Image {
	return getNewImage(edit, edit.mapData)
}

// NewEditor creates a new pixel editor that is ready to have a file loaded
func NewEditor() api.Editor {
	edit := &editor{
		zoom:      1,
		zoomLabel: widget.NewLabel("100%"),
		status:    widget.NewLabel("Open a file"),
		tileX:     -1,
		tileY:     -1,
		graphicsOptions: mapdraw.GraphicsOptions{
			ShowResourcesImprovements: false,
			ShowRoads:                 false,
			ShowUnits:                 false,
		},
	}
	edit.drawSurface = newInteractiveRaster(edit)
	edit.tileProperties = NewTileProperties(edit)

	return edit
}
