package ui

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/samuelyuan/PolytopiaMapEditor/internal/api"
	"github.com/samuelyuan/PolytopiaMapEditor/internal/fileio"
	"github.com/samuelyuan/PolytopiaMapEditor/internal/mapdraw"
)

type editor struct {
	drawSurface             *interactiveRaster
	status                  *widget.Label
	zoomLabel               *widget.Label
	leftPropertiesBar       fyne.CanvasObject
	tileProperties          TileProperties
	playerPropertiesBar     fyne.CanvasObject
	cache                   *image.RGBA
	cacheWidth, cacheHeight int

	// map properties
	uri       string
	img       *image.RGBA
	mapData   *fileio.PolytopiaSaveOutput
	mapHeight int
	mapWidth  int
	zoom      int

	tileX int
	tileY int

	win        fyne.Window
	recentMenu *fyne.Menu
}

func colorToBytes(col color.Color) []uint8 {
	r, g, b, a := col.RGBA()
	return []uint8{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
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

func (edit *editor) draw(w, h int) image.Image {
	if edit.cacheWidth == 0 || edit.cacheHeight == 0 {
		return image.NewRGBA(image.Rect(0, 0, w, h))
	}

	if w > edit.cacheWidth || h > edit.cacheHeight {
		bigger := image.NewRGBA(image.Rect(0, 0, w, h))
		draw.Draw(bigger, edit.cache.Bounds(), edit.cache, image.Point{}, draw.Over)
		return bigger
	}

	return edit.cache
}

func (e *editor) updateSizes() {
	if e.img == nil {
		return
	}
	e.cacheWidth = e.img.Bounds().Dx() * e.zoom
	e.cacheHeight = e.img.Bounds().Dy() * e.zoom

	c := fyne.CurrentApp().Driver().CanvasForObject(e.status)
	scale := float32(1.0)
	if c != nil {
		scale = c.Scale()
	}
	e.drawSurface.SetMinSize(fyne.NewSize(
		float32(e.cacheWidth)/scale,
		float32(e.cacheHeight)/scale))

	e.renderCache()
}

func (e *editor) pixAt(x, y int) []uint8 {
	ix := x / e.zoom
	iy := y / e.zoom

	if ix >= e.img.Bounds().Dx() || iy >= e.img.Bounds().Dy() {
		return []uint8{0, 0, 0, 0}
	}

	return colorToBytes(e.img.At(ix, iy))
}

func (e *editor) renderCache() {
	e.cache = image.NewRGBA(image.Rect(0, 0, e.cacheWidth, e.cacheHeight))
	for y := 0; y < e.cacheHeight; y++ {
		for x := 0; x < e.cacheWidth; x++ {
			i := (y*e.cacheWidth + x) * 4
			col := e.pixAt(x, y)
			e.cache.Pix[i] = col[0]
			e.cache.Pix[i+1] = col[1]
			e.cache.Pix[i+2] = col[2]
			e.cache.Pix[i+3] = col[3]
		}
	}

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

func (e *editor) LoadFile(read fyne.URIReadCloser) {
	defer read.Close()

	inputFilename := read.URI().String()[7:]
	fmt.Println("Input filename: ", inputFilename)

	fileio.DecompressFile(inputFilename)
	saveFileData, err := fileio.ReadPolytopiaDecompressedFile(inputFilename + ".decomp")
	if err != nil {
		fyne.LogError("Failed to read input file: ", err)
		e.status.SetText(err.Error())
		return
	}

	img := mapdraw.DrawMap(saveFileData, e.tileX, e.tileY)

	e.addRecent(read.URI())
	e.uri = read.URI().String()
	e.img = fixEncoding(img)
	e.mapData = saveFileData
	e.mapHeight = len(saveFileData.TileData)
	e.mapWidth = len(saveFileData.TileData[0])
	e.status.SetText(fmt.Sprintf("File: %s | Map Rows: %d | Map Cols: %d",
		filepath.Base(read.URI().String()), e.mapHeight, e.mapWidth))
	e.tileProperties.UpdateOwnerOptions(e)

	content := e.createContainer()
	e.win.SetContent(content)

	e.updateSizes()
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

	fileInfo := fileio.FileInfo{
		InputFilename: decompressedFilename,
		GameVersion:   int(edit.mapData.GameVersion),
	}
	fileio.WriteMapToFile(fileInfo, edit.mapData.TileData)
	fileio.WritePlayersToFile(decompressedFilename, edit.mapData.PlayerData)
	fileio.WriteMapHeaderToFile(decompressedFilename, edit.mapData.MapHeaderOutput)
	fmt.Println("Exporting map to", outputFilename)
	fileio.CompressFile(decompressedFilename, outputFilename)
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
	return strings.LastIndex(strings.ToLower(path), "png") == len(path)-3
}

func (edit *editor) updateZoom(val int) {
	if val < 1 {
		val = 1
	} else if val > 16 {
		val = 16
	}
	edit.setZoom(val)
	edit.zoomLabel.SetText(fmt.Sprintf("%d%%", edit.zoom*100))
}

func (edit *editor) createContainer() fyne.CanvasObject {
	toolbar := buildToolbar(edit)
	edit.leftPropertiesBar = newLeftPropertiesBar(edit)
	edit.playerPropertiesBar = newPlayerPropertiesBar(edit)
	return fyne.NewContainerWithLayout(
		layout.NewBorderLayout(toolbar, edit.status, edit.leftPropertiesBar, edit.playerPropertiesBar),
		toolbar, edit.status, edit.leftPropertiesBar, edit.buildUI(), edit.playerPropertiesBar,
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

	options := append([]fyne.CanvasObject{
		container.NewGridWithColumns(1),
		zoom,
	})
	options = append(options, edit.tileProperties.GetOptions()...)

	return container.NewVBox(options...)
}

func (edit *editor) refreshMapImage() {
	img := mapdraw.DrawMap(edit.mapData, edit.tileX, edit.tileY)
	edit.img = fixEncoding(img)
	edit.updateSizes()
}

// NewEditor creates a new pixel editor that is ready to have a file loaded
func NewEditor() api.Editor {
	edit := &editor{
		zoom:      1,
		zoomLabel: widget.NewLabel("100%"),
		status:    widget.NewLabel("Open a file"),
		tileX:     -1,
		tileY:     -1,
	}
	edit.drawSurface = newInteractiveRaster(edit)
	edit.tileProperties = NewTileProperties(edit)

	return edit
}
