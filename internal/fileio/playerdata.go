package fileio

import (
	"encoding/binary"
	"image/color"
	"io"
	"log"
)

type PlayerData struct {
	Id                   int
	Name                 string
	AccountId            string
	AutoPlay             bool
	StartTileCoordinates [2]int
	Tribe                int
	UnknownByte1         int
	UnknownInt1          int
	UnknownArr1          []PlayerUnknownData
	Currency             int
	Score                int
	UnknownInt2          int
	NumCities            int
	AvailableTech        []int
	EncounteredPlayers   []int
	Tasks                []PlayerTaskData
	TotalUnitsKilled     int
	TotalUnitsLost       int
	TotalTribesDestroyed int
	OverrideColor        []int
	UnknownByte2         byte
	UniqueImprovements   []int
	DiplomacyArr         []DiplomacyData
	DiplomacyMessages    []DiplomacyMessage
	DestroyedByTribe     int
	DestroyedTurn        int
	UnknownBuffer2       []int
}

type PlayerUnknownData struct {
	PlayerId int
	Unknown1 int
	Unknown2 int
	Unknown3 int
	Unknown4 int
}

type PlayerTaskData struct {
	Type   int
	Buffer []int
}

type DiplomacyMessage struct {
	MessageType int
	Sender      int
}

type DiplomacyData struct {
	PlayerId               uint8
	DiplomacyRelationState uint8
	LastAttackTurn         int32
	EmbassyLevel           uint8
	LastPeaceBrokenTurn    int32
	FirstMeet              int32
	EmbassyBuildTurn       int32
	PreviousAttackTurn     int32
}

func DeserializePlayerDataFromBytes(streamReader *io.SectionReader) PlayerData {
	playerId := unsafeReadUint8(streamReader)
	playerName := readVarString(streamReader, "playerName")
	playerAccountId := readVarString(streamReader, "playerAccountId")
	autoPlay := unsafeReadUint8(streamReader)
	startTileCoordinates1 := unsafeReadInt32(streamReader)
	startTileCoordinates2 := unsafeReadInt32(streamReader)
	tribe := unsafeReadUint16(streamReader)
	unknownByte1 := unsafeReadUint8(streamReader)
	unknownInt1 := unsafeReadUint32(streamReader)

	playerArr1Key := buildPlayerArr1Key(int(playerId))
	updateFileOffsetMap(fileOffsetMap, streamReader, playerArr1Key)
	unknownArrLen1 := unsafeReadUint16(streamReader)
	unknownArr1 := make([]PlayerUnknownData, 0)
	for i := 0; i < int(unknownArrLen1); i++ {
		playerIdOther := unsafeReadUint8(streamReader)
		value2 := readFixedList(streamReader, 4)
		unknownArr1 = append(unknownArr1, PlayerUnknownData{
			PlayerId: int(playerIdOther),
			Unknown1: int(value2[0]),
			Unknown2: int(value2[1]),
			Unknown3: int(value2[2]),
			Unknown4: int(value2[3]),
		})
	}

	playerCurrencyKey := buildPlayerCurrencyKey(int(playerId))
	updateFileOffsetMap(fileOffsetMap, streamReader, playerCurrencyKey)
	currency := unsafeReadUint32(streamReader)

	score := unsafeReadUint32(streamReader)
	unknownInt2 := unsafeReadUint32(streamReader)
	numCities := unsafeReadUint16(streamReader)

	techArrayLen := unsafeReadUint16(streamReader)
	techArray := make([]int, techArrayLen)
	for i := 0; i < int(techArrayLen); i++ {
		techType := unsafeReadUint16(streamReader)
		techArray[i] = int(techType)
	}

	encounteredPlayersLen := unsafeReadUint16(streamReader)
	encounteredPlayers := make([]int, 0)
	for i := 0; i < int(encounteredPlayersLen); i++ {
		playerId := unsafeReadUint8(streamReader)
		encounteredPlayers = append(encounteredPlayers, int(playerId))
	}

	numTasks := unsafeReadInt16(streamReader)
	taskArr := make([]PlayerTaskData, int(numTasks))
	for i := 0; i < int(numTasks); i++ {
		taskType := unsafeReadInt16(streamReader)

		var buffer []byte
		if taskType == 1 || taskType == 5 { // Task type 1 is Pacifist, type 5 is Killer
			buffer = readFixedList(streamReader, 6) // Extra buffer contains a uint32
		} else if taskType >= 1 && taskType <= 8 {
			buffer = readFixedList(streamReader, 2)
		} else {
			log.Fatal("Invalid task type:", taskType)
		}
		taskArr[i] = PlayerTaskData{
			Type:   int(taskType),
			Buffer: convertByteListToInt(buffer),
		}
	}

	totalKills := unsafeReadInt32(streamReader)
	totalLosses := unsafeReadInt32(streamReader)
	totalTribesDestroyed := unsafeReadInt32(streamReader)
	overrideColor := convertByteListToInt(readFixedList(streamReader, 4))
	unknownByte2 := unsafeReadUint8(streamReader)

	playerUniqueImprovementsSize := unsafeReadUint16(streamReader)
	playerUniqueImprovements := make([]int, int(playerUniqueImprovementsSize))
	for i := 0; i < int(playerUniqueImprovementsSize); i++ {
		improvement := unsafeReadUint16(streamReader)
		playerUniqueImprovements[i] = int(improvement)
	}

	diplomacyArrLen := unsafeReadUint16(streamReader)
	diplomacyArr := make([]DiplomacyData, int(diplomacyArrLen))
	for i := 0; i < len(diplomacyArr); i++ {
		diplomacyData := DiplomacyData{}
		if err := binary.Read(streamReader, binary.LittleEndian, &diplomacyData); err != nil {
			log.Fatal("Failed to load diplomacyData: ", err)
		}
		diplomacyArr[i] = diplomacyData
	}

	diplomacyMessagesSize := unsafeReadUint16(streamReader)
	diplomacyMessagesArr := make([]DiplomacyMessage, int(diplomacyMessagesSize))
	for i := 0; i < int(diplomacyMessagesSize); i++ {
		messageType := unsafeReadUint8(streamReader)
		sender := unsafeReadUint8(streamReader)

		diplomacyMessagesArr[i] = DiplomacyMessage{
			MessageType: int(messageType),
			Sender:      int(sender),
		}
	}

	destroyedByTribe := unsafeReadUint8(streamReader)
	destroyedTurn := unsafeReadUint32(streamReader)
	unknownBuffer2 := convertByteListToInt(readFixedList(streamReader, 14))

	return PlayerData{
		Id:                   int(playerId),
		Name:                 playerName,
		AccountId:            playerAccountId,
		AutoPlay:             int(autoPlay) != 0,
		StartTileCoordinates: [2]int{int(startTileCoordinates1), int(startTileCoordinates2)},
		Tribe:                int(tribe),
		UnknownByte1:         int(unknownByte1),
		UnknownInt1:          int(unknownInt1),
		UnknownArr1:          unknownArr1,
		Currency:             int(currency),
		Score:                int(score),
		UnknownInt2:          int(unknownInt2),
		NumCities:            int(numCities),
		AvailableTech:        techArray,
		EncounteredPlayers:   encounteredPlayers,
		Tasks:                taskArr,
		TotalUnitsKilled:     int(totalKills),
		TotalUnitsLost:       int(totalLosses),
		TotalTribesDestroyed: int(totalTribesDestroyed),
		OverrideColor:        overrideColor,
		UnknownByte2:         unknownByte2,
		UniqueImprovements:   playerUniqueImprovements,
		DiplomacyArr:         diplomacyArr,
		DiplomacyMessages:    diplomacyMessagesArr,
		DestroyedByTribe:     int(destroyedByTribe),
		DestroyedTurn:        int(destroyedTurn),
		UnknownBuffer2:       unknownBuffer2,
	}
}

func SerializePlayerDataToBytes(playerData PlayerData) []byte {
	allPlayerData := make([]byte, 0)

	allPlayerData = append(allPlayerData, byte(playerData.Id))
	allPlayerData = append(allPlayerData, ConvertVarString(playerData.Name)...)
	allPlayerData = append(allPlayerData, ConvertVarString(playerData.AccountId)...)
	allPlayerData = append(allPlayerData, ConvertBoolToByte(playerData.AutoPlay))
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(playerData.StartTileCoordinates[0])...)
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(playerData.StartTileCoordinates[1])...)
	allPlayerData = append(allPlayerData, ConvertUint16Bytes(playerData.Tribe)...)
	allPlayerData = append(allPlayerData, byte(playerData.UnknownByte1))
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(playerData.UnknownInt1)...)

	allPlayerData = append(allPlayerData, ConvertUint16Bytes(len(playerData.UnknownArr1))...)
	for i := 0; i < len(playerData.UnknownArr1); i++ {
		allPlayerData = append(allPlayerData, byte(playerData.UnknownArr1[i].PlayerId), byte(playerData.UnknownArr1[i].Unknown1),
			byte(playerData.UnknownArr1[i].Unknown2), byte(playerData.UnknownArr1[i].Unknown3), byte(playerData.UnknownArr1[i].Unknown4))
	}

	allPlayerData = append(allPlayerData, ConvertUint32Bytes(playerData.Currency)...)
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(playerData.Score)...)
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(playerData.UnknownInt2)...)
	allPlayerData = append(allPlayerData, ConvertUint16Bytes(playerData.NumCities)...)

	allPlayerData = append(allPlayerData, ConvertUint16Bytes(len(playerData.AvailableTech))...)
	for i := 0; i < len(playerData.AvailableTech); i++ {
		allPlayerData = append(allPlayerData, ConvertUint16Bytes(playerData.AvailableTech[i])...)
	}

	allPlayerData = append(allPlayerData, ConvertUint16Bytes(len(playerData.EncounteredPlayers))...)
	for i := 0; i < len(playerData.EncounteredPlayers); i++ {
		allPlayerData = append(allPlayerData, byte(playerData.EncounteredPlayers[i]))
	}

	allPlayerData = append(allPlayerData, ConvertUint16Bytes(len(playerData.Tasks))...)
	for i := 0; i < len(playerData.Tasks); i++ {
		allPlayerData = append(allPlayerData, ConvertUint16Bytes(playerData.Tasks[i].Type)...)
		allPlayerData = append(allPlayerData, ConvertByteList(playerData.Tasks[i].Buffer)...)
	}

	allPlayerData = append(allPlayerData, ConvertUint32Bytes(playerData.TotalUnitsKilled)...)
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(playerData.TotalUnitsLost)...)
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(playerData.TotalTribesDestroyed)...)
	allPlayerData = append(allPlayerData, ConvertByteList(playerData.OverrideColor)...)

	allPlayerData = append(allPlayerData, playerData.UnknownByte2)

	allPlayerData = append(allPlayerData, ConvertUint16Bytes(len(playerData.UniqueImprovements))...)
	for i := 0; i < len(playerData.UniqueImprovements); i++ {
		allPlayerData = append(allPlayerData, ConvertUint16Bytes(playerData.UniqueImprovements[i])...)
	}

	allPlayerData = append(allPlayerData, ConvertUint16Bytes(len(playerData.DiplomacyArr))...)
	for i := 0; i < len(playerData.DiplomacyArr); i++ {
		allPlayerData = append(allPlayerData, SerializeDiplomacyDataToBytes(playerData.DiplomacyArr[i])...)
	}

	allPlayerData = append(allPlayerData, ConvertUint16Bytes(len(playerData.DiplomacyMessages))...)
	for i := 0; i < len(playerData.DiplomacyMessages); i++ {
		allPlayerData = append(allPlayerData, byte(playerData.DiplomacyMessages[i].MessageType))
		allPlayerData = append(allPlayerData, byte(playerData.DiplomacyMessages[i].Sender))
	}

	allPlayerData = append(allPlayerData, byte(playerData.DestroyedByTribe))
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(playerData.DestroyedTurn)...)
	allPlayerData = append(allPlayerData, ConvertByteList(playerData.UnknownBuffer2)...)

	return allPlayerData
}

func SerializeDiplomacyDataToBytes(diplomacyData DiplomacyData) []byte {
	data := make([]byte, 0)
	data = append(data, byte(diplomacyData.PlayerId))
	data = append(data, byte(diplomacyData.DiplomacyRelationState))
	data = append(data, ConvertUint32Bytes(int(diplomacyData.LastAttackTurn))...)
	data = append(data, byte(diplomacyData.EmbassyLevel))
	data = append(data, ConvertUint32Bytes(int(diplomacyData.LastPeaceBrokenTurn))...)
	data = append(data, ConvertUint32Bytes(int(diplomacyData.FirstMeet))...)
	data = append(data, ConvertUint32Bytes(int(diplomacyData.EmbassyBuildTurn))...)
	data = append(data, ConvertUint32Bytes(int(diplomacyData.PreviousAttackTurn))...)
	return data
}

func BuildEmptyPlayer(index int, playerName string, overrideColor color.RGBA) PlayerData {
	if index >= 254 {
		log.Fatal("Over 255 players")
	}

	// unknown array
	newArraySize := index + 1
	unknownArr1 := make([]PlayerUnknownData, 0)
	for i := 1; i <= int(newArraySize); i++ {
		playerId := i
		if i == newArraySize {
			playerId = 255
		}
		unknownArr1 = append(unknownArr1, PlayerUnknownData{
			PlayerId: playerId,
			Unknown1: 0,
			Unknown2: 0,
			Unknown3: 0,
			Unknown4: 0,
		})
	}

	playerData := PlayerData{
		Id:                   index,
		Name:                 playerName,
		AccountId:            "00000000-0000-0000-0000-000000000000",
		AutoPlay:             true,
		StartTileCoordinates: [2]int{0, 0},
		Tribe:                2, // Ai-mo
		UnknownByte1:         1,
		UnknownInt1:          2,
		UnknownArr1:          unknownArr1,
		Currency:             5,
		Score:                0,
		UnknownInt2:          0,
		NumCities:            1,
		AvailableTech:        []int{},
		EncounteredPlayers:   []int{},
		Tasks:                []PlayerTaskData{},
		TotalUnitsKilled:     0,
		TotalUnitsLost:       0,
		TotalTribesDestroyed: 0,
		OverrideColor:        []int{int(overrideColor.B), int(overrideColor.G), int(overrideColor.R), 0},
		UnknownByte2:         0,
		UniqueImprovements:   []int{},
		DiplomacyArr:         []DiplomacyData{},
		DiplomacyMessages:    []DiplomacyMessage{},
		DestroyedByTribe:     0,
		DestroyedTurn:        0,
		UnknownBuffer2:       []int{255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 255, 255, 255, 255},
	}

	return playerData
}