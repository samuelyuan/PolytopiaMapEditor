package fileio

import (
	"fmt"
	"reflect"
	"testing"
)

var (
	oldPlayerRelationData = []PlayerAggression{
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
	}
	newPlayerRelationData = []PlayerAggression{
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
		{PlayerId: 17, Aggression: 0}, // new line for player id 17
		{PlayerId: 255, Aggression: 0},
	}
)

func TestBuildNewPlayerUnknownArr(t *testing.T) {
	oldArr := oldPlayerRelationData
	resultBytesNoChange := BuildNewPlayerUnknownArr(oldArr, 16)
	expectedBytesNoChange := oldPlayerRelationData
	if !reflect.DeepEqual(resultBytesNoChange, expectedBytesNoChange) {
		t.Fatalf(`No change failed. Result = %v, expected = %v`, resultBytesNoChange, expectedBytesNoChange)
	}

	resultBytesWithChange := BuildNewPlayerUnknownArr(oldArr, 17)
	expectedBytesWithChange := newPlayerRelationData
	if !reflect.DeepEqual(resultBytesWithChange, expectedBytesWithChange) {
		t.Fatalf(`Change to include player 17 failed. Result = %v, expected = %v`, resultBytesWithChange, expectedBytesWithChange)
	}
}

func compareArrays(t *testing.T, resultBytes []byte, expectedBytes []byte) {
	if !reflect.DeepEqual(len(resultBytes), len(expectedBytes)) {
		t.Fatalf(`Size not equal. Result = %v (size = %v), expected = %v (size = %v)`,
			resultBytes, len(resultBytes), expectedBytes, len(expectedBytes))
	}
	if !reflect.DeepEqual(resultBytes, expectedBytes) {
		findArrayDifference(resultBytes, expectedBytes)
		t.Fatalf(`Contents not equal. Result = %v, expected = %v`, resultBytes, expectedBytes)
	}
}

func findArrayDifference(resultBytes []byte, expectedBytes []byte) {
	if len(resultBytes) != len(expectedBytes) {
		fmt.Println("Array lengths not equal.")
		return
	}

	for i := 0; i < len(resultBytes); i++ {
		if resultBytes[i] != expectedBytes[i] {
			fmt.Println("Not equal at index", i, ", result:", resultBytes[i], ", expected:", expectedBytes[i])
		}
	}
}
