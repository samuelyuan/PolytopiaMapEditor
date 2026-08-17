package ui

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

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
