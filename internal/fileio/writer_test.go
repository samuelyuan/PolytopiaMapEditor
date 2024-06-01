package fileio

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBuildNewPlayerRelationArr(t *testing.T) {
	oldArr := []PlayerRelationData{
		{PlayerId: 1, Unknown1: 0, Unknown2: 0, Unknown3: 0, Unknown4: 0},
		{PlayerId: 2, Unknown1: 80, Unknown2: 69, Unknown3: 0, Unknown4: 0},
		{PlayerId: 3, Unknown1: 88, Unknown2: 29, Unknown3: 1, Unknown4: 0},
		{PlayerId: 4, Unknown1: 39, Unknown2: 95, Unknown3: 0, Unknown4: 0},
		{PlayerId: 5, Unknown1: 222, Unknown2: 34, Unknown3: 1, Unknown4: 0},
		{PlayerId: 6, Unknown1: 218, Unknown2: 77, Unknown3: 1, Unknown4: 0},
		{PlayerId: 7, Unknown1: 134, Unknown2: 250, Unknown3: 0, Unknown4: 0},
		{PlayerId: 8, Unknown1: 243, Unknown2: 153, Unknown3: 0, Unknown4: 0},
		{PlayerId: 9, Unknown1: 131, Unknown2: 143, Unknown3: 0, Unknown4: 0},
		{PlayerId: 10, Unknown1: 180, Unknown2: 147, Unknown3: 0, Unknown4: 0},
		{PlayerId: 11, Unknown1: 74, Unknown2: 89, Unknown3: 0, Unknown4: 0},
		{PlayerId: 12, Unknown1: 7, Unknown2: 125, Unknown3: 0, Unknown4: 0},
		{PlayerId: 13, Unknown1: 74, Unknown2: 69, Unknown3: 0, Unknown4: 0},
		{PlayerId: 14, Unknown1: 66, Unknown2: 163, Unknown3: 0, Unknown4: 0},
		{PlayerId: 15, Unknown1: 165, Unknown2: 216, Unknown3: 0, Unknown4: 0},
		{PlayerId: 16, Unknown1: 41, Unknown2: 125, Unknown3: 0, Unknown4: 0},
		{PlayerId: 255, Unknown1: 0, Unknown2: 0, Unknown3: 0, Unknown4: 0},
	}
	resultBytesNoChange := BuildNewPlayerRelationArr(oldArr, 16)
	expectedBytesNoChange := []PlayerRelationData{
		{PlayerId: 1, Unknown1: 0, Unknown2: 0, Unknown3: 0, Unknown4: 0},
		{PlayerId: 2, Unknown1: 80, Unknown2: 69, Unknown3: 0, Unknown4: 0},
		{PlayerId: 3, Unknown1: 88, Unknown2: 29, Unknown3: 1, Unknown4: 0},
		{PlayerId: 4, Unknown1: 39, Unknown2: 95, Unknown3: 0, Unknown4: 0},
		{PlayerId: 5, Unknown1: 222, Unknown2: 34, Unknown3: 1, Unknown4: 0},
		{PlayerId: 6, Unknown1: 218, Unknown2: 77, Unknown3: 1, Unknown4: 0},
		{PlayerId: 7, Unknown1: 134, Unknown2: 250, Unknown3: 0, Unknown4: 0},
		{PlayerId: 8, Unknown1: 243, Unknown2: 153, Unknown3: 0, Unknown4: 0},
		{PlayerId: 9, Unknown1: 131, Unknown2: 143, Unknown3: 0, Unknown4: 0},
		{PlayerId: 10, Unknown1: 180, Unknown2: 147, Unknown3: 0, Unknown4: 0},
		{PlayerId: 11, Unknown1: 74, Unknown2: 89, Unknown3: 0, Unknown4: 0},
		{PlayerId: 12, Unknown1: 7, Unknown2: 125, Unknown3: 0, Unknown4: 0},
		{PlayerId: 13, Unknown1: 74, Unknown2: 69, Unknown3: 0, Unknown4: 0},
		{PlayerId: 14, Unknown1: 66, Unknown2: 163, Unknown3: 0, Unknown4: 0},
		{PlayerId: 15, Unknown1: 165, Unknown2: 216, Unknown3: 0, Unknown4: 0},
		{PlayerId: 16, Unknown1: 41, Unknown2: 125, Unknown3: 0, Unknown4: 0},
		{PlayerId: 255, Unknown1: 0, Unknown2: 0, Unknown3: 0, Unknown4: 0},
	}
	if !reflect.DeepEqual(resultBytesNoChange, expectedBytesNoChange) {
		t.Fatalf(`No change failed. Result = %v, expected = %v`, resultBytesNoChange, expectedBytesNoChange)
	}

	resultBytesWithChange := BuildNewPlayerRelationArr(oldArr, 17)
	expectedBytesWithChange := []PlayerRelationData{
		{PlayerId: 1, Unknown1: 0, Unknown2: 0, Unknown3: 0, Unknown4: 0},
		{PlayerId: 2, Unknown1: 80, Unknown2: 69, Unknown3: 0, Unknown4: 0},
		{PlayerId: 3, Unknown1: 88, Unknown2: 29, Unknown3: 1, Unknown4: 0},
		{PlayerId: 4, Unknown1: 39, Unknown2: 95, Unknown3: 0, Unknown4: 0},
		{PlayerId: 5, Unknown1: 222, Unknown2: 34, Unknown3: 1, Unknown4: 0},
		{PlayerId: 6, Unknown1: 218, Unknown2: 77, Unknown3: 1, Unknown4: 0},
		{PlayerId: 7, Unknown1: 134, Unknown2: 250, Unknown3: 0, Unknown4: 0},
		{PlayerId: 8, Unknown1: 243, Unknown2: 153, Unknown3: 0, Unknown4: 0},
		{PlayerId: 9, Unknown1: 131, Unknown2: 143, Unknown3: 0, Unknown4: 0},
		{PlayerId: 10, Unknown1: 180, Unknown2: 147, Unknown3: 0, Unknown4: 0},
		{PlayerId: 11, Unknown1: 74, Unknown2: 89, Unknown3: 0, Unknown4: 0},
		{PlayerId: 12, Unknown1: 7, Unknown2: 125, Unknown3: 0, Unknown4: 0},
		{PlayerId: 13, Unknown1: 74, Unknown2: 69, Unknown3: 0, Unknown4: 0},
		{PlayerId: 14, Unknown1: 66, Unknown2: 163, Unknown3: 0, Unknown4: 0},
		{PlayerId: 15, Unknown1: 165, Unknown2: 216, Unknown3: 0, Unknown4: 0},
		{PlayerId: 16, Unknown1: 41, Unknown2: 125, Unknown3: 0, Unknown4: 0},
		{PlayerId: 17, Unknown1: 0, Unknown2: 0, Unknown3: 0, Unknown4: 0}, // new line for player id 17
		{PlayerId: 255, Unknown1: 0, Unknown2: 0, Unknown3: 0, Unknown4: 0},
	}
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
