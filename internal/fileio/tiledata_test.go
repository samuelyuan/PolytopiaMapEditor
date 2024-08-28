package fileio

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestDeserializeEmptyTileDataFromBytes(t *testing.T) {
	inputByteData := []byte{3, 0, 0, 0, 1, 0, 0, 0, 3, 0, 8, 0, 1, 0, 0, 0,
		// coordinates
		255, 255, 255, 255, 255, 255, 255, 255,
		// resource
		0,
		// improvement
		0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}
	streamReader := io.NewSectionReader(bytes.NewReader(inputByteData), 0, int64(len(inputByteData)))
	result := DeserializeTileDataFromBytes(streamReader, 1, 3, 104)
	expected := TileData{
		WorldCoordinates:           [2]int{3, 1},
		Terrain:                    3,
		Climate:                    8,
		Altitude:                   1,
		Owner:                      0,
		Capital:                    0,
		CapitalCoordinates:         [2]int{-1, -1},
		ResourceExists:             false,
		ResourceType:               -1,
		ImprovementExists:          false,
		ImprovementType:            -1,
		ImprovementData:            (*ImprovementData)(nil),
		Unit:                       (*UnitData)(nil),
		PassengerUnit:              (*UnitData)(nil),
		UnitEffectData:             []int{},
		UnitDirectionData:          []int{},
		PassengerUnitEffectData:    []int{},
		PassengerUnitDirectionData: []int{},
		PlayerVisibility:           []int{},
		HasRoad:                    false,
		HasWaterRoute:              false,
		TileSkin:                   0,
		Unknown:                    []int{0, 0},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf(`Values not equal, result = %v, expected = %v`, result, expected)
	}
}

func TestSerializeEmptyTileDataToBytes(t *testing.T) {
	tileData := TileData{
		WorldCoordinates:           [2]int{3, 1},
		Terrain:                    3,
		Climate:                    8,
		Altitude:                   1,
		Owner:                      0,
		Capital:                    0,
		CapitalCoordinates:         [2]int{-1, -1},
		ResourceExists:             false,
		ResourceType:               -1,
		ImprovementExists:          false,
		ImprovementType:            -1,
		ImprovementData:            nil,
		Unit:                       nil,
		PassengerUnit:              nil,
		UnitEffectData:             []int{},
		UnitDirectionData:          []int{},
		PassengerUnitEffectData:    []int{},
		PassengerUnitDirectionData: []int{},
		PlayerVisibility:           []int{},
		HasRoad:                    false,
		HasWaterRoute:              false,
		TileSkin:                   0,
		Unknown:                    []int{0, 0},
	}
	versionNumber := 104
	resultBytes := SerializeTileToBytes(tileData, versionNumber)
	expectedBytes := []byte{3, 0, 0, 0, 1, 0, 0, 0, 3, 0, 8, 0, 1, 0, 0, 0,
		// coordinates
		255, 255, 255, 255, 255, 255, 255, 255,
		// resource
		0,
		// improvement
		0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}
	if !reflect.DeepEqual(resultBytes, expectedBytes) {
		t.Fatalf(`result = %v, expected = %v`, resultBytes, expectedBytes)
	}
}

func TestSerializeTileDataWithImprovementToBytes(t *testing.T) {
	tileData := TileData{
		WorldCoordinates:   [2]int{3, 1},
		Terrain:            3,
		Climate:            8,
		Altitude:           1,
		Owner:              0,
		Capital:            0,
		CapitalCoordinates: [2]int{-1, -1},
		ResourceExists:     false,
		ResourceType:       -1,
		ImprovementExists:  true,
		ImprovementType:    1,
		ImprovementData: &ImprovementData{
			Level:                  1,
			FoundedTurn:            0,
			CurrentPopulation:      0,
			TotalPopulation:        0,
			Production:             1,
			BaseScore:              0,
			BorderSize:             1,
			UpgradeCount:           0,
			ConnectedPlayerCapital: 0,
			HasCityName:            0,
			FoundedTribe:           0,
			CityRewards:            []int{},
			RebellionFlag:          0,
			RebellionBuffer:        []int{},
		},
		Unit:                       nil,
		UnitEffectData:             []int{},
		UnitDirectionData:          []int{},
		PassengerUnitEffectData:    []int{},
		PassengerUnitDirectionData: []int{},
		PlayerVisibility:           []int{},
		HasRoad:                    false,
		HasWaterRoute:              false,
		TileSkin:                   0,
		Unknown:                    []int{0, 0},
	}

	versionNumber := 104
	resultBytes := SerializeTileToBytes(tileData, versionNumber)
	expectedBytes := []byte{3, 0, 0, 0, 1, 0, 0, 0, 3, 0, 8, 0, 1, 0, 0, 0,
		// coordinates
		255, 255, 255, 255, 255, 255, 255, 255,
		// resource
		0,
		// improvement
		1, 1, 0,
		1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	if !reflect.DeepEqual(resultBytes, expectedBytes) {
		t.Fatalf(`result = %v, expected = %v`, resultBytes, expectedBytes)
	}
}

func TestSerializeTileDataWithUnitToBytes(t *testing.T) {
	tileData := TileData{
		WorldCoordinates:   [2]int{3, 1},
		Terrain:            3,
		Climate:            8,
		Altitude:           1,
		Owner:              0,
		Capital:            0,
		CapitalCoordinates: [2]int{-1, -1},
		ResourceExists:     false,
		ResourceType:       -1,
		ImprovementExists:  false,
		ImprovementType:    0,
		ImprovementData:    nil,
		Unit: &UnitData{
			Id:                 4,
			Owner:              4,
			UnitType:           2,
			FollowerUnitId:     0,
			LeaderUnitId:       0,
			CurrentCoordinates: [2]int32{3, 1},
			HomeCoordinates:    [2]int32{3, 1},
			Health:             100,
			PromotionLevel:     0,
			Experience:         0,
			Moved:              false,
			Attacked:           false,
			Flipped:            false,
			CreatedTurn:        0,
		},
		UnitEffectData:             []int{1},
		UnitDirectionData:          []int{255, 255, 1, 0, 0},
		PassengerUnitEffectData:    []int{},
		PassengerUnitDirectionData: []int{},
		PlayerVisibility:           []int{},
		HasRoad:                    false,
		HasWaterRoute:              false,
		TileSkin:                   0,
		Unknown:                    []int{0, 0},
	}

	versionNumber := 104
	resultBytes := SerializeTileToBytes(tileData, versionNumber)
	expectedBytes := []byte{3, 0, 0, 0, 1, 0, 0, 0, 3, 0, 8, 0, 1, 0, 0, 0,
		// coordinates
		255, 255, 255, 255, 255, 255, 255, 255,
		// resource
		0,
		// improvement
		0,
		1, 4, 0, 0, 0, 4, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 1, 0, 0, 0, 3, 0, 0, 0, 1, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 255, 255, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	compareArrays(t, resultBytes, expectedBytes)
}

func TestSerializeUnitDataToBytes(t *testing.T) {
	unitData := UnitData{
		Id:                 4,
		Owner:              4,
		UnitType:           2,
		FollowerUnitId:     0,
		LeaderUnitId:       0,
		CurrentCoordinates: [2]int32{8, 2},
		HomeCoordinates:    [2]int32{8, 2},
		Health:             100,
		PromotionLevel:     0,
		Experience:         0,
		Moved:              false,
		Attacked:           false,
		Flipped:            false,
		CreatedTurn:        0,
	}
	resultBytes := SerializeUnitDataToBytes(unitData)
	expectedBytes := []byte{4, 0, 0, 0, 4, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 0, 0, 0, 2, 0, 0, 0, 8, 0, 0, 0, 2, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	if !reflect.DeepEqual(resultBytes, expectedBytes) {
		t.Fatalf(`result = %v, expected = %v`, resultBytes, expectedBytes)
	}
}
