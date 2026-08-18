package ui

import (
	"image"
	"image/color"
	"testing"

	"fyne.io/fyne/v2"

	polytopiamapmodel "github.com/samuelyuan/polytopiamapmodelgo"
)

func TestClampZoom(t *testing.T) {
	tests := []struct {
		in   int
		want int
	}{
		{0, 1},
		{-5, 1},
		{1, 1},
		{8, 8},
		{16, 16},
		{17, 16},
		{1000, 16},
	}

	for _, tt := range tests {
		if got := clampZoom(tt.in); got != tt.want {
			t.Errorf("clampZoom(%d) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestClampCurrency(t *testing.T) {
	tests := []struct {
		in   int
		want int
	}{
		{-1, 0},
		{-100, 0},
		{0, 0},
		{100, 100},
		{2147483647, 2147483647},
		{2147483648, 2147483647},
		{1 << 40, 2147483647},
	}

	for _, tt := range tests {
		if got := clampCurrency(tt.in); got != tt.want {
			t.Errorf("clampCurrency(%d) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestTerrainToAltitude(t *testing.T) {
	tests := []struct {
		terrain int
		want    int
	}{
		{1, -1}, // water
		{2, -2}, // ocean
		{3, 0},  // field (no case body; preserves existing behavior)
		{4, 2},  // mountain
		{5, 1},  // forest
		{6, 0},  // ice, unhandled terrain
	}

	for _, tt := range tests {
		if got := terrainToAltitude(tt.terrain); got != tt.want {
			t.Errorf("terrainToAltitude(%d) = %d, want %d", tt.terrain, got, tt.want)
		}
	}
}

func TestParseLeadingInt(t *testing.T) {
	tests := []struct {
		in      string
		want    int
		wantErr bool
	}{
		{"3 (Field)", 3, false},
		{"12  (Lumber Hut)", 12, false},
		{"None", 0, true},
	}

	for _, tt := range tests {
		got, err := parseLeadingInt(tt.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("parseLeadingInt(%q) error = %v, wantErr %v", tt.in, err, tt.wantErr)
			continue
		}
		if err == nil && got != tt.want {
			t.Errorf("parseLeadingInt(%q) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestParsePlayerOptionId(t *testing.T) {
	tests := []struct {
		in      string
		want    int
		wantErr bool
	}{
		{"Player 2 (Bob)", 2, false},
		{"Player 10 (Alice)", 10, false},
		{"None", 0, true},
	}

	for _, tt := range tests {
		got, err := parsePlayerOptionId(tt.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("parsePlayerOptionId(%q) error = %v, wantErr %v", tt.in, err, tt.wantErr)
			continue
		}
		if err == nil && got != tt.want {
			t.Errorf("parsePlayerOptionId(%q) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestBuildPlayerOwnerOptions(t *testing.T) {
	players := []polytopiamapmodel.PlayerData{
		{PlayerId: 1, Name: "Alice"},
		{PlayerId: 255, Name: "Barbarians"}, // should be excluded
		{PlayerId: 2, Name: "Bob"},
	}

	got := buildPlayerOwnerOptions(players)
	want := []string{"None", "Player 1 (Alice)", "Player 2 (Bob)"}

	if len(got) != len(want) {
		t.Fatalf("buildPlayerOwnerOptions() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("buildPlayerOwnerOptions()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestBuildPlayerOwnerOptions_Empty(t *testing.T) {
	got := buildPlayerOwnerOptions(nil)
	if len(got) != 1 || got[0] != "None" {
		t.Errorf("buildPlayerOwnerOptions(nil) = %v, want [None]", got)
	}
}

func TestIsPNGPath(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"map.png", true},
		{"map.PNG", true},
		{"map.PnG", true},
		{"map.jpg", false},
		{"map.pngx", false},
		{"", false},
	}

	for _, tt := range tests {
		if got := isPNGPath(tt.path); got != tt.want {
			t.Errorf("isPNGPath(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestOverrideColorFromRGBA(t *testing.T) {
	got := overrideColorFromRGBA(color.RGBA{R: 30, G: 20, B: 10, A: 255})
	want := []int{10, 20, 30, 0}

	if len(got) != len(want) {
		t.Fatalf("overrideColorFromRGBA() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("overrideColorFromRGBA()[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestContainsInt(t *testing.T) {
	list := []int{1, 3, 5}

	if !containsInt(list, 3) {
		t.Error("containsInt(list, 3) = false, want true")
	}
	if containsInt(list, 4) {
		t.Error("containsInt(list, 4) = true, want false")
	}
	if containsInt(nil, 1) {
		t.Error("containsInt(nil, 1) = true, want false")
	}
}

func TestPushRecent(t *testing.T) {
	t.Run("prepends and caps", func(t *testing.T) {
		list := []string{"b", "c"}
		got := pushRecent(list, "a", 5)
		want := []string{"a", "b", "c"}

		if len(got) != len(want) {
			t.Fatalf("pushRecent() = %v, want %v", got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("pushRecent()[%d] = %q, want %q", i, got[i], want[i])
			}
		}
	})

	t.Run("caps at max", func(t *testing.T) {
		list := []string{"1", "2", "3", "4", "5"}
		got := pushRecent(list, "0", 5)
		want := []string{"0", "1", "2", "3", "4"}

		if len(got) != 5 {
			t.Fatalf("pushRecent() len = %d, want 5", len(got))
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("pushRecent()[%d] = %q, want %q", i, got[i], want[i])
			}
		}
	})

	t.Run("does not mutate input slice", func(t *testing.T) {
		list := []string{"b", "c"}
		_ = pushRecent(list, "a", 5)
		if list[0] != "b" || list[1] != "c" {
			t.Errorf("pushRecent() mutated input slice: %v", list)
		}
	})
}

func TestStripFileURIScheme(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"file:///C:/maps/foo.state", "/C:/maps/foo.state"},
		{"file://C:/maps/foo.state", "C:/maps/foo.state"},
		{"/no/scheme/here", "/no/scheme/here"},
		{"", ""},
	}

	for _, tt := range tests {
		if got := stripFileURIScheme(tt.in); got != tt.want {
			t.Errorf("stripFileURIScheme(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestBuildStatusText(t *testing.T) {
	got := buildStatusText("foo.state", 12, 34, 114)
	want := "File: foo.state | Map Rows: 12 | Map Cols: 34 | Game Version: 114"

	if got != want {
		t.Errorf("buildStatusText() = %q, want %q", got, want)
	}
}

func TestComputeDisplaySize(t *testing.T) {
	tests := []struct {
		name             string
		imgW, imgH, zoom int
		scale            float32
		want             fyne.Size
	}{
		{"1x zoom, 1x scale", 100, 200, 1, 1, fyne.NewSize(100, 200)},
		{"4x zoom, 1x scale", 100, 200, 4, 1, fyne.NewSize(400, 800)},
		{"1x zoom, 1.25x scale", 100, 200, 1, 1.25, fyne.NewSize(80, 160)},
		{"zero scale treated as 1x", 100, 200, 1, 0, fyne.NewSize(100, 200)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := computeDisplaySize(tt.imgW, tt.imgH, tt.zoom, tt.scale)
			if got != tt.want {
				t.Errorf("computeDisplaySize(%d, %d, %d, %v) = %v, want %v",
					tt.imgW, tt.imgH, tt.zoom, tt.scale, got, tt.want)
			}
		})
	}
}

func TestCanHaveRoad(t *testing.T) {
	tests := []struct {
		name string
		tile polytopiamapmodel.TileData
		want bool
	}{
		{"field", polytopiamapmodel.TileData{Terrain: 3}, true},
		{"forest", polytopiamapmodel.TileData{Terrain: 5}, true},
		{"water", polytopiamapmodel.TileData{Terrain: 1}, false},
		{"ocean", polytopiamapmodel.TileData{Terrain: 2}, false},
		{"mountain", polytopiamapmodel.TileData{Terrain: 4}, false},
		{
			"city on otherwise-valid terrain",
			polytopiamapmodel.TileData{Terrain: 3, ImprovementType: 1, ImprovementData: &polytopiamapmodel.ImprovementData{}},
			false,
		},
		{
			"non-city improvement on valid terrain",
			polytopiamapmodel.TileData{Terrain: 3, ImprovementType: 2, ImprovementData: &polytopiamapmodel.ImprovementData{}},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := canHaveRoad(tt.tile); got != tt.want {
				t.Errorf("canHaveRoad(%+v) = %v, want %v", tt.tile, got, tt.want)
			}
		})
	}
}

func TestFixEncoding(t *testing.T) {
	t.Run("already RGBA is returned as-is", func(t *testing.T) {
		src := image.NewRGBA(image.Rect(0, 0, 4, 4))
		src.Set(1, 1, color.RGBA{R: 200, G: 0, B: 0, A: 255})

		got := fixEncoding(src)
		if got != src {
			t.Error("fixEncoding() allocated a new image for an already-RGBA input, want the same pointer")
		}
	})

	t.Run("non-RGBA is converted", func(t *testing.T) {
		src := image.NewGray(image.Rect(0, 0, 2, 2))
		src.SetGray(0, 0, color.Gray{Y: 128})

		got := fixEncoding(src)
		if got.Bounds() != src.Bounds() {
			t.Errorf("fixEncoding() bounds = %v, want %v", got.Bounds(), src.Bounds())
		}
		r, g, b, _ := got.At(0, 0).RGBA()
		if r != g || g != b {
			t.Errorf("fixEncoding() did not preserve grayscale pixel: got RGB(%d,%d,%d)", r, g, b)
		}
	})
}

func TestBgPattern(t *testing.T) {
	// Same checkerboard cell should be the same shade; adjacent cells should differ.
	c00 := bgPattern(0, 0, 0, 0)
	c1010 := bgPattern(10, 10, 0, 0)
	c250 := bgPattern(25, 0, 0, 0)
	c025 := bgPattern(0, 25, 0, 0)

	if c00 != c1010 {
		t.Errorf("bgPattern(0,0) = %v, bgPattern(10,10) = %v; expected same cell shade", c00, c1010)
	}
	if c00 == c250 {
		t.Errorf("bgPattern(0,0) and bgPattern(25,0) should differ across a cell boundary")
	}
	if c00 == c025 {
		t.Errorf("bgPattern(0,0) and bgPattern(0,25) should differ across a cell boundary")
	}
}
