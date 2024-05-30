package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

const (
	DefaultPlayerId = 1
)

func (edit *editor) buildCheatsMenu() *fyne.Menu {
	return fyne.NewMenu("Cheats",
		fyne.NewMenuItem("Reveal All Tiles", edit.revealAllTiles),
		fyne.NewMenuItem("Unlock All Tech", edit.unlockAllTech),
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
