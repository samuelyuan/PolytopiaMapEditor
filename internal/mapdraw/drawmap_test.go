package mapdraw

import (
	"image"
	"image/color"
	"testing"

	"github.com/fogleman/gg"
	polytopiamapmodel "github.com/samuelyuan/polytopiamapmodelgo"
)

// imageIsBlank reports whether every pixel is fully transparent, i.e.
// nothing was drawn onto a fresh gg.Context.
func imageIsBlank(img image.Image) bool {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if _, _, _, a := img.At(x, y).RGBA(); a != 0 {
				return false
			}
		}
	}
	return true
}

func TestGetImagePosition(t *testing.T) {
	tests := []struct {
		i, j  int
		wantX float64
		wantY float64
	}{
		{0, 0, 0, 0},
		{1, 0, 0, 30},
		{0, 1, 30, 0},
		{2, 3, 90, 60},
	}

	for _, tt := range tests {
		x, y := getImagePosition(tt.i, tt.j)
		if x != tt.wantX || y != tt.wantY {
			t.Errorf("getImagePosition(%d, %d) = (%v, %v), want (%v, %v)", tt.i, tt.j, x, y, tt.wantX, tt.wantY)
		}
	}
}

func TestGetNeighbors(t *testing.T) {
	got := getNeighbors(5, 5)
	want := [4][2]int{{5, 6}, {4, 5}, {5, 4}, {6, 5}}
	if got != want {
		t.Errorf("getNeighbors(5, 5) = %v, want %v", got, want)
	}
}

func TestGetNeighborsWithDiagonals(t *testing.T) {
	got := getNeighborsWithDiagonals(5, 5)
	want := [8][2]int{{5, 6}, {6, 6}, {6, 5}, {6, 4}, {5, 4}, {4, 4}, {4, 5}, {4, 6}}
	if got != want {
		t.Errorf("getNeighborsWithDiagonals(5, 5) = %v, want %v", got, want)
	}
}

func TestGetPhysicalMapTileColor(t *testing.T) {
	tests := []struct {
		terrain int
		want    color.RGBA
	}{
		{1, color.RGBA{95, 149, 149, 255}},  // water
		{2, color.RGBA{47, 74, 93, 255}},    // ocean
		{3, color.RGBA{105, 125, 54, 255}},  // flat land
		{4, color.RGBA{105, 125, 54, 255}},  // mountain (flat land base color)
		{5, color.RGBA{105, 125, 54, 255}},  // forest (flat land base color)
		{6, color.RGBA{238, 249, 255, 255}}, // ice
		{0, color.RGBA{0, 0, 0, 255}},       // unknown/default
		{99, color.RGBA{0, 0, 0, 255}},      // unknown/default
	}

	for _, tt := range tests {
		got := getPhysicalMapTileColor(tt.terrain)
		if got != tt.want {
			t.Errorf("getPhysicalMapTileColor(%d) = %v, want %v", tt.terrain, got, tt.want)
		}
	}
}

func TestGetPlayerColor_KnownTribe(t *testing.T) {
	playerData := polytopiamapmodel.PlayerData{
		Tribe:         7, // Imperius
		OverrideColor: []int{0, 0, 0, 255},
	}
	got := GetPlayerColor(playerData)
	want := color.RGBA{0, 0, 255, 255}
	if got != want {
		t.Errorf("GetPlayerColor(Imperius) = %v, want %v", got, want)
	}
}

func TestGetPlayerColor_UnknownTribe(t *testing.T) {
	playerData := polytopiamapmodel.PlayerData{
		Tribe:         999,
		OverrideColor: []int{0, 0, 0, 255},
	}
	got := GetPlayerColor(playerData)
	want := color.RGBA{128, 128, 128, 255}
	if got != want {
		t.Errorf("GetPlayerColor(unknown tribe) = %v, want %v", got, want)
	}
}

func TestGetPlayerColor_OverrideColor(t *testing.T) {
	// OverrideColor is stored as [B, G, R, flag]; a flag != 255 means the
	// override is active and should take precedence over the tribe color.
	playerData := polytopiamapmodel.PlayerData{
		Tribe:         7, // Imperius (blue), should be ignored
		OverrideColor: []int{10, 20, 30, 0},
	}
	got := GetPlayerColor(playerData)
	want := color.RGBA{30, 20, 10, 255}
	if got != want {
		t.Errorf("GetPlayerColor(override) = %v, want %v", got, want)
	}
}

func TestLookupPlayerColor(t *testing.T) {
	playerColors := map[int]color.RGBA{
		1: {255, 0, 0, 255},
	}

	if got := lookupPlayerColor(playerColors, 1); got != (color.RGBA{255, 0, 0, 255}) {
		t.Errorf("lookupPlayerColor(known owner) = %v, want {255 0 0 255}", got)
	}

	want := color.RGBA{128, 128, 128, 255}
	if got := lookupPlayerColor(playerColors, 42); got != want {
		t.Errorf("lookupPlayerColor(unknown owner) = %v, want %v", got, want)
	}
}

func TestBuildPlayerColorMap(t *testing.T) {
	saveData := &polytopiamapmodel.PolytopiaSaveOutput{
		PlayerData: []polytopiamapmodel.PlayerData{
			{PlayerId: 1, Tribe: 7, OverrideColor: []int{0, 0, 0, 255}},  // Imperius
			{PlayerId: 2, Tribe: 10, OverrideColor: []int{0, 0, 0, 255}}, // Oumaji
		},
	}

	got := buildPlayerColorMap(saveData)

	if len(got) != 2 {
		t.Fatalf("buildPlayerColorMap() len = %d, want 2", len(got))
	}
	if got[1] != (color.RGBA{0, 0, 255, 255}) {
		t.Errorf("buildPlayerColorMap()[1] = %v, want Imperius blue", got[1])
	}
	if got[2] != (color.RGBA{255, 255, 0, 255}) {
		t.Errorf("buildPlayerColorMap()[2] = %v, want Oumaji yellow", got[2])
	}
}

// Below: the per-tile draw functions, now testable without a whole map.

func TestDrawSingleTerritoryTile_WaterFill(t *testing.T) {
	dc := gg.NewContext(int(radius), int(radius))
	saveData := &polytopiamapmodel.PolytopiaSaveOutput{
		TileData: [][]polytopiamapmodel.TileData{{{Terrain: 1}}}, // water
	}

	drawSingleTerritoryTile(dc, saveData, 0, 0, GraphicsOptions{}, nil)

	r, g, b, a := dc.Image().At(15, 15).RGBA()
	got := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
	want := color.RGBA{95, 149, 149, 255}
	if got != want {
		t.Errorf("center pixel = %v, want water color %v", got, want)
	}
}

func TestDrawSingleTerritoryTile_CapitalCityUsesPlayerColor(t *testing.T) {
	dc := gg.NewContext(int(radius), int(radius))
	saveData := &polytopiamapmodel.PolytopiaSaveOutput{
		TileData: [][]polytopiamapmodel.TileData{{{
			Terrain:         3,
			Owner:           1,
			ImprovementType: 1,
			ImprovementData: &polytopiamapmodel.ImprovementData{},
		}}},
	}
	playerColors := map[int]color.RGBA{1: {200, 10, 10, 255}}

	drawSingleTerritoryTile(dc, saveData, 0, 0, GraphicsOptions{}, playerColors)

	r, g, b, a := dc.Image().At(15, 15).RGBA()
	got := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
	want := color.RGBA{200, 10, 10, 255}
	if got != want {
		t.Errorf("city icon center pixel = %v, want owner color %v", got, want)
	}
}

func TestDrawRoadsForTile(t *testing.T) {
	t.Run("isolated road tile draws a stub", func(t *testing.T) {
		dc := gg.NewContext(int(radius), int(radius))
		saveData := &polytopiamapmodel.PolytopiaSaveOutput{
			TileData: [][]polytopiamapmodel.TileData{{{HasRoad: true}}},
		}
		drawRoadsForTile(dc, saveData, 1, 1, 0, 0)
		if imageIsBlank(dc.Image()) {
			t.Error("expected a road stub to be drawn for an isolated road tile")
		}
	})

	t.Run("no road and no city draws nothing", func(t *testing.T) {
		dc := gg.NewContext(int(radius), int(radius))
		saveData := &polytopiamapmodel.PolytopiaSaveOutput{
			TileData: [][]polytopiamapmodel.TileData{{{HasRoad: false}}},
		}
		drawRoadsForTile(dc, saveData, 1, 1, 0, 0)
		if !imageIsBlank(dc.Image()) {
			t.Error("expected nothing to be drawn for a tile with no road or city")
		}
	})
}

func TestDrawBordersForTile(t *testing.T) {
	playerColors := map[int]color.RGBA{1: {200, 10, 10, 255}}
	newSaveData := func(ownerA, ownerB int) *polytopiamapmodel.PolytopiaSaveOutput {
		return &polytopiamapmodel.PolytopiaSaveOutput{
			TileData: [][]polytopiamapmodel.TileData{{{Owner: ownerA}, {Owner: ownerB}}},
		}
	}

	t.Run("differing owner neighbor draws a border", func(t *testing.T) {
		dc := gg.NewContext(int(radius)*2, int(radius))
		drawBordersForTile(dc, newSaveData(1, 2), 1, 2, playerColors, 0, 0)
		if imageIsBlank(dc.Image()) {
			t.Error("expected a border to be drawn when the neighbor owner differs")
		}
	})

	t.Run("same owner neighbor draws nothing", func(t *testing.T) {
		dc := gg.NewContext(int(radius)*2, int(radius))
		drawBordersForTile(dc, newSaveData(1, 1), 1, 2, playerColors, 0, 0)
		if !imageIsBlank(dc.Image()) {
			t.Error("expected nothing to be drawn when the neighbor owner matches")
		}
	})

	t.Run("unowned tile draws nothing", func(t *testing.T) {
		dc := gg.NewContext(int(radius)*2, int(radius))
		drawBordersForTile(dc, newSaveData(0, 1), 1, 2, playerColors, 0, 0)
		if !imageIsBlank(dc.Image()) {
			t.Error("expected nothing to be drawn for an unowned tile")
		}
	})
}

func TestGetTileCoordinates(t *testing.T) {
	tests := []struct {
		name                string
		pixelX, pixelY      int
		mapWidth, mapHeight int
		wantX, wantY        int
	}{
		{"origin (top-left of inverted map)", 0, 0, 10, 10, 0, 9},
		{"one radius right", 30, 0, 10, 10, 1, 9},
		{"one radius down", 0, 30, 10, 10, 0, 8},
		{"mid-tile stays within same tile", 44, 44, 10, 10, 1, 8},
	}

	for _, tt := range tests {
		x, y := GetTileCoordinates(tt.pixelX, tt.pixelY, tt.mapWidth, tt.mapHeight)
		if x != tt.wantX || y != tt.wantY {
			t.Errorf("%s: GetTileCoordinates(%d, %d, %d, %d) = (%d, %d), want (%d, %d)",
				tt.name, tt.pixelX, tt.pixelY, tt.mapWidth, tt.mapHeight, x, y, tt.wantX, tt.wantY)
		}
	}
}
