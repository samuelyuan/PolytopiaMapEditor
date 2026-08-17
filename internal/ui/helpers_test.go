package ui

import (
	"image/color"
	"testing"

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
