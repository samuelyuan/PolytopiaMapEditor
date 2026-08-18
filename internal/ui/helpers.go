package ui

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
	"strings"

	"golang.org/x/image/draw"

	"fyne.io/fyne/v2"

	polytopiamapmodel "github.com/samuelyuan/polytopiamapmodelgo"
)

// clampZoom keeps the zoom level within the [1, 16] range supported by the toolbar.
func clampZoom(val int) int {
	if val < 1 {
		return 1
	}
	if val > 16 {
		return 16
	}
	return val
}

// clampCurrency keeps a player's currency within [0, math.MaxInt32], the range
// the save format can represent.
func clampCurrency(v int) int {
	if v < 0 {
		return 0
	}
	if v > 2147483647 {
		// no overflow for max uint32
		return 2147483647
	}
	return v
}

// terrainToAltitude returns the tile altitude implied by a terrain selection.
func terrainToAltitude(terrain int) int {
	altitude := 0
	switch terrain {
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
	return altitude
}

// parseLeadingInt parses the leading integer token of a select option string
// such as "3 (Field)" or "None".
func parseLeadingInt(s string) (int, error) {
	return strconv.Atoi(strings.Split(s, " ")[0])
}

// parsePlayerOptionId parses the player id out of an option string produced by
// buildPlayerOwnerOptions, e.g. "Player 2 (Bob)" -> 2.
func parsePlayerOptionId(s string) (int, error) {
	fields := strings.Split(s, " ")
	if len(fields) < 2 {
		return 0, fmt.Errorf("unexpected player option format: %q", s)
	}
	return strconv.Atoi(fields[1])
}

// buildPlayerOwnerOptions builds the "None"-prefixed list of selectable owner
// labels shared by the tile owner and unit owner dropdowns.
func buildPlayerOwnerOptions(players []polytopiamapmodel.PlayerData) []string {
	options := []string{"None"}
	for _, playerData := range players {
		if playerData.PlayerId == 255 {
			continue
		}
		options = append(options, fmt.Sprintf("Player %d (%s)", playerData.PlayerId, playerData.Name))
	}
	return options
}

// isPNGPath reports whether path ends with a "png" extension, case-insensitively.
func isPNGPath(path string) bool {
	return strings.LastIndex(strings.ToLower(path), "png") == len(path)-3
}

// overrideColorFromRGBA converts a picked color into the [B, G, R, flag] triple
// the save format stores as a player's OverrideColor, with the override flag enabled.
func overrideColorFromRGBA(newColor color.Color) []int {
	r, g, b, _ := newColor.RGBA()
	return []int{int(b >> 8), int(g >> 8), int(r >> 8), 0}
}

// containsInt reports whether val is present in list.
func containsInt(list []int, val int) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}

// pushRecent prepends item to list and caps the result at max entries, used to
// maintain the "recently opened" file list.
func pushRecent[T any](list []T, item T, max int) []T {
	updated := append([]T{item}, list...)
	if len(updated) > max {
		updated = updated[:max]
	}
	return updated
}

const fileURIScheme = "file://"

// stripFileURIScheme strips a leading "file://" from a Fyne file URI string,
// returning the underlying filesystem path. Returns uri unchanged if it
// doesn't have that prefix.
func stripFileURIScheme(uri string) string {
	if strings.HasPrefix(uri, fileURIScheme) {
		return uri[len(fileURIScheme):]
	}
	return uri
}

// buildStatusText formats the status bar text shown after a map loads.
func buildStatusText(filename string, mapHeight, mapWidth, gameVersion int) string {
	return fmt.Sprintf("File: %s | Map Rows: %d | Map Cols: %d | Game Version: %d",
		filename, mapHeight, mapWidth, gameVersion)
}

// computeDisplaySize returns the widget size, in Fyne logical units, needed
// to show an imgW x imgH image at the given zoom level such that 1 image
// pixel maps to 1 physical screen pixel at zoom=1, regardless of the
// canvas's DPI scale factor.
func computeDisplaySize(imgW, imgH, zoom int, scale float32) fyne.Size {
	if scale <= 0 {
		scale = 1
	}
	return fyne.NewSize(
		float32(imgW*zoom)/scale,
		float32(imgH*zoom)/scale,
	)
}

// canHaveRoad reports whether roads are allowed on this tile: not water,
// ocean, or mountain, and not a city.
func canHaveRoad(tile polytopiamapmodel.TileData) bool {
	if tile.Terrain == 1 || tile.Terrain == 2 || tile.Terrain == 4 {
		return false
	}
	if tile.ImprovementData != nil && tile.ImprovementType == 1 {
		return false
	}
	return true
}

// fixEncoding returns img as *image.RGBA, converting it if it isn't already
// one (canvas.Image and PNG encoding both need a concrete RGBA buffer).
func fixEncoding(img image.Image) *image.RGBA {
	if rgba, ok := img.(*image.RGBA); ok {
		return rgba
	}

	newImg := image.NewRGBA(img.Bounds())
	draw.Draw(newImg, newImg.Bounds(), img, img.Bounds().Min, draw.Over)
	return newImg
}
