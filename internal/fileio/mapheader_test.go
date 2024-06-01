package fileio

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

var (
	mapHeaderOutput = MapHeaderOutput{
		MapHeaderInput: MapHeaderInput{
			Version1:           0,
			Version2:           104,
			TotalActions:       300,
			CurrentTurn:        5,
			CurrentPlayerIndex: 0,
			MaxUnitId:          50,
			CurrentGameState:   2,
			Seed:               999999999,
			TurnLimit:          0,
			ScoreLimit:         0,
			WinByCapital:       0,
			UnknownSettings: [6]byte{
				0,
				1,
				0,
				1,
				0,
				0,
			},
			GameModeBase:  5,
			GameModeRules: 6,
		},
		MapName:           "Test Map",
		MapSquareSize:     16,
		DisabledTribesArr: []int{3, 5, 16, 17},
		UnlockedTribesArr: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
		GameDifficulty:    1,
		NumOpponents:      15,
		GameType:          0,
		MapPreset:         3,
		UnknownArr:        []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		SelectedTribeSkins: []TribeSkin{
			{
				Tribe: 11,
				Skin:  14,
			},
		},
		MapWidth:  16,
		MapHeight: 16,
	}

	mapHeaderBytes = []byte{0, 0, 0, 0, 104, 0, 0, 0, 44, 1, 5, 0, 0, 0, 0, 50, 0, 0, 0, 2,
		// seed
		255, 201, 154, 59,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 5, 6,
		// map name
		8, 84, 101, 115, 116, 32, 77, 97, 112,
		16, 0, 0, 0,
		// disabled tribes
		4, 0, 3, 0, 5, 0, 16, 0, 17, 0,
		// unlocked tribes
		18, 0, 0, 0, 1, 0, 2, 0, 3, 0, 4, 0, 5, 0, 6, 0, 7, 0, 8, 0, 9, 0, 10, 0, 11, 0, 12, 0, 13, 0, 14, 0, 15, 0, 16, 0, 17, 0,
		1, 0, 15, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 11, 0, 14, 0, 16, 0, 16, 0,
	}
)

func TestDeserializeMapHeaderFromBytes(t *testing.T) {
	inputByteData := mapHeaderBytes
	streamReader := io.NewSectionReader(bytes.NewReader(inputByteData), 0, int64(len(inputByteData)))
	result := DeserializeMapHeaderFromBytes(streamReader)
	expected := mapHeaderOutput

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf(`Contents not equal. Result = %v, expected = %v`, result, expected)
	}
}

func TestSerializeMapHeaderToBytes(t *testing.T) {
	resultBytes := SerializeMapHeaderToBytes(mapHeaderOutput)
	expectedBytes := mapHeaderBytes

	if !reflect.DeepEqual(resultBytes, expectedBytes) {
		t.Fatalf(`Contents not equal. Result = %v, expected = %v`, resultBytes, expectedBytes)
	}
}
