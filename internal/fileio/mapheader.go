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
	CurrentGameState   uint8
	Seed               int32
	TurnLimit          uint32
	ScoreLimit         uint32
	WinByCapital       uint8
	UnknownSettings    [6]byte
	GameModeBase       uint8
	GameModeRules      uint8
}

type MapHeaderOutput struct {
	MapHeaderInput       MapHeaderInput
	MapName              string
	MapSquareSize        int
	DisabledTribesArr    []int
	UnlockedTribesArr    []int
	GameDifficulty       int
	NumOpponents         int
	GameType             int
	MapPreset            int
	TurnTimeLimitMinutes int
	UnknownFloat1        float32
	UnknownFloat2        float32
	BaseTimeSeconds      float32
	TimeSettings         []int
	SelectedTribeSkins   []TribeSkin
	MapWidth             int
	MapHeight            int
}

type TribeSkin struct {
	Tribe int
	Skin  int
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
	gameType := unsafeReadUint16(streamReader)
	mapPreset := unsafeReadUint8(streamReader)
	turnTimeLimitMinutes := unsafeReadInt32(streamReader)
	unknownFloat1 := unsafeReadFloat32(streamReader)
	unknownFloat2 := unsafeReadFloat32(streamReader)
	baseTimeSeconds := unsafeReadFloat32(streamReader)
	timeSettings := readFixedList(streamReader, 4)

	selectedTribeSkinSize := unsafeReadUint32(streamReader)
	selectedTribeSkins := make([]TribeSkin, int(selectedTribeSkinSize))
	for i := 0; i < int(selectedTribeSkinSize); i++ {
		tribe := unsafeReadUint16(streamReader)
		skin := unsafeReadUint16(streamReader)
		selectedTribeSkins[i] = TribeSkin{
			Tribe: int(tribe),
			Skin:  int(skin),
		}
	}

	updateFileOffsetMap(fileOffsetMap, streamReader, "MapWidth")
	mapWidth := unsafeReadUint16(streamReader)
	updateFileOffsetMap(fileOffsetMap, streamReader, "MapHeight")
	mapHeight := unsafeReadUint16(streamReader)

	return MapHeaderOutput{
		MapHeaderInput:       mapHeaderInput,
		MapName:              mapName,
		MapSquareSize:        squareSize,
		DisabledTribesArr:    disabledTribesArr,
		UnlockedTribesArr:    unlockedTribesArr,
		GameDifficulty:       int(gameDifficulty),
		NumOpponents:         int(numOpponents),
		GameType:             int(gameType),
		MapPreset:            int(mapPreset),
		TurnTimeLimitMinutes: int(turnTimeLimitMinutes),
		UnknownFloat1:        unknownFloat1,
		UnknownFloat2:        unknownFloat2,
		BaseTimeSeconds:      baseTimeSeconds,
		TimeSettings:         convertByteListToInt(timeSettings),
		SelectedTribeSkins:   selectedTribeSkins,
		MapWidth:             int(mapWidth),
		MapHeight:            int(mapHeight),
	}
}

func SerializeMapHeaderToBytes(mapHeaderOutput MapHeaderOutput) []byte {
	serializedData := make([]byte, 0)

	serializedData = append(serializedData, SerializeMapHeaderInputToBytes(mapHeaderOutput.MapHeaderInput)...)
	serializedData = append(serializedData, ConvertVarString(mapHeaderOutput.MapName)...)
	serializedData = append(serializedData, ConvertUint32Bytes(mapHeaderOutput.MapSquareSize)...)

	serializedData = append(serializedData, ConvertUint16Bytes(len(mapHeaderOutput.DisabledTribesArr))...)
	for i := 0; i < len(mapHeaderOutput.DisabledTribesArr); i++ {
		serializedData = append(serializedData, ConvertUint16Bytes(mapHeaderOutput.DisabledTribesArr[i])...)
	}

	serializedData = append(serializedData, ConvertUint16Bytes(len(mapHeaderOutput.UnlockedTribesArr))...)
	for i := 0; i < len(mapHeaderOutput.UnlockedTribesArr); i++ {
		serializedData = append(serializedData, ConvertUint16Bytes(mapHeaderOutput.UnlockedTribesArr[i])...)
	}

	serializedData = append(serializedData, ConvertUint16Bytes(mapHeaderOutput.GameDifficulty)...)
	serializedData = append(serializedData, ConvertUint32Bytes(mapHeaderOutput.NumOpponents)...)
	serializedData = append(serializedData, ConvertUint16Bytes(mapHeaderOutput.GameType)...)
	serializedData = append(serializedData, byte(mapHeaderOutput.MapPreset))
	serializedData = append(serializedData, ConvertUint32Bytes(mapHeaderOutput.TurnTimeLimitMinutes)...)
	serializedData = append(serializedData, ConvertFloat32Bytes(mapHeaderOutput.UnknownFloat1)...)
	serializedData = append(serializedData, ConvertFloat32Bytes(mapHeaderOutput.UnknownFloat2)...)
	serializedData = append(serializedData, ConvertFloat32Bytes(mapHeaderOutput.BaseTimeSeconds)...)
	serializedData = append(serializedData, ConvertByteList(mapHeaderOutput.TimeSettings)...)

	serializedData = append(serializedData, ConvertUint32Bytes(len(mapHeaderOutput.SelectedTribeSkins))...)
	for i := 0; i < len(mapHeaderOutput.SelectedTribeSkins); i++ {
		serializedData = append(serializedData, ConvertUint16Bytes(mapHeaderOutput.SelectedTribeSkins[i].Tribe)...)
		serializedData = append(serializedData, ConvertUint16Bytes(mapHeaderOutput.SelectedTribeSkins[i].Skin)...)
	}

	serializedData = append(serializedData, ConvertUint16Bytes(mapHeaderOutput.MapWidth)...)
	serializedData = append(serializedData, ConvertUint16Bytes(mapHeaderOutput.MapHeight)...)

	return serializedData
}

func SerializeMapHeaderInputToBytes(mapHeaderInput MapHeaderInput) []byte {
	serializedData := make([]byte, 0)

	serializedData = append(serializedData, ConvertUint32Bytes(int(mapHeaderInput.Version1))...)
	serializedData = append(serializedData, ConvertUint32Bytes(int(mapHeaderInput.Version2))...)
	serializedData = append(serializedData, ConvertUint16Bytes(int(mapHeaderInput.TotalActions))...)
	serializedData = append(serializedData, ConvertUint32Bytes(int(mapHeaderInput.CurrentTurn))...)
	serializedData = append(serializedData, mapHeaderInput.CurrentPlayerIndex)
	serializedData = append(serializedData, ConvertUint32Bytes(int(mapHeaderInput.MaxUnitId))...)
	serializedData = append(serializedData, mapHeaderInput.CurrentGameState)
	serializedData = append(serializedData, ConvertUint32Bytes(int(mapHeaderInput.Seed))...)
	serializedData = append(serializedData, ConvertUint32Bytes(int(mapHeaderInput.TurnLimit))...)
	serializedData = append(serializedData, ConvertUint32Bytes(int(mapHeaderInput.ScoreLimit))...)
	serializedData = append(serializedData, mapHeaderInput.WinByCapital)
	serializedData = append(serializedData, mapHeaderInput.UnknownSettings[:]...)
	serializedData = append(serializedData, mapHeaderInput.GameModeBase)
	serializedData = append(serializedData, mapHeaderInput.GameModeRules)

	return serializedData
}
