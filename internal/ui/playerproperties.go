package ui

import (
	"fmt"
	"image/color"
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/samuelyuan/PolytopiaMapEditor/internal/mapdraw"
)

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
		optionInt, err := parseLeadingInt(s)
		if err != nil {
			log.Fatal(err)
		}
		edit.mapData.PlayerData[playerIndex].Tribe = optionInt

		edit.refreshMapImage()
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
		if clamped := clampCurrency(currencyInt); clamped != currencyInt {
			currencyInt = clamped
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
			edit.mapData.PlayerData[playerIndex].OverrideColor = overrideColorFromRGBA(newColor)

			edit.refreshMapImage()
		}, edit.win)
		picker.Advanced = true
		picker.SetColor(playerColor)
		picker.Show()
	})

	return colorPickerButton
}

func newPlayerPropertiesBar(edit *editor) fyne.CanvasObject {
	options := []fyne.CanvasObject{
		container.NewGridWithColumns(1),
		widget.NewLabel("Players"),
	}

	if edit.mapData != nil {
		for i := 0; i < len(edit.mapData.PlayerData); i++ {
			playerData := edit.mapData.PlayerData[i]
			if playerData.PlayerId == 255 {
				continue
			}

			playerIdLabel := widget.NewLabel(fmt.Sprintf("Player %d", playerData.PlayerId))
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
	// NewScroll (not NewVScroll) so the panel can be narrowed via the split divider.
	return container.NewScroll(container.NewVBox(options...))
}
