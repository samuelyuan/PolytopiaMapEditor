package fileio

import (
	"bytes"
	"image/color"
	"io"
	"reflect"
	"testing"
)

var (
	playerData = PlayerData{
		PlayerId:             1,
		Name:                 "TestPlayer",
		AccountId:            "00000000-0000-0000-0000-000000000000",
		AutoPlay:             true,
		StartTileCoordinates: [2]int{6, 22},
		Tribe:                15,
		UnknownByte1:         1,
		DifficultyHandicap:   1,
		AggressionsByPlayers: []PlayerAggression{
			{PlayerId: 1, Aggression: 0},
			{PlayerId: 2, Aggression: 17744},
			{PlayerId: 3, Aggression: 73048},
			{PlayerId: 4, Aggression: 24359},
			{PlayerId: 5, Aggression: 74462},
			{PlayerId: 6, Aggression: 85466},
			{PlayerId: 7, Aggression: 64134},
			{PlayerId: 8, Aggression: 39411},
			{PlayerId: 9, Aggression: 36739},
			{PlayerId: 10, Aggression: 37812},
			{PlayerId: 11, Aggression: 22858},
			{PlayerId: 12, Aggression: 32007},
			{PlayerId: 13, Aggression: 17738},
			{PlayerId: 14, Aggression: 41794},
			{PlayerId: 15, Aggression: 55461},
			{PlayerId: 16, Aggression: 32041},
			{PlayerId: 255, Aggression: 0},
		},
		Currency:           900,
		Score:              10000,
		UnknownInt2:        0,
		NumCities:          11,
		AvailableTech:      []int{0, 8, 15, 10, 39, 18, 13, 1, 4, 14, 20},
		EncounteredPlayers: []int{7, 11, 3, 5, 10},
		Tasks: []PlayerTaskData{
			{Type: 6, Buffer: []int{1, 1}},
			{Type: 5, Buffer: []int{1, 1, 10, 0, 0, 0}},
			{Type: 8, Buffer: []int{1, 0}},
			{Type: 3, Buffer: []int{1, 1}},
		},
		TotalUnitsKilled:     28,
		TotalUnitsLost:       32,
		TotalTribesDestroyed: 1,
		OverrideColor:        []int{153, 0, 255, 255},
		OverrideTribe:        0,
		UniqueImprovements:   []int{27},
		DiplomacyArr: []DiplomacyData{
			{PlayerId: 1, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 2, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 3, DiplomacyRelationState: 0, LastAttackTurn: 21, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: 7, EmbassyBuildTurn: -100, PreviousAttackTurn: 21},
			{PlayerId: 4, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 5, DiplomacyRelationState: 0, LastAttackTurn: 19, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: 8, EmbassyBuildTurn: -100, PreviousAttackTurn: 21},
			{PlayerId: 6, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 7, DiplomacyRelationState: 0, LastAttackTurn: 20, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: 0, EmbassyBuildTurn: -100, PreviousAttackTurn: 21},
			{PlayerId: 8, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 9, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 10, DiplomacyRelationState: 0, LastAttackTurn: 15, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: 13, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 11, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: 0, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 12, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 13, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 14, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 15, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 16, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 255, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 0, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
		},
		DiplomacyMessages: []DiplomacyMessage{},
		DestroyedByTribe:  0,
		DestroyedTurn:     0,
		UnknownBuffer2:    []int{255, 255, 255, 255},
		EndScore:          -1,
		PlayerSkin:        0,
		UnknownBuffer3:    []int{255, 255, 255, 255},
	}

	playerBytes = []byte{1,
		// Player name
		10, 84, 101, 115, 116, 80, 108, 97, 121, 101, 114,
		// Account Id
		36, 48, 48, 48, 48, 48, 48, 48, 48, 45, 48, 48, 48, 48, 45, 48, 48, 48, 48, 45, 48, 48, 48, 48, 45, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48,
		1, 6, 0, 0, 0, 22, 0, 0, 0, 15, 0, 1, 1, 0, 0, 0,
		// Unknown Array 1
		17, 0,
		1, 0, 0, 0, 0, 2, 80, 69, 0, 0, 3, 88, 29, 1, 0, 4, 39, 95, 0, 0, 5, 222, 34, 1, 0,
		6, 218, 77, 1, 0, 7, 134, 250, 0, 0, 8, 243, 153, 0, 0, 9, 131, 143, 0, 0, 10, 180, 147, 0, 0,
		11, 74, 89, 0, 0, 12, 7, 125, 0, 0, 13, 74, 69, 0, 0, 14, 66, 163, 0, 0, 15, 165, 216, 0, 0,
		16, 41, 125, 0, 0, 255, 0, 0, 0, 0,
		// currency
		132, 3, 0, 0,
		// score
		16, 39, 0, 0,
		0, 0, 0, 0,
		// num cities
		11, 0,
		// tech
		11, 0, 0, 0, 8, 0, 15, 0, 10, 0, 39, 0, 18, 0, 13, 0, 1, 0, 4, 0, 14, 0, 20, 0,
		// encountered players
		5, 0, 7, 11, 3, 5, 10,
		// tasks
		4, 0, 6, 0, 1, 1, 5, 0, 1, 1, 10, 0, 0, 0, 8, 0, 1, 0, 3, 0, 1, 1,
		28, 0, 0, 0,
		32, 0, 0, 0,
		1, 0, 0, 0,
		// override color
		153, 0, 255, 255,
		0,
		// improvements
		1, 0, 27, 0,
		18, 0, 1, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 2, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 3, 0, 21, 0, 0, 0, 0, 156, 255, 255, 255, 7, 0, 0, 0, 156, 255, 255, 255, 21, 0, 0, 0, 4, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 5, 0, 19, 0, 0, 0, 0, 156, 255, 255, 255, 8, 0, 0, 0, 156, 255, 255, 255, 21, 0, 0, 0, 6, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 7, 0, 20, 0, 0, 0, 0, 156, 255, 255, 255, 0, 0, 0, 0, 156, 255, 255, 255, 21, 0, 0, 0, 8, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 9, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 10, 0, 15, 0, 0, 0, 0, 156, 255, 255, 255, 13, 0, 0, 0, 156, 255, 255, 255, 156, 255, 255, 255, 11, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 0, 0, 0, 0, 156, 255, 255, 255, 156, 255, 255, 255, 12, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 13, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 14, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 15, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 16, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 255, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 0, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 255, 255, 255, 255,
	}
)

func TestDeserializePlayerDataFromBytes(t *testing.T) {
	inputByteData := playerBytes
	streamReader := io.NewSectionReader(bytes.NewReader(inputByteData), 0, int64(len(inputByteData)))
	result := DeserializePlayerDataFromBytes(streamReader)
	expected := playerData

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf(`Contents not equal. Result = %v, expected = %v`, result, expected)
	}
}

func TestSerializePlayerDataToBytes(t *testing.T) {
	resultBytes := SerializePlayerDataToBytes(playerData)
	expectedBytes := playerBytes

	if !reflect.DeepEqual(resultBytes, expectedBytes) {
		t.Fatalf(`Contents not equal. Result = %v, expected = %v`, resultBytes, expectedBytes)
	}
}

func TestSerializeEmptyPlayer(t *testing.T) {
	resultBytes := SerializePlayerDataToBytes(BuildEmptyPlayer(17, "Player17", color.RGBA{100, 150, 200, 255}))
	expectedBytes := []byte{17,
		// Player name
		8, 80, 108, 97, 121, 101, 114, 49, 55,
		// 00000000-0000-0000-0000-000000000000
		36, 48, 48, 48, 48, 48, 48, 48, 48, 45, 48, 48, 48, 48, 45, 48, 48, 48, 48, 45, 48, 48, 48, 48, 45, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48,
		1, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 1, 2, 0, 0, 0,
		18, 0, 1, 0, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, 0, 0, 0, 4, 0, 0, 0, 0, 5, 0, 0, 0, 0, 6, 0, 0, 0, 0, 7, 0, 0, 0, 0,
		8, 0, 0, 0, 0, 9, 0, 0, 0, 0, 10, 0, 0, 0, 0, 11, 0, 0, 0, 0, 12, 0, 0, 0, 0, 13, 0, 0, 0, 0, 14, 0, 0, 0, 0,
		15, 0, 0, 0, 0, 16, 0, 0, 0, 0, 17, 0, 0, 0, 0, 255, 0, 0, 0, 0,
		5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		// override color
		200, 150, 100, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 255, 255, 255, 255,
	}

	compareArrays(t, resultBytes, expectedBytes)
}
