package ui

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/samuelyuan/PolytopiaMapEditor/internal/fileio"
)

type TileProperties struct {
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
}

func NewTileProperties(edit *editor) TileProperties {
	tileProperties := TileProperties{}
	tileProperties.tileCoordinatesProperties = widget.NewLabel("Tile (-1, -1)")
	tileProperties.terrainSelect = createTerrainSelect(edit)
	tileProperties.climateSelect = createClimateSelect(edit)
	tileProperties.resourceSelect = createResourceSelect(edit)
	tileProperties.improvementSelect = createImprovementSelect(edit)
	tileProperties.tileOwnerSelect = createTileOwnerSelect(edit)
	tileProperties.hasRoadCheckbox = createHasRoadCheckbox(edit)
	tileProperties.cityNameEntry = createCityNameEntry(edit)
	tileProperties.cityTileCoordinatesEntry = createCityTileCoordinatesEntry(edit)
	tileProperties.unitOwnerSelect = createUnitOwnerSelect(edit)
	tileProperties.unitTypeSelect = createUnitTypeSelect(edit)
	tileProperties.unitHasMovedCheckbox = createUnitHasMovedCheckbox(edit)
	tileProperties.unitHasAttackedCheckbox = createUnitHasAttackedCheckbox(edit)
	return tileProperties
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

		edit.refreshMapImage()
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

		edit.refreshMapImage()
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

		edit.refreshMapImage()
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

		edit.refreshMapImage()
	}
	return widget.NewSelect(options, changeFunc)
}

func createTileOwnerSelect(edit *editor) *widget.Select {
	options := []string{"None"}
	if edit.mapData != nil {
		for i := 0; i < len(edit.mapData.PlayerData); i++ {
			playerData := edit.mapData.PlayerData[i]
			if playerData.PlayerId == 255 {
				continue
			}
			options = append(options, fmt.Sprintf("Player %d (%s)", playerData.PlayerId, playerData.Name))
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

		edit.refreshMapImage()
	}
	return widget.NewSelect(options, changeFunc)
}

func createHasRoadCheckbox(edit *editor) *widget.Check {
	return widget.NewCheck("", func(value bool) {
		if edit.tileX == -1 || edit.tileY == -1 {
			return
		}

		edit.mapData.TileData[edit.tileY][edit.tileX].HasRoad = value

		edit.refreshMapImage()
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

func createCityTileCoordinatesEntry(edit *editor) *widget.Entry {
	entry := widget.NewEntry()
	return entry
}

func createUnitOwnerSelect(edit *editor) *widget.Select {
	options := []string{"None"}
	if edit.mapData != nil {
		for i := 0; i < len(edit.mapData.PlayerData); i++ {
			playerData := edit.mapData.PlayerData[i]
			if playerData.PlayerId == 255 {
				continue
			}
			options = append(options, fmt.Sprintf("Player %d (%s)", playerData.PlayerId, playerData.Name))
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

		edit.refreshMapImage()
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

		edit.refreshMapImage()
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

func (tileProperties *TileProperties) UpdateTileProperties(tileX int, tileY int, tile fileio.TileData) {
	tileProperties.tileCoordinatesProperties.SetText(fmt.Sprintf("Tile (%d, %d)", tileX, tileY))

	tileProperties.tileOwnerSelect.SetSelectedIndex(tile.Owner)
	tileProperties.terrainSelect.SetSelectedIndex(tile.Terrain - 1)
	tileProperties.climateSelect.SetSelectedIndex(tile.Climate - 1)

	if tile.ResourceExists {
		tileProperties.resourceSelect.SetSelectedIndex(tile.ResourceType)
	} else {
		tileProperties.resourceSelect.SetSelectedIndex(0)
	}

	if tile.ImprovementExists {
		tileProperties.improvementSelect.SetSelectedIndex(tile.ImprovementType)
	} else {
		tileProperties.improvementSelect.SetSelectedIndex(0)
	}

	tileProperties.hasRoadCheckbox.SetChecked(tile.HasRoad)
	// Road shouldn't be allowed for water or mountains or cities
	if tile.Terrain == 1 || tile.Terrain == 2 || tile.Terrain == 4 || (tile.ImprovementData != nil && tile.ImprovementType == 1) {
		tileProperties.hasRoadCheckbox.Disable()
	} else {
		tileProperties.hasRoadCheckbox.Enable()
	}

	if tile.ImprovementData != nil && tile.ImprovementType == 1 {
		tileProperties.cityNameEntry.SetText(tile.ImprovementData.CityName)
		tileProperties.cityNameEntry.Enable()
	} else {
		tileProperties.cityNameEntry.SetText("")
		tileProperties.cityNameEntry.Disable()
	}

	tileProperties.cityTileCoordinatesEntry.SetText(fmt.Sprintf("(%v, %v)", tile.CapitalCoordinates[0], tile.CapitalCoordinates[1]))
	tileProperties.cityTileCoordinatesEntry.Disable()

	if tile.Unit != nil {
		tileProperties.unitOwnerSelect.SetSelectedIndex(int(tile.Unit.Owner))
		tileProperties.unitOwnerSelect.Enable()
		tileProperties.unitTypeSelect.SetSelectedIndex(int(tile.Unit.UnitType))
		tileProperties.unitTypeSelect.Enable()
		tileProperties.unitHasMovedCheckbox.SetChecked(tile.Unit.Moved)
		tileProperties.unitHasMovedCheckbox.Enable()
		tileProperties.unitHasAttackedCheckbox.SetChecked(tile.Unit.Attacked)
		tileProperties.unitHasAttackedCheckbox.Enable()
	} else {
		tileProperties.unitOwnerSelect.SetSelectedIndex(0)
		tileProperties.unitOwnerSelect.Disable()
		tileProperties.unitTypeSelect.SetSelectedIndex(0)
		tileProperties.unitTypeSelect.Disable()
		tileProperties.unitHasMovedCheckbox.SetChecked(false)
		tileProperties.unitHasMovedCheckbox.Disable()
		tileProperties.unitHasAttackedCheckbox.SetChecked(false)
		tileProperties.unitHasAttackedCheckbox.Disable()
	}
}

func (tileProperties *TileProperties) GetOptions() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		tileProperties.tileCoordinatesProperties,
		container.NewHBox(widget.NewLabel("Terrain"), tileProperties.terrainSelect),
		container.NewHBox(widget.NewLabel("Climate"), tileProperties.climateSelect),
		container.NewHBox(widget.NewLabel("Resource"), tileProperties.resourceSelect),
		container.NewHBox(widget.NewLabel("Improvement"), tileProperties.improvementSelect),
		container.NewHBox(widget.NewLabel("Tile Owner"), tileProperties.tileOwnerSelect),
		container.NewHBox(widget.NewLabel("Has Road"), tileProperties.hasRoadCheckbox),
		container.NewHBox(widget.NewLabel("City Name"), tileProperties.cityNameEntry),
		container.NewHBox(widget.NewLabel("City Tile Coordinates"), tileProperties.cityTileCoordinatesEntry),
		container.NewHBox(widget.NewLabel("Unit Owner"), tileProperties.unitOwnerSelect),
		container.NewHBox(widget.NewLabel("Unit Type"), tileProperties.unitTypeSelect),
		container.NewHBox(widget.NewLabel("Unit Has Moved"), tileProperties.unitHasMovedCheckbox),
		container.NewHBox(widget.NewLabel("Unit Has Attacked"), tileProperties.unitHasAttackedCheckbox),
	}
}

func (tileProperties *TileProperties) UpdateOwnerOptions(edit *editor) {
	tileProperties.tileOwnerSelect = createTileOwnerSelect(edit)
	tileProperties.unitOwnerSelect = createUnitOwnerSelect(edit)
}
