package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"

	"github.com/samuelyuan/PolytopiaMapEditor/internal/fileio"
)

const (
	DefaultPlayerId = 1
)

func (edit *editor) buildCheatsMenu() *fyne.Menu {
	return fyne.NewMenu("Cheats",
		fyne.NewMenuItem("Reveal All Tiles", edit.revealAllTiles),
		fyne.NewMenuItem("Unlock All Tech", edit.unlockAllTech),
		fyne.NewMenuItem("Complete All Tasks", edit.completeAllTasks),
		fyne.NewMenuItem("Convert All Units", edit.convertAllUnits),
	)
}

func (edit *editor) revealAllTiles() {
	win := fyne.CurrentApp().Driver().AllWindows()[0]
	dialog.ShowConfirm("Reveal all tiles?", "Are you sure you want to reveal all tiles?",
		func(ok bool) {
			if !ok {
				return
			}

			edit.RevealAllTilesForPlayer(DefaultPlayerId)
		}, win)
}

func (edit *editor) unlockAllTech() {
	win := fyne.CurrentApp().Driver().AllWindows()[0]
	dialog.ShowConfirm("Unlock all tech?", "Are you sure you want to unlock all tech?",
		func(ok bool) {
			if !ok {
				return
			}

			edit.UnlockAllTechForPlayer(DefaultPlayerId)
		}, win)
}

func (edit *editor) completeAllTasks() {
	win := fyne.CurrentApp().Driver().AllWindows()[0]
	dialog.ShowConfirm(
		"Complete all tasks?",
		"Are you sure you want to complete all tasks?",
		func(ok bool) {
			if !ok {
				return
			}

			edit.CompleteAllTasksForPlayer(DefaultPlayerId)
		}, win)
}

func (edit *editor) convertAllUnits() {
	win := fyne.CurrentApp().Driver().AllWindows()[0]
	dialog.ShowConfirm(
		"Convert all units to your tribe?",
		"Are you sure you want to convert all units? This will change all units to your tribe.",
		func(ok bool) {
			if !ok {
				return
			}

			edit.ConvertAllUnitsForPlayer(DefaultPlayerId)
		}, win)
}

func (edit *editor) RevealAllTilesForPlayer(newTribe int) {
	for i := 0; i < edit.mapData.MapHeight; i++ {
		for j := 0; j < edit.mapData.MapWidth; j++ {
			visibilityData := edit.mapData.TileData[i][j].PlayerVisibility
			isAlreadyVisible := false
			for visibilityIndex := 0; visibilityIndex < len(visibilityData); visibilityIndex++ {
				if int(visibilityData[visibilityIndex]) == newTribe {
					isAlreadyVisible = true
					break
				}
			}
			if !isAlreadyVisible {
				edit.mapData.TileData[i][j].PlayerVisibility = append(edit.mapData.TileData[i][j].PlayerVisibility, newTribe)
			}
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
		11, // Whaling
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

func (edit *editor) CompleteAllTasksForPlayer(newTribe int) {
	allTasksUnlocked := []fileio.PlayerTaskData{
		{
			Type:   1, // Pacifist
			Buffer: []int{1, 1, 5, 0, 0, 0},
		},
		{
			Type:   2, // Genius
			Buffer: []int{1, 1},
		},
		{
			Type:   3, // Network
			Buffer: []int{1, 1},
		},
		{
			Type:   4, // Wealth
			Buffer: []int{1, 1},
		},
		{
			Type:   5, // Killer
			Buffer: []int{1, 1, 10, 0, 0, 0},
		},
		{
			Type:   6, // Metropolis
			Buffer: []int{1, 1},
		},
		{
			Type:   8, // Explorer
			Buffer: []int{1, 1},
		},
	}

	for i := 0; i < len(edit.mapData.PlayerData); i++ {
		if edit.mapData.PlayerData[i].Id == newTribe {
			edit.mapData.PlayerData[i].Tasks = allTasksUnlocked
			break
		}
	}
}

func (edit *editor) ConvertAllUnitsForPlayer(newTribe int) {
	for i := 0; i < int(edit.mapData.MapHeight); i++ {
		for j := 0; j < int(edit.mapData.MapWidth); j++ {
			if edit.mapData.TileData[i][j].Unit == nil {
				continue
			}
			tribeOwner := edit.mapData.TileData[i][j].Owner
			if tribeOwner == newTribe {
				continue
			}
			edit.mapData.TileData[i][j].Unit.Owner = uint8(newTribe)
			if edit.mapData.TileData[i][j].PassengerUnit != nil {
				edit.mapData.TileData[i][j].PassengerUnit.Owner = uint8(newTribe)
			}
		}
	}

	edit.refreshMapImage()
}
