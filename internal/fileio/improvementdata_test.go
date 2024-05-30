package fileio

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestDeserializeCityDataFromBytes(t *testing.T) {
	inputByteData := []byte{3, 0, 0, 0, 1, 0, 6, 0, 1, 0, 0, 0, 1, 0, 254, 255, 1, 1, 4, 84, 101, 115, 116, 0, 2, 0, 4, 0, 7, 0, 0, 0}
	streamReader := io.NewSectionReader(bytes.NewReader(inputByteData), 0, int64(len(inputByteData)))
	result := DeserializeImprovementDataFromBytes(streamReader)
	expected := ImprovementData{
		Level:                  3,
		CurrentPopulation:      1,
		TotalPopulation:        6,
		Production:             1,
		BaseScore:              0,
		BorderSize:             1,
		UpgradeCount:           -2,
		ConnectedPlayerCapital: 1,
		HasCityName:            1,
		CityName:               "Test",
		FoundedTribe:           0,
		CityRewards:            []int{4, 7},
		RebellionFlag:          0,
		RebellionBuffer:        []int{},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf(`Values not equal, result = %v, expected = %v`, result, expected)
	}
}

func TestDeserializeImprovementDataFromBytes(t *testing.T) {
	inputByteData := []byte{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	streamReader := io.NewSectionReader(bytes.NewReader(inputByteData), 0, int64(len(inputByteData)))
	result := DeserializeImprovementDataFromBytes(streamReader)
	expected := ImprovementData{
		Level:                  1,
		FoundedTurn:            0,
		CurrentPopulation:      0,
		TotalPopulation:        0,
		Production:             1,
		BaseScore:              0,
		BorderSize:             0,
		UpgradeCount:           0,
		ConnectedPlayerCapital: 0,
		HasCityName:            0,
		FoundedTribe:           0,
		CityRewards:            []int{},
		RebellionFlag:          0,
		RebellionBuffer:        []int{},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf(`result = %v, expected = %v`, result, expected)
	}
}

func TestSerializeCityDataToBytes(t *testing.T) {
	cityData := ImprovementData{
		Level:                  3,
		CurrentPopulation:      1,
		TotalPopulation:        6,
		Production:             1,
		BaseScore:              0,
		BorderSize:             1,
		UpgradeCount:           -2,
		ConnectedPlayerCapital: 1,
		HasCityName:            1,
		CityName:               "Test",
		FoundedTribe:           0,
		CityRewards:            []int{4, 7},
		RebellionFlag:          0,
		RebellionBuffer:        []int{},
	}
	resultBytes := SerializeImprovementDataToBytes(cityData)
	expectedBytes := []byte{3, 0, 0, 0, 1, 0, 6, 0, 1, 0, 0, 0, 1, 0, 254, 255, 1, 1, 4, 84, 101, 115, 116, 0, 2, 0, 4, 0, 7, 0, 0, 0}
	if !reflect.DeepEqual(resultBytes, expectedBytes) {
		t.Fatalf(`result = %v, expected = %v`, resultBytes, expectedBytes)
	}
}

func TestSerializeOtherImprovementDataToBytes(t *testing.T) {
	improvementData := ImprovementData{
		Level:                  1,
		FoundedTurn:            0,
		CurrentPopulation:      0,
		TotalPopulation:        0,
		Production:             1,
		BaseScore:              0,
		BorderSize:             0,
		UpgradeCount:           0,
		ConnectedPlayerCapital: 0,
		HasCityName:            0,
		FoundedTribe:           0,
		CityRewards:            []int{},
		RebellionFlag:          0,
		RebellionBuffer:        []int{},
	}
	resultBytes := SerializeImprovementDataToBytes(improvementData)
	expectedBytes := []byte{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	if !reflect.DeepEqual(resultBytes, expectedBytes) {
		t.Fatalf(`result = %v, expected = %v`, resultBytes, expectedBytes)
	}
}
