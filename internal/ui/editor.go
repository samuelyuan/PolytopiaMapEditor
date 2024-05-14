package ui

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/image/draw"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/samuelyuan/PolytopiaMapEditor/internal/api"
	"github.com/samuelyuan/PolytopiaMapEditor/internal/fileio"
	"github.com/samuelyuan/PolytopiaMapEditor/internal/mapdraw"
)

type editor struct {
	drawSurface               *interactiveRaster
	status                    *widget.Label
	cache                     *image.RGBA
	cacheWidth, cacheHeight   int
	tileCoordinatesProperties *widget.Label
	terrainSelect             *widget.Select
	mapTileProperties         *widget.Label

	uri                   string
	img                   *image.RGBA
	mapData               *fileio.PolytopiaSaveOutput
	mapHeight             int
	mapWidth              int
	currentTileProperties string
	zoom                  int
	tileX                 int
	tileY                 int
	tileTerrain           int

	win        fyne.Window
	recentMenu *fyne.Menu
}

func (e *editor) PixelColor(x, y int) color.Color {
	return e.img.At(x, y)
}

func colorToBytes(col color.Color) []uint8 {
	r, g, b, a := col.RGBA()
	return []uint8{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

func (e *editor) Clicked(x, y int, edit api.Editor) {
	edit.SetHexCoordinates(x, y)
}

func (e *editor) SetHexCoordinates(x int, y int) {
	tileX, tileY := mapdraw.GetTileCoordinates(x, y, e.mapHeight, e.mapWidth)
	e.tileX = tileX
	e.tileY = tileY
	e.tileCoordinatesProperties.SetText(fmt.Sprintf("Tile (%d, %d)", e.tileX, e.tileY))

	e.currentTileProperties = ""

	tile := e.mapData.TileData[tileY][tileX]
	if tile.ImprovementData != nil && tile.ImprovementType == 1 {
		e.currentTileProperties += fmt.Sprintf("City: %s\n", tile.ImprovementData.CityName)
	}
	e.currentTileProperties += fmt.Sprintf("Owner: Player %d\n", tile.Owner)

	terrain := tile.Terrain
	e.tileTerrain = terrain
	e.terrainSelect.SetSelectedIndex(terrain - 1)

	resource := "Unknown"
	if tile.ResourceExists {
		resource = fmt.Sprintf("%v", tile.ResourceType)
	} else {
		resource = "No Resource"
	}
	e.currentTileProperties += fmt.Sprintf("Resource: %s\n", resource)

	improvement := "Unknown"
	if tile.ImprovementExists {
		improvement = fmt.Sprintf("%v", tile.ImprovementType)
	} else {
		improvement = "None"
	}
	e.currentTileProperties += fmt.Sprintf("Tile Improvement: %s\n", improvement)

	routeType := "Unknown"
	if tile.HasRoad {
		routeType = "Road"
	} else if tile.HasWaterRoute {
		routeType = "Water Route"
	} else {
		routeType = "No Route"
	}
	e.currentTileProperties += fmt.Sprintf("Route: %s\n", routeType)

	e.mapTileProperties.SetText(e.currentTileProperties)
}

func (e *editor) buildUI() fyne.CanvasObject {
	return container.NewScroll(e.drawSurface)
}

func (e *editor) setZoom(zoom int) {
	e.zoom = zoom
	e.updateSizes()
	e.drawSurface.Refresh()
}

func (e *editor) draw(w, h int) image.Image {
	if e.cacheWidth == 0 || e.cacheHeight == 0 {
		return image.NewRGBA(image.Rect(0, 0, w, h))
	}

	if w > e.cacheWidth || h > e.cacheHeight {
		bigger := image.NewRGBA(image.Rect(0, 0, w, h))
		draw.Draw(bigger, e.cache.Bounds(), e.cache, image.Point{}, draw.Over)
		return bigger
	}

	return e.cache
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

	img := mapdraw.DrawMap(saveFileData)

	e.addRecent(read.URI())
	e.uri = read.URI().String()
	e.img = fixEncoding(img)
	e.mapData = saveFileData
	e.mapHeight = len(saveFileData.TileData)
	e.mapWidth = len(saveFileData.TileData[0])
	e.status.SetText(fmt.Sprintf("File: %s | Map Rows: %d | Map Cols: %d",
		filepath.Base(read.URI().String()), e.mapHeight, e.mapWidth))
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

func (e *editor) WriteMapState() {
	originalFilname := e.uri[7:]
	decompressedFilename := originalFilname + ".decomp"
	outputFilename := originalFilname

	fileio.WriteMapToFile(decompressedFilename, e.mapData.TileData)
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

// NewEditor creates a new pixel editor that is ready to have a file loaded
func NewEditor() api.Editor {
	mapTileProperties := widget.NewLabel("Tile Properties")
	tileCoordinatesProperties := widget.NewLabel("Tile (-1, -1)")

	edit := &editor{
		zoom:                      1,
		mapTileProperties:         mapTileProperties,
		tileCoordinatesProperties: tileCoordinatesProperties,
		status:                    newStatusBar(),
		tileX:                     -1,
		tileY:                     -1,
		tileTerrain:               -1,
	}
	edit.drawSurface = newInteractiveRaster(edit)

	// terrainData := binding.BindInt(&edit.tileTerrain)
	terrainOptions := []string{"1 (Coast)", "2 (Ocean)", "3 (Field)", "4 (Mountain)", "5 (Forest)"}
	changeTerrainFunc := func(s string) {
		if edit.tileX == -1 || edit.tileY == -1 {
			return
		}
		optionInt, err := strconv.Atoi(strings.Split(s, " ")[0])
		if err != nil {
			log.Fatal(err)
		}
		edit.mapData.TileData[edit.tileY][edit.tileX].Terrain = optionInt

		altitude := 0
		switch optionInt {
		case 1:
			altitude = -1 // water altitude is -1
		case 2:
			altitude = -2 // ocean altitude is -2
		case 3:
		case 5:
			altitude = 1 // flat tile altitude is 1
		case 4:
			altitude = 2 // mountain altitude is 2
		}
		edit.mapData.TileData[edit.tileY][edit.tileX].Altitude = altitude

		img := mapdraw.DrawMap(edit.mapData)
		edit.img = fixEncoding(img)
		edit.updateSizes()
	}
	terrainSelect := widget.NewSelect(terrainOptions, changeTerrainFunc)
	edit.terrainSelect = terrainSelect
	return edit
}
