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
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
	tilePropertiesBar       fyne.CanvasObject
	playerPropertiesBar     fyne.CanvasObject
	cache                   *image.RGBA
	cacheWidth, cacheHeight int

	// Tile properties bar
	tileCoordinatesProperties *widget.Label
	terrainSelect             *widget.Select
	climateSelect             *widget.Select
	resourceSelect            *widget.Select
	improvementSelect         *widget.Select
	tileOwnerSelect           *widget.Select
	hasRoadCheckbox           *widget.Check
	cityNameEntry             *widget.Entry
	cityTileCoordinatesEntry  *widget.Entry
	unitOwnerSelect           *widget.Select
	unitTypeSelect            *widget.Select
	unitHasMovedCheckbox      *widget.Check
	unitHasAttackedCheckbox   *widget.Check
	mapTileProperties         *widget.Label

	// map properties
	uri                   string
	img                   *image.RGBA
	mapData               *fileio.PolytopiaSaveOutput
	mapHeight             int
	mapWidth              int
	currentTileProperties string
	zoom                  int

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
	e.tileCoordinatesProperties.SetText(fmt.Sprintf("Tile (%d, %d)", e.tileX, e.tileY))
	e.currentTileProperties = ""

	tile := e.mapData.TileData[tileY][tileX]
	e.tileOwnerSelect.SetSelectedIndex(tile.Owner)
	e.terrainSelect.SetSelectedIndex(tile.Terrain - 1)
	e.climateSelect.SetSelectedIndex(tile.Climate - 1)

	if tile.ResourceExists {
		e.resourceSelect.SetSelectedIndex(tile.ResourceType)
	} else {
		e.resourceSelect.SetSelectedIndex(0)
	}

	if tile.ImprovementExists {
		e.improvementSelect.SetSelectedIndex(tile.ImprovementType)
	} else {
		e.improvementSelect.SetSelectedIndex(0)
	}

	e.hasRoadCheckbox.SetChecked(tile.HasRoad)

	if tile.ImprovementData != nil && tile.ImprovementType == 1 {
		e.cityNameEntry.SetText(tile.ImprovementData.CityName)
		e.cityNameEntry.Enable()
	} else {
		e.cityNameEntry.SetText("")
		e.cityNameEntry.Disable()
	}

	e.cityTileCoordinatesEntry.SetText(fmt.Sprintf("(%v, %v)", tile.CapitalCoordinates[0], tile.CapitalCoordinates[1]))
	e.cityTileCoordinatesEntry.Disable()

	if tile.Unit != nil {
		e.unitOwnerSelect.SetSelectedIndex(int(tile.Unit.Owner))
		e.unitOwnerSelect.Enable()
		e.unitTypeSelect.SetSelectedIndex(int(tile.Unit.UnitType))
		e.unitTypeSelect.Enable()
		e.unitHasMovedCheckbox.SetChecked(tile.Unit.Moved)
		e.unitHasMovedCheckbox.Enable()
		e.unitHasAttackedCheckbox.SetChecked(tile.Unit.Attacked)
		e.unitHasAttackedCheckbox.Enable()
	} else {
		e.unitOwnerSelect.SetSelectedIndex(0)
		e.unitOwnerSelect.Disable()
		e.unitTypeSelect.SetSelectedIndex(0)
		e.unitTypeSelect.Disable()
		e.unitHasMovedCheckbox.SetChecked(false)
		e.unitHasMovedCheckbox.Disable()
		e.unitHasAttackedCheckbox.SetChecked(false)
		e.unitHasAttackedCheckbox.Disable()
	}

	e.mapTileProperties.SetText(e.currentTileProperties)
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
	e.tileOwnerSelect = createTileOwnerSelect(e)
	e.unitOwnerSelect = createUnitOwnerSelect(e)

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

func (edit *editor) RevealAllTilesForPlayer(newTribe int) {
	for i := 0; i < edit.mapData.MapHeight; i++ {
		for j := 0; j < edit.mapData.MapWidth; j++ {
			targetX := j
			targetY := i

			visibilityData := edit.mapData.TileData[i][j].PlayerVisibility
			fmt.Println("Existing visibility data:", visibilityData)
			isAlreadyVisible := false
			for visibilityIndex := 0; visibilityIndex < len(visibilityData); visibilityIndex++ {
				if int(visibilityData[visibilityIndex]) == newTribe {
					fmt.Printf("Tile is already visible to tribe %v. No change will be made to visibility data.\n", newTribe)
					isAlreadyVisible = true
					break
				}
			}
			if !isAlreadyVisible {
				edit.mapData.TileData[i][j].PlayerVisibility = append(edit.mapData.TileData[i][j].PlayerVisibility, newTribe)
				fmt.Println(fmt.Sprintf("Revealed (%v, %v) for tribe %v", targetX, targetY, newTribe))
			}
		}
	}

	for i := 0; i < edit.mapData.MapHeight; i++ {
		for j := 0; j < edit.mapData.MapWidth; j++ {
			fmt.Println(fmt.Sprintf("Tile (%v, %v) visibility: %v", j, i, edit.mapData.TileData[i][j].PlayerVisibility))
		}
	}
}

func (edit *editor) UnlockAllTechForPlayer(newTribe int) {
	// all tech in numerical order, but not the same order as in game
	allTech := []int{
		0,
		1,
		2,
		3,
		4,
		5,
		6,
		7,
		8,
		9,
		10,
		12,
		13,
		14,
		15,
		16,
		17,
		18,
		19,
		20,
		21,
		22,
		23,
		24,
		25, // Aquarion only, Free Diving
		26, // Aquarion only, Spearing
		27, // Aquarion only, Riding
		28, // Elyrion only, Forest Magic
		30, // Polaris only, Frostwork
		31, // Polaris only, Polar Warfare
		32, // Polaris only, Polarism
		33, // Cymanti only, Oceanology
		35, // Cymanti only, Shock Tactics
		36, // Cymanti only, Recycling
		37, // Cymanti only, Hydrology
		38,
		39,
		40, // Cymanti only, Fishing
		41, // Polaris only, Sledding
		42, // Polaris only, Ice Fishing
		43, // Cymanti only, Pescetism
	}
	for i := 0; i < len(edit.mapData.PlayerData); i++ {
		if edit.mapData.PlayerData[i].Id == newTribe {
			edit.mapData.PlayerData[i].AvailableTech = allTech
			break
		}
	}
}

func (e *editor) WriteMapState() {
	originalFilname := e.uri[7:]
	decompressedFilename := originalFilname + ".decomp"
	outputFilename := originalFilname

	fileio.WriteMapToFile(decompressedFilename, e.mapData.TileData)
	fileio.WritePlayersToFile(decompressedFilename, e.mapData.PlayerData)
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
	edit.tilePropertiesBar = newTilePropertiesBar(edit)
	edit.playerPropertiesBar = newPlayerPropertiesBar(edit)
	return fyne.NewContainerWithLayout(
		layout.NewBorderLayout(toolbar, edit.status, edit.tilePropertiesBar, edit.playerPropertiesBar),
		toolbar, edit.status, edit.tilePropertiesBar, edit.buildUI(), edit.playerPropertiesBar,
	)
}

func createTerrainSelect(edit *editor) *widget.Select {
	options := getTerrainOptions()
	changeFunc := func(s string) {
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

		img := mapdraw.DrawMap(edit.mapData, edit.tileX, edit.tileY)
		edit.img = fixEncoding(img)
		edit.updateSizes()
	}
	return widget.NewSelect(options, changeFunc)
}

func createClimateSelect(edit *editor) *widget.Select {
	options := getClimateOptions()
	changeFunc := func(s string) {
		if edit.tileX == -1 || edit.tileY == -1 {
			return
		}
		optionInt, err := strconv.Atoi(strings.Split(s, " ")[0])
		if err != nil {
			log.Fatal(err)
		}
		edit.mapData.TileData[edit.tileY][edit.tileX].Climate = optionInt

		img := mapdraw.DrawMap(edit.mapData, edit.tileX, edit.tileY)
		edit.img = fixEncoding(img)
		edit.updateSizes()
	}
	return widget.NewSelect(options, changeFunc)
}

func createResourceSelect(edit *editor) *widget.Select {
	options := getResourceSelectOptions()
	changeFunc := func(s string) {
		if edit.tileX == -1 || edit.tileY == -1 {
			return
		}

		if s == "None" {
			edit.mapData.TileData[edit.tileY][edit.tileX].ResourceExists = false
			edit.mapData.TileData[edit.tileY][edit.tileX].ResourceType = -1
			return
		}

		optionInt, err := strconv.Atoi(strings.Split(s, " ")[0])
		if err != nil {
			log.Fatal(err)
		}
		edit.mapData.TileData[edit.tileY][edit.tileX].ResourceExists = true
		edit.mapData.TileData[edit.tileY][edit.tileX].ResourceType = optionInt

		img := mapdraw.DrawMap(edit.mapData, edit.tileX, edit.tileY)
		edit.img = fixEncoding(img)
		edit.updateSizes()
	}
	return widget.NewSelect(options, changeFunc)
}

func createImprovementSelect(edit *editor) *widget.Select {
	options := getImprovementOptions()
	changeFunc := func(s string) {
		if edit.tileX == -1 || edit.tileY == -1 {
			return
		}

		if s == "None" {
			edit.mapData.TileData[edit.tileY][edit.tileX].ImprovementExists = false
			edit.mapData.TileData[edit.tileY][edit.tileX].ImprovementType = -1
			edit.mapData.TileData[edit.tileY][edit.tileX].ImprovementData = nil
			return
		}

		optionInt, err := strconv.Atoi(strings.Split(s, " ")[0])
		if err != nil {
			log.Fatal(err)
		}
		if edit.mapData.TileData[edit.tileY][edit.tileX].ImprovementData == nil {
			edit.mapData.TileData[edit.tileY][edit.tileX].ImprovementData = &fileio.ImprovementData{
				Level:                  1,
				FoundedTurn:            edit.mapData.MaxTurn,
				CurrentPopulation:      0,
				TotalPopulation:        0,
				Production:             1,
				BaseScore:              0,
				BorderSize:             1,
				UpgradeCount:           0,
				ConnectedPlayerCapital: 0,
				HasCityName:            0,
				CityName:               "",
				FoundedTribe:           edit.mapData.TileData[edit.tileY][edit.tileX].Owner,
				CityRewards:            []int{},
				RebellionFlag:          0,
				RebellionBuffer:        []int{},
			}
		}

		edit.mapData.TileData[edit.tileY][edit.tileX].ImprovementExists = true
		edit.mapData.TileData[edit.tileY][edit.tileX].ImprovementType = optionInt

		img := mapdraw.DrawMap(edit.mapData, edit.tileX, edit.tileY)
		edit.img = fixEncoding(img)
		edit.updateSizes()
	}
	return widget.NewSelect(options, changeFunc)
}

func createTileOwnerSelect(edit *editor) *widget.Select {
	options := []string{"None"}
	if edit.mapData != nil {
		for i := 0; i < len(edit.mapData.PlayerData); i++ {
			playerData := edit.mapData.PlayerData[i]
			if playerData.Id == 255 {
				continue
			}
			options = append(options, fmt.Sprintf("Player %d (%s)", playerData.Id, playerData.Name))
		}
	}

	changeFunc := func(s string) {
		if edit.tileX == -1 || edit.tileY == -1 {
			return
		}

		if s == "None" {
			edit.mapData.TileData[edit.tileY][edit.tileX].Owner = 0
			return
		}

		optionInt, err := strconv.Atoi(strings.Split(s, " ")[1])
		if err != nil {
			log.Fatal(err)
		}
		edit.mapData.TileData[edit.tileY][edit.tileX].Owner = optionInt

		img := mapdraw.DrawMap(edit.mapData, edit.tileX, edit.tileY)
		edit.img = fixEncoding(img)
		edit.updateSizes()
	}
	return widget.NewSelect(options, changeFunc)
}

func createHasRoadCheckbox(edit *editor) *widget.Check {
	return widget.NewCheck("", func(value bool) {
		if edit.tileX == -1 || edit.tileY == -1 {
			return
		}

		edit.mapData.TileData[edit.tileY][edit.tileX].HasRoad = value

		img := mapdraw.DrawMap(edit.mapData, edit.tileX, edit.tileY)
		edit.img = fixEncoding(img)
		edit.updateSizes()
	})
}

func createCityNameEntry(edit *editor) *widget.Entry {
	entry := widget.NewEntry()
	entry.OnChanged = func(value string) {
		if edit.tileX == -1 || edit.tileY == -1 {
			return
		}

		if edit.mapData.TileData[edit.tileY][edit.tileX].ImprovementData != nil &&
			edit.mapData.TileData[edit.tileY][edit.tileX].ImprovementType == 1 {
			if value == "" {
				edit.mapData.TileData[edit.tileY][edit.tileX].ImprovementData.HasCityName = 0
			} else {
				edit.mapData.TileData[edit.tileY][edit.tileX].ImprovementData.HasCityName = 1
			}
			edit.mapData.TileData[edit.tileY][edit.tileX].ImprovementData.CityName = value
		}

	}
	return entry
}

func createcityTileCoordinatesEntry(edit *editor) *widget.Entry {
	entry := widget.NewEntry()
	return entry
}

func createUnitOwnerSelect(edit *editor) *widget.Select {
	options := []string{"None"}
	if edit.mapData != nil {
		for i := 0; i < len(edit.mapData.PlayerData); i++ {
			playerData := edit.mapData.PlayerData[i]
			if playerData.Id == 255 {
				continue
			}
			options = append(options, fmt.Sprintf("Player %d (%s)", playerData.Id, playerData.Name))
		}
	}

	changeFunc := func(s string) {
		if edit.tileX == -1 || edit.tileY == -1 {
			return
		}

		if edit.mapData.TileData[edit.tileY][edit.tileX].Unit == nil {
			return
		}

		if s == "None" {
			edit.mapData.TileData[edit.tileY][edit.tileX].Unit.Owner = 0
			if edit.mapData.TileData[edit.tileY][edit.tileX].PassengerUnit != nil {
				edit.mapData.TileData[edit.tileY][edit.tileX].PassengerUnit.Owner = 0
			}
			return
		}

		optionInt, err := strconv.Atoi(strings.Split(s, " ")[1])
		if err != nil {
			log.Fatal(err)
		}
		edit.mapData.TileData[edit.tileY][edit.tileX].Unit.Owner = uint8(optionInt)
		if edit.mapData.TileData[edit.tileY][edit.tileX].PassengerUnit != nil {
			edit.mapData.TileData[edit.tileY][edit.tileX].PassengerUnit.Owner = uint8(optionInt)
		}

		img := mapdraw.DrawMap(edit.mapData, edit.tileX, edit.tileY)
		edit.img = fixEncoding(img)
		edit.updateSizes()
	}
	return widget.NewSelect(options, changeFunc)
}

func createUnitTypeSelect(edit *editor) *widget.Select {
	options := getUnitTypeOptions()

	changeFunc := func(s string) {
		if edit.tileX == -1 || edit.tileY == -1 {
			return
		}

		if edit.mapData.TileData[edit.tileY][edit.tileX].Unit == nil {
			return
		}

		if s == "None" {
			edit.mapData.TileData[edit.tileY][edit.tileX].Unit.UnitType = 0
			return
		}

		optionInt, err := strconv.Atoi(strings.Split(s, " ")[0])
		if err != nil {
			log.Fatal(err)
		}
		edit.mapData.TileData[edit.tileY][edit.tileX].Unit.UnitType = uint16(optionInt)

		img := mapdraw.DrawMap(edit.mapData, edit.tileX, edit.tileY)
		edit.img = fixEncoding(img)
		edit.updateSizes()
	}
	return widget.NewSelect(options, changeFunc)
}

func createUnitHasMovedCheckbox(edit *editor) *widget.Check {
	return widget.NewCheck("", func(value bool) {
		if edit.tileX == -1 || edit.tileY == -1 {
			return
		}

		if edit.mapData.TileData[edit.tileY][edit.tileX].Unit != nil {
			edit.mapData.TileData[edit.tileY][edit.tileX].Unit.Moved = value
		}
	})
}

func createUnitHasAttackedCheckbox(edit *editor) *widget.Check {
	return widget.NewCheck("", func(value bool) {
		if edit.tileX == -1 || edit.tileY == -1 {
			return
		}

		if edit.mapData.TileData[edit.tileY][edit.tileX].Unit != nil {
			edit.mapData.TileData[edit.tileY][edit.tileX].Unit.Attacked = value
		}
	})
}

func newTilePropertiesBar(edit *editor) fyne.CanvasObject {
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
		edit.tileCoordinatesProperties,
		container.NewHBox(widget.NewLabel("Terrain"), edit.terrainSelect),
		container.NewHBox(widget.NewLabel("Climate"), edit.climateSelect),
		container.NewHBox(widget.NewLabel("Resource"), edit.resourceSelect),
		container.NewHBox(widget.NewLabel("Improvement"), edit.improvementSelect),
		container.NewHBox(widget.NewLabel("Tile Owner"), edit.tileOwnerSelect),
		container.NewHBox(widget.NewLabel("Has Road"), edit.hasRoadCheckbox),
		container.NewHBox(widget.NewLabel("City Name"), edit.cityNameEntry),
		container.NewHBox(widget.NewLabel("City Tile Coordinates"), edit.cityTileCoordinatesEntry),
		container.NewHBox(widget.NewLabel("Unit Owner"), edit.unitOwnerSelect),
		container.NewHBox(widget.NewLabel("Unit Type"), edit.unitTypeSelect),
		container.NewHBox(widget.NewLabel("Unit Has Moved"), edit.unitHasMovedCheckbox),
		container.NewHBox(widget.NewLabel("Unit Has Attacked"), edit.unitHasAttackedCheckbox),
	})

	return container.NewVBox(options...)
}

func createPlayerNameEntry(edit *editor, playerIndex int, playerName string) *widget.Entry {
	entry := widget.NewEntry()
	entry.SetText(playerName)
	entry.OnChanged = func(value string) {
		edit.mapData.PlayerData[playerIndex].Name = value
	}
	return entry
}

func createPlayerTribeSelect(edit *editor, playerIndex int) *widget.Select {
	options := getPlayerTribeOptions()
	changeFunc := func(s string) {
		optionInt, err := strconv.Atoi(strings.Split(s, " ")[0])
		if err != nil {
			log.Fatal(err)
		}
		edit.mapData.PlayerData[playerIndex].Tribe = optionInt

		img := mapdraw.DrawMap(edit.mapData, edit.tileX, edit.tileY)
		edit.img = fixEncoding(img)
		edit.updateSizes()
	}
	return widget.NewSelect(options, changeFunc)
}

func createPlayerCurrencyEntry(edit *editor, playerIndex int, playerCurrency int) *widget.Entry {
	entry := widget.NewEntry()
	entry.SetText(fmt.Sprintf("%v", playerCurrency))
	entry.OnChanged = func(value string) {
		currencyInt, err := strconv.Atoi(value)
		if err != nil {
			entry.SetText(fmt.Sprintf("%v", playerCurrency))
			return
		}
		if currencyInt < 0 {
			currencyInt = 0
			entry.SetText(fmt.Sprintf("%v", currencyInt))
		} else if currencyInt > 2147483647 {
			// no overflow for max uint32
			currencyInt = 2147483647
			entry.SetText(fmt.Sprintf("%v", currencyInt))
		}
		edit.mapData.PlayerData[playerIndex].Currency = currencyInt
	}
	return entry
}

func createPlayerColorBox(edit *editor, playerIndex int, overrideColor []int) *canvas.Rectangle {
	playerColor := mapdraw.GetPlayerColor(edit.mapData.PlayerData[playerIndex])
	playerColorBox := canvas.NewRectangle(playerColor)
	playerColorBox.Resize(fyne.NewSize(1500, 1000))
	return playerColorBox
}

func createPlayerColorPickerButton(edit *editor, playerIndex int, overrideColor []int) *widget.Button {
	playerColor := mapdraw.GetPlayerColor(edit.mapData.PlayerData[playerIndex])
	colorPickerButton := widget.NewButton("Edit", func() {
		picker := dialog.NewColorPicker("Pick a Color", "Please pick tribe color:", func(newColor color.Color) {
			newColorR, newColorG, newColorB, _ := newColor.RGBA()
			edit.mapData.PlayerData[playerIndex].OverrideColor = []int{int(newColorB >> 8), int(newColorG >> 8), int(newColorR >> 8), 0}

			img := mapdraw.DrawMap(edit.mapData, edit.tileX, edit.tileY)
			edit.img = fixEncoding(img)
			edit.updateSizes()
		}, edit.win)
		picker.Advanced = true
		picker.SetColor(playerColor)
		picker.Show()
	})

	return colorPickerButton
}

func newPlayerPropertiesBar(edit *editor) fyne.CanvasObject {
	options := append([]fyne.CanvasObject{
		container.NewGridWithColumns(1),
		widget.NewLabel("Players"),
	})

	if edit.mapData != nil {
		for i := 0; i < len(edit.mapData.PlayerData); i++ {
			playerData := edit.mapData.PlayerData[i]
			if playerData.Id == 255 {
				continue
			}

			playerIdLabel := widget.NewLabel(fmt.Sprintf("Player %d", playerData.Id))
			playerNameEntry := createPlayerNameEntry(edit, i, playerData.Name)
			playerTribeSelect := createPlayerTribeSelect(edit, i)
			playerTribeSelect.SetSelectedIndex(playerData.Tribe - 2) // offset by 2 because options 0 and 1 aren't there
			playerStarsEntry := createPlayerCurrencyEntry(edit, i, playerData.Currency)
			playerColorBox := createPlayerColorBox(edit, i, playerData.OverrideColor)
			playerColorPickerButton := createPlayerColorPickerButton(edit, i, playerData.OverrideColor)

			options = append(options, []fyne.CanvasObject{
				container.NewHBox(playerIdLabel),
				playerNameEntry,
				playerTribeSelect,
				container.NewHBox(widget.NewLabel("Stars"), playerStarsEntry),
				container.NewHBox(widget.NewLabel("Color"), playerColorBox, playerColorPickerButton),
			}...)
		}
	}
	return container.NewVScroll(container.NewVBox(options...))
}

// NewEditor creates a new pixel editor that is ready to have a file loaded
func NewEditor() api.Editor {
	mapTileProperties := widget.NewLabel("Tile Properties")
	tileCoordinatesProperties := widget.NewLabel("Tile (-1, -1)")

	edit := &editor{
		zoom:                      1,
		zoomLabel:                 widget.NewLabel("100%"),
		mapTileProperties:         mapTileProperties,
		tileCoordinatesProperties: tileCoordinatesProperties,
		status:                    widget.NewLabel("Open a file"),
		tileX:                     -1,
		tileY:                     -1,
	}
	edit.drawSurface = newInteractiveRaster(edit)
	edit.terrainSelect = createTerrainSelect(edit)
	edit.climateSelect = createClimateSelect(edit)
	edit.resourceSelect = createResourceSelect(edit)
	edit.improvementSelect = createImprovementSelect(edit)
	edit.tileOwnerSelect = createTileOwnerSelect(edit)
	edit.hasRoadCheckbox = createHasRoadCheckbox(edit)
	edit.cityNameEntry = createCityNameEntry(edit)
	edit.cityTileCoordinatesEntry = createcityTileCoordinatesEntry(edit)
	edit.unitOwnerSelect = createUnitOwnerSelect(edit)
	edit.unitTypeSelect = createUnitTypeSelect(edit)
	edit.unitHasMovedCheckbox = createUnitHasMovedCheckbox(edit)
	edit.unitHasAttackedCheckbox = createUnitHasAttackedCheckbox(edit)

	return edit
}
