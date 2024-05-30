package fileio

import (
	"encoding/binary"
	"io"
	"log"
)

type MapHeaderInput struct {
	Version1           uint32
	Version2           uint32
	TotalActions       uint16
	CurrentTurn        uint32
	CurrentPlayerIndex uint8
	MaxUnitId          uint32
	UnknownByte1       uint8
	Seed               uint32
	TurnLimit          uint32
	Unknown1           [11]byte
	GameMode1          uint8
	GameMode2          uint8
}

type MapHeaderOutput struct {
	MapHeaderInput    MapHeaderInput
	MapName           string
	MapSquareSize     int
	DisabledTribesArr []int
	UnlockedTribesArr []int
	GameDifficulty    int
	NumOpponents      int
	UnknownArr        []int
	SelectedTribes    map[int]int
	MapWidth          int
	MapHeight         int
}

func DeserializeMapHeaderFromBytes(streamReader *io.SectionReader) MapHeaderOutput {
	mapHeaderInput := MapHeaderInput{}
	if err := binary.Read(streamReader, binary.LittleEndian, &mapHeaderInput); err != nil {
		log.Fatal("Failed to load MapHeaderInput: ", err)
	}

	mapName := readVarString(streamReader, "MapName")

	// map dimenions is a square: squareSize x squareSize
	updateFileOffsetMap(fileOffsetMap, streamReader, "SquareSizeKey")
	squareSize := int(unsafeReadUint32(streamReader))

	disabledTribesSize := unsafeReadUint16(streamReader)
	disabledTribesArr := make([]int, disabledTribesSize)
	for i := 0; i < int(disabledTribesSize); i++ {
		disabledTribesArr[i] = int(unsafeReadUint16(streamReader))
	}

	unlockedTribesSize := unsafeReadUint16(streamReader)
	unlockedTribesArr := make([]int, unlockedTribesSize)
	for i := 0; i < int(unlockedTribesSize); i++ {
		unlockedTribesArr[i] = int(unsafeReadUint16(streamReader))
	}

	gameDifficulty := unsafeReadUint16(streamReader)
	numOpponents := unsafeReadUint32(streamReader)
	unknownArr := readFixedList(streamReader, 5+int(unlockedTribesSize))

	selectedTribeSkinSize := unsafeReadUint32(streamReader)
	selectedTribeSkins := make(map[int]int)
	for i := 0; i < int(selectedTribeSkinSize); i++ {
		tribe := unsafeReadUint16(streamReader)
		skin := unsafeReadUint16(streamReader)
		selectedTribeSkins[int(tribe)] = int(skin)
	}

	updateFileOffsetMap(fileOffsetMap, streamReader, "MapWidth")
	mapWidth := unsafeReadUint16(streamReader)
	updateFileOffsetMap(fileOffsetMap, streamReader, "MapHeight")
	mapHeight := unsafeReadUint16(streamReader)
	if mapWidth == 0 && mapHeight == 0 {
		updateFileOffsetMap(fileOffsetMap, streamReader, "MapWidth")
		mapWidth = unsafeReadUint16(streamReader)
		updateFileOffsetMap(fileOffsetMap, streamReader, "MapHeight")
		mapHeight = unsafeReadUint16(streamReader)
	}

	return MapHeaderOutput{
		MapHeaderInput:    mapHeaderInput,
		MapName:           mapName,
		MapSquareSize:     squareSize,
		DisabledTribesArr: disabledTribesArr,
		UnlockedTribesArr: unlockedTribesArr,
		GameDifficulty:    int(gameDifficulty),
		NumOpponents:      int(numOpponents),
		UnknownArr:        convertByteListToInt(unknownArr),
		SelectedTribes:    selectedTribeSkins,
		MapWidth:          int(mapWidth),
		MapHeight:         int(mapHeight),
	}
}
