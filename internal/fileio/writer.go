package fileio

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

func WriteUint8AtFileOffset(inputFilename string, offset int, value int) {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state: ", err)
	}

	if value >= 256 {
		log.Fatal("Value is too large for uint8")
	}
	if _, err := inputFile.WriteAt([]byte{uint8(value)}, int64(offset)); err != nil {
		log.Fatal("Failed to write uint8 to file:", err)
	}
}

func WriteUint16AtFileOffset(inputFilename string, offset int, updatedValue int) {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state: ", err)
	}

	if updatedValue >= 65536 {
		log.Fatal("Value is too large for uint16")
	}
	byteArrUnitType := make([]byte, 2)
	binary.LittleEndian.PutUint16(byteArrUnitType, uint16(updatedValue))
	if _, err := inputFile.WriteAt(byteArrUnitType, int64(offset)); err != nil {
		log.Fatal(err)
	}
}

func WriteUint32AtFileOffset(inputFilename string, offset int, updatedValue int) {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state: ", err)
	}

	if updatedValue >= 4294967295 {
		log.Fatal("Value is too large for uint32")
	}
	byteArrUnitType := make([]byte, 4)
	binary.LittleEndian.PutUint32(byteArrUnitType, uint32(updatedValue))
	if _, err := inputFile.WriteAt(byteArrUnitType, int64(offset)); err != nil {
		log.Fatal(err)
	}
}

func GetFileRemainingData(inputFile *os.File, offset int) []byte {
	if _, err := inputFile.Seek(int64(offset), 0); err != nil {
		log.Fatal(err)
	}
	remainder, err := ioutil.ReadAll(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	return remainder
}

func WriteAndShiftData(inputFilename string, offsetStartOriginalBlockKey string, offsetEndOriginalBlockKey string, newData []byte) {
	// Update file offsets to make sure they are up to date
	_, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}

	// Open file to modify
	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state:", err)
	}

	offsetOriginalBlockStart, ok := fileOffsetMap[offsetStartOriginalBlockKey]
	if !ok {
		log.Fatal(fmt.Sprintf("Error: Unable to find start of data block with key %v. Command not run.", offsetStartOriginalBlockKey))
	}
	offsetOriginalBlockEnd, ok := fileOffsetMap[offsetEndOriginalBlockKey]
	if !ok {
		log.Fatal(fmt.Sprintf("Error: Unable to find end of data block with key %v. Command not run.", offsetStartOriginalBlockKey))
	}
	// Get all data after end of block
	remainder := GetFileRemainingData(inputFile, offsetOriginalBlockEnd)

	// overwrite block with new data at original block start
	if _, err := inputFile.WriteAt(newData, int64(offsetOriginalBlockStart)); err != nil {
		log.Fatal(err)
	}

	// shift remaining data and write after new data instead of original end start
	if _, err := inputFile.WriteAt(remainder, int64(offsetOriginalBlockStart+len(newData))); err != nil {
		log.Fatal(err)
	}
}

func ConvertUint32Bytes(value int) []byte {
	byteArr := make([]byte, 4)
	binary.LittleEndian.PutUint32(byteArr, uint32(value))
	return byteArr
}

func ConvertUint16Bytes(value int) []byte {
	byteArr := make([]byte, 2)
	binary.LittleEndian.PutUint16(byteArr, uint16(value))
	return byteArr
}

func ConvertVarString(value string) []byte {
	byteArr := make([]byte, 0)
	byteArr = append(byteArr, byte(len(value)))
	byteArr = append(byteArr, []byte(value)...)
	return byteArr
}

func ConvertByteList(oldArr []int) []byte {
	// the values are stored as ints but they were originally bytes
	newArr := make([]byte, len(oldArr))
	for i := 0; i < len(oldArr); i++ {
		if oldArr[i] > 255 {
			log.Fatal(fmt.Sprintf("Byte list has value over 255. Value is %v for index %v", oldArr[i], i))
		}
		newArr[i] = byte(oldArr[i])
	}
	return newArr
}

func ConvertBoolToByte(value bool) byte {
	if value {
		return 1
	} else {
		return 0
	}
}

func ConvertImprovementDataToBytes(improvementData ImprovementData) []byte {
	data := make([]byte, 0)
	data = append(data, ConvertUint16Bytes(int(improvementData.Level))...)
	data = append(data, ConvertUint16Bytes(int(improvementData.FoundedTurn))...)
	data = append(data, ConvertUint16Bytes(int(improvementData.CurrentPopulation))...)
	data = append(data, ConvertUint16Bytes(int(improvementData.TotalPopulation))...)
	data = append(data, ConvertUint16Bytes(int(improvementData.Production))...)
	data = append(data, ConvertUint16Bytes(int(improvementData.BaseScore))...)
	data = append(data, ConvertUint16Bytes(int(improvementData.BorderSize))...)
	data = append(data, ConvertUint16Bytes(int(improvementData.UpgradeCount))...)
	data = append(data, byte(improvementData.ConnectedPlayerCapital))
	data = append(data, byte(improvementData.HasCityName))
	if improvementData.HasCityName == 1 {
		data = append(data, ConvertVarString(improvementData.CityName)...)
	}
	data = append(data, byte(improvementData.FoundedTribe))
	data = append(data, ConvertUint16Bytes(len(improvementData.CityRewards))...)
	for i := 0; i < len(improvementData.CityRewards); i++ {
		data = append(data, ConvertUint16Bytes(int(improvementData.CityRewards[i]))...)
	}
	data = append(data, ConvertUint16Bytes(int(improvementData.RebellionFlag))...)
	if improvementData.RebellionFlag != 0 {
		data = append(data, ConvertByteList(improvementData.RebellionBuffer)...)
	}
	return data
}

func ConvertUnitDataToBytes(unitData UnitData) []byte {
	data := make([]byte, 0)
	data = append(data, ConvertUint32Bytes(int(unitData.Id))...)
	data = append(data, byte(unitData.Owner))
	data = append(data, ConvertUint16Bytes(int(unitData.UnitType))...)
	data = append(data, ConvertUint32Bytes(int(unitData.FollowerUnitId))...)
	data = append(data, ConvertUint32Bytes(int(unitData.LeaderUnitId))...)
	data = append(data, ConvertUint32Bytes(int(unitData.CurrentCoordinates[0]))...)
	data = append(data, ConvertUint32Bytes(int(unitData.CurrentCoordinates[1]))...)
	data = append(data, ConvertUint32Bytes(int(unitData.HomeCoordinates[0]))...)
	data = append(data, ConvertUint32Bytes(int(unitData.HomeCoordinates[1]))...)
	data = append(data, ConvertUint16Bytes(int(unitData.Health))...)
	data = append(data, ConvertUint16Bytes(int(unitData.PromotionLevel))...)
	data = append(data, ConvertUint16Bytes(int(unitData.Experience))...)
	data = append(data, ConvertBoolToByte(unitData.Moved))
	data = append(data, ConvertBoolToByte(unitData.Attacked))
	data = append(data, ConvertBoolToByte(unitData.Flipped))
	data = append(data, ConvertUint16Bytes(int(unitData.CreatedTurn))...)
	return data
}

func ConvertTileToBytes(tileData TileData) []byte {
	tileBytes := make([]byte, 0)

	headerBytes := make([]byte, 0)
	headerBytes = append(headerBytes, ConvertUint32Bytes(int(tileData.WorldCoordinates[0]))...)
	headerBytes = append(headerBytes, ConvertUint32Bytes(int(tileData.WorldCoordinates[1]))...)
	headerBytes = append(headerBytes, ConvertUint16Bytes(int(tileData.Terrain))...)
	headerBytes = append(headerBytes, ConvertUint16Bytes(int(tileData.Climate))...)
	headerBytes = append(headerBytes, ConvertUint16Bytes(int(tileData.Altitude))...) // should be int16
	headerBytes = append(headerBytes, byte(tileData.Owner))
	headerBytes = append(headerBytes, byte(tileData.Capital))
	headerBytes = append(headerBytes, ConvertUint32Bytes(int(tileData.CapitalCoordinates[0]))...) // should be int32
	headerBytes = append(headerBytes, ConvertUint32Bytes(int(tileData.CapitalCoordinates[1]))...) // should be int32
	tileBytes = append(tileBytes, headerBytes...)

	if tileData.ResourceExists {
		tileBytes = append(tileBytes, byte(1))
		tileBytes = append(tileBytes, ConvertUint16Bytes(tileData.ResourceType)...)
	} else {
		tileBytes = append(tileBytes, byte(0))
	}

	if tileData.ImprovementExists {
		tileBytes = append(tileBytes, byte(1))
		tileBytes = append(tileBytes, ConvertUint16Bytes(tileData.ImprovementType)...)
	} else {
		tileBytes = append(tileBytes, byte(0))
	}

	if tileData.ImprovementData != nil {
		tileBytes = append(tileBytes, ConvertImprovementDataToBytes(*tileData.ImprovementData)...)
	}

	// no unit
	if tileData.Unit != nil {
		tileBytes = append(tileBytes, 1)
		tileBytes = append(tileBytes, ConvertUnitDataToBytes(*tileData.Unit)...)

		if tileData.PassengerUnit != nil {
			tileBytes = append(tileBytes, 1) // has other unit flag is 1
			tileBytes = append(tileBytes, ConvertUnitDataToBytes(*tileData.PassengerUnit)...)

			// unknown flag, might be to check if passnger unit has another passnger unit
			// should always be zero
			tileBytes = append(tileBytes, 0)

			tileBytes = append(tileBytes, ConvertUint16Bytes(len(tileData.PassengerUnitEffectData))...)
			for i := 0; i < len(tileData.PassengerUnitEffectData); i++ {
				tileBytes = append(tileBytes, ConvertUint16Bytes(tileData.PassengerUnitEffectData[i])...)
			}
			tileBytes = append(tileBytes, ConvertByteList(tileData.PassengerUnitDirectionData)...)

			tileBytes = append(tileBytes, ConvertUint16Bytes(len(tileData.UnitEffectData))...)
			for i := 0; i < len(tileData.UnitEffectData); i++ {
				tileBytes = append(tileBytes, ConvertUint16Bytes(tileData.UnitEffectData[i])...)
			}
			tileBytes = append(tileBytes, ConvertByteList(tileData.UnitDirectionData)...)
		} else {
			tileBytes = append(tileBytes, 0) // has other unit flag is 0

			tileBytes = append(tileBytes, ConvertUint16Bytes(len(tileData.UnitEffectData))...)
			for i := 0; i < len(tileData.UnitEffectData); i++ {
				tileBytes = append(tileBytes, ConvertUint16Bytes(tileData.UnitEffectData[i])...)
			}
			tileBytes = append(tileBytes, ConvertByteList(tileData.UnitDirectionData)...)
		}
	} else {
		tileBytes = append(tileBytes, 0)
	}

	tileBytes = append(tileBytes, byte(len(tileData.PlayerVisibility)))
	tileBytes = append(tileBytes, ConvertByteList(tileData.PlayerVisibility)...)
	tileBytes = append(tileBytes, ConvertBoolToByte(tileData.HasRoad))
	tileBytes = append(tileBytes, ConvertBoolToByte(tileData.HasWaterRoute))
	tileBytes = append(tileBytes, ConvertByteList(tileData.Unknown)...)
	return tileBytes
}

func ConvertMapDataToBytes(tileData [][]TileData) []byte {
	mapHeight := len(tileData)
	mapWidth := len(tileData[0])

	allMapBytes := make([]byte, 0)
	for i := 0; i < mapHeight; i++ {
		for j := 0; j < mapWidth; j++ {
			tileBytes := ConvertTileToBytes(tileData[i][j])
			allMapBytes = append(allMapBytes, tileBytes...)
		}
	}
	return allMapBytes
}

func ConvertDiplomacyDataToBytes(diplomacyData DiplomacyData) []byte {
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

func ConvertPlayerDataToBytes(playerData PlayerData) []byte {
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
		allPlayerData = append(allPlayerData, ConvertDiplomacyDataToBytes(playerData.DiplomacyArr[i])...)
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

func ConvertAllPlayerDataToBytes(allPlayerData []PlayerData) []byte {
	allPlayerBytes := make([]byte, 0)
	allPlayerBytes = append(allPlayerBytes, ConvertUint16Bytes(len(allPlayerData))...)
	for i := 0; i < len(allPlayerData); i++ {
		allPlayerBytes = append(allPlayerBytes, ConvertPlayerDataToBytes(allPlayerData[i])...)
	}
	return allPlayerBytes
}

func WriteTileToFile(inputFilename string, tileDataOverwrite TileData, targetX int, targetY int) {
	tileBytes := ConvertTileToBytes(tileDataOverwrite)
	WriteAndShiftData(inputFilename, buildTileStartKey(targetX, targetY), buildTileEndKey(targetX, targetY), tileBytes)
}

func WriteMapToFile(inputFilename string, tileDataOverwrite [][]TileData) {
	allTileBytes := ConvertMapDataToBytes(tileDataOverwrite)
	WriteAndShiftData(inputFilename, buildMapStartKey(), buildMapEndKey(), allTileBytes)
}

func WritePlayersToFile(inputFilename string, playersList []PlayerData) {
	allPlayerBytes := ConvertAllPlayerDataToBytes(playersList)
	WriteAndShiftData(inputFilename, buildAllPlayersStartKey(), buildAllPlayersEndKey(), allPlayerBytes)
}

func ModifyTileTerrain(inputFilename string, targetX int, targetY int, updatedValue int) {
	saveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}

	updatedTile := saveOutput.TileData[targetY][targetX]

	// write terrain
	updatedTile.Terrain = updatedValue

	// write altitude
	altitude := 0
	if updatedValue == 1 { // water altitude is -1
		altitude = -1
	} else if updatedValue == 2 { // ocean altitude is -2
		altitude = -2
	} else if updatedValue == 3 || updatedValue == 5 { // flat tile altitude is 1
		altitude = 1
	} else if updatedValue == 4 { // mountain altitude is 2
		altitude = 2
	}
	updatedTile.Altitude = altitude

	WriteTileToFile(inputFilename, updatedTile, targetX, targetY)
}

func ModifyUnitTribe(inputFilename string, targetX int, targetY int, updatedValue int) {
	saveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}

	updatedTile := saveOutput.TileData[targetY][targetX]
	if updatedTile.Unit != nil {
		fmt.Println(fmt.Sprintf("Before changing unit's owner on tile (%v, %v), current owner is %v",
			targetX, targetY, updatedTile.Unit.Owner))
		updatedTile.Unit.Owner = uint8(updatedValue)
	} else {
		fmt.Println(fmt.Sprintf("No unit on tile (%v, %v)", targetX, targetY))
	}
	if updatedTile.PassengerUnit != nil {
		fmt.Println(fmt.Sprintf("Before changing transition unit's owner on tile (%v, %v), current owner is %v",
			targetX, targetY, updatedTile.PassengerUnit.Owner))
		updatedTile.PassengerUnit.Owner = uint8(updatedValue)
	} else {
		fmt.Println(fmt.Sprintf("No transition unit on tile (%v, %v)", targetX, targetY))
	}
	WriteTileToFile(inputFilename, updatedTile, targetX, targetY)
}

func ConvertTribe(inputFilename string, oldTribe int, newTribe int) {
	saveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}

	tribeUnits, ok := saveOutput.TribeUnitMap[oldTribe]
	if !ok {
		log.Fatal(fmt.Sprintf("Tribe %v doesn't exist", oldTribe))
	}

	for i := 0; i < len(tribeUnits); i++ {
		targetX := tribeUnits[i].X
		targetY := tribeUnits[i].Y

		updatedTile := saveOutput.TileData[targetY][targetX]
		if updatedTile.Unit != nil {
			updatedTile.Unit.Owner = uint8(newTribe)
		}
		if updatedTile.PassengerUnit != nil {
			updatedTile.PassengerUnit.Owner = uint8(newTribe)
		}
		fmt.Println(fmt.Sprintf("Converted unit on (%v, %v) from tribe %v to %v", targetX, targetY, oldTribe, newTribe))

		saveOutput.TileData[targetY][targetX] = updatedTile
	}

	WriteMapToFile(inputFilename, saveOutput.TileData)
	fmt.Println(fmt.Sprintf("Changed all units under tribe %v to tribe %v. Total of %v units converted.", oldTribe, newTribe, len(tribeUnits)))
}

func ModifyUnitType(inputFilename string, targetX int, targetY int, updatedValue int) {
	saveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}

	updatedTile := saveOutput.TileData[targetY][targetX]
	if updatedTile.Unit != nil {
		fmt.Println(fmt.Sprintf("Before changing unit's owner on tile (%v, %v), current type is %v",
			targetX, targetY, updatedTile.Unit.UnitType))
		updatedTile.Unit.UnitType = uint16(updatedValue)
	} else {
		fmt.Println(fmt.Sprintf("No unit on tile (%v, %v)", targetX, targetY))
	}
	WriteTileToFile(inputFilename, updatedTile, targetX, targetY)
}

func BuildEmptyTile(x int, y int) TileData {
	return TileData{
		WorldCoordinates:   [2]int{x, y},
		Terrain:            3,
		Climate:            1,
		Altitude:           1,
		Owner:              0,
		Capital:            0,
		CapitalCoordinates: [2]int{-1, -1},
		ResourceExists:     false,
		ResourceType:       -1,
		ImprovementExists:  false,
		ImprovementType:    -1,
		ImprovementData:    nil,
		Unit:               nil,
		UnitDirectionData:  []int{},
		PlayerVisibility:   []int{},
		HasRoad:            false,
		HasWaterRoute:      false,
		Unknown:            []int{0, 0, 0, 0},
	}
}

func ModifyMapDimensions(inputFilename string, width int, height int) {
	minSquareSize := width
	if minSquareSize > height {
		minSquareSize = height
	}
	squareSizeOffset, ok := fileOffsetMap["SquareSizeKey"]
	if !ok {
		log.Fatal("Error: No square size key. Command not run.")
	}
	WriteUint32AtFileOffset(inputFilename, squareSizeOffset, minSquareSize)

	widthOffset, ok := fileOffsetMap["MapWidth"]
	if !ok {
		log.Fatal("Error: No map width key. Command not run.")
	}
	WriteUint16AtFileOffset(inputFilename, widthOffset, width)

	heightOffset, ok := fileOffsetMap["MapHeight"]
	if !ok {
		log.Fatal("Error: No map height key. Command not run.")
	}
	WriteUint16AtFileOffset(inputFilename, heightOffset, height)
}

func BuildEmptyCity(cityName string) ImprovementData {
	return ImprovementData{
		Level:                  1,
		FoundedTurn:            0,
		CurrentPopulation:      0,
		TotalPopulation:        0,
		Production:             1,
		BaseScore:              0,
		BorderSize:             1,
		UpgradeCount:           0,
		ConnectedPlayerCapital: 0,
		HasCityName:            1,
		CityName:               cityName,
		FoundedTribe:           0,
		CityRewards:            []int{},
		RebellionFlag:          0,
		RebellionBuffer:        []int{},
	}
}

func AddCityToTile(inputFilename string, targetX int, targetY int, cityName string, tribe int) {
	saveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}
	// Overwrite tile header tribe
	// world coordinates, terrain, climate, altitude are the same for all players
	// the difference is the tribe that owns this city
	saveOutput.TileData[targetY][targetX].Owner = tribe
	// set capital to 0 unless this city is designated as capital city
	saveOutput.TileData[targetY][targetX].Capital = 0
	saveOutput.TileData[targetY][targetX].CapitalCoordinates = [2]int{targetX, targetY}
	// Overwrite improvement data and set city
	saveOutput.TileData[targetY][targetX].ImprovementExists = true
	saveOutput.TileData[targetY][targetX].ImprovementType = 1
	improvementData := BuildEmptyCity(cityName)
	saveOutput.TileData[targetY][targetX].ImprovementData = &improvementData
	WriteTileToFile(inputFilename, saveOutput.TileData[targetY][targetX], targetX, targetY)
}

func ResetTile(inputFilename string, targetX int, targetY int) {
	updatedTile := BuildEmptyTile(targetX, targetY)
	WriteTileToFile(inputFilename, updatedTile, targetX, targetY)
}

func ExpandRows(inputFilename string, newRowDimensions int) {
	if newRowDimensions >= 256 {
		log.Fatal("Updated value is over 256")
	}

	saveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}
	fmt.Println(fmt.Sprintf("Old dimensions width: %v, height: %v", saveOutput.MapWidth, saveOutput.MapHeight))

	if newRowDimensions <= saveOutput.MapHeight {
		log.Fatal(fmt.Sprintf("New row dimensions are less than existing dimensions, new value: %v, existing height: %v",
			newRowDimensions, saveOutput.MapHeight))
	}

	for y := saveOutput.MapHeight; y < newRowDimensions; y++ {
		newTileRow := make([]TileData, saveOutput.MapWidth)
		for x := 0; x < saveOutput.MapWidth; x++ {
			newTileRow[x] = BuildEmptyTile(x, y)
		}
		saveOutput.TileData = append(saveOutput.TileData, newTileRow)
	}
	WriteMapToFile(inputFilename, saveOutput.TileData)
	ModifyMapDimensions(inputFilename, saveOutput.MapWidth, newRowDimensions)

	finalSaveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	fmt.Println(fmt.Sprintf("New dimensions, width: %v, height: %v", finalSaveOutput.MapWidth, finalSaveOutput.MapHeight))
}

func ExpandColumns(inputFilename string, newColDimensions int) {
	if newColDimensions >= 256 {
		log.Fatal("Updated value is over 256")
	}

	saveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}
	fmt.Println(fmt.Sprintf("Old dimensions width: %v, height: %v", saveOutput.MapWidth, saveOutput.MapHeight))

	if newColDimensions <= saveOutput.MapWidth {
		log.Fatal(fmt.Sprintf("New column dimensions are less than existing dimensions, new value: %v, existing width: %v",
			newColDimensions, saveOutput.MapWidth))
	}

	for y := saveOutput.MapHeight - 1; y >= 0; y-- {
		for x := saveOutput.MapWidth; x < newColDimensions; x++ {
			saveOutput.TileData[y] = append(saveOutput.TileData[y], BuildEmptyTile(x, y))
		}
	}
	WriteMapToFile(inputFilename, saveOutput.TileData)
	ModifyMapDimensions(inputFilename, newColDimensions, saveOutput.MapHeight)

	finalSaveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	fmt.Println(fmt.Sprintf("New dimensions, width: %v, height: %v", finalSaveOutput.MapWidth, finalSaveOutput.MapHeight))
}

func ExpandTiles(inputFilename string, newSquareSizeDimensions int) {
	if newSquareSizeDimensions >= 256 {
		log.Fatal("Updated value is over 256")
	}

	saveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}

	if newSquareSizeDimensions <= saveOutput.MapWidth || newSquareSizeDimensions <= saveOutput.MapHeight {
		log.Fatal(fmt.Sprintf("New dimensions are less than existing dimensions, new value: %v, existing width: %v, height: %v",
			newSquareSizeDimensions, saveOutput.MapWidth, saveOutput.MapHeight))
	}

	ExpandColumns(inputFilename, newSquareSizeDimensions)
	ExpandRows(inputFilename, newSquareSizeDimensions)
}

func RevealAllTiles(inputFilename string, newTribe int) {
	saveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}

	for i := saveOutput.MapHeight - 1; i >= 0; i-- {
		for j := saveOutput.MapWidth - 1; j >= 0; j-- {
			targetX := j
			targetY := i

			visibilityData := saveOutput.TileData[i][j].PlayerVisibility
			fmt.Println("Existing visibility data:", visibilityData)
			isAlreadyVisible := false
			for visibilityIndex := 0; visibilityIndex < len(visibilityData); visibilityIndex++ {
				if int(visibilityData[visibilityIndex]) == newTribe {
					fmt.Printf("Tile is already visible to tribe %v. No change will be made to visibility data.\n", newTribe)
					isAlreadyVisible = true
					break
				}
			}
			if !isAlreadyVisible {
				saveOutput.TileData[i][j].PlayerVisibility = append(saveOutput.TileData[i][j].PlayerVisibility, newTribe)
				fmt.Println(fmt.Sprintf("Revealed (%v, %v) for tribe %v", targetX, targetY, newTribe))
			}
		}
	}

	for i := saveOutput.MapHeight - 1; i >= 0; i-- {
		for j := saveOutput.MapWidth - 1; j >= 0; j-- {
			fmt.Println(fmt.Sprintf("Tile (%v, %v) visibility: %v", j, i, saveOutput.TileData[i][j].PlayerVisibility))
		}
	}

	WriteMapToFile(inputFilename, saveOutput.TileData)
}

func RevealTileForTribe(inputFilename string, targetX int, targetY int, newTribe int) {
	saveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}

	visibilityData := saveOutput.TileData[targetY][targetX].PlayerVisibility
	fmt.Println("Existing visibility data:", visibilityData)
	isAlreadyVisible := false
	for visibilityIndex := 0; visibilityIndex < len(visibilityData); visibilityIndex++ {
		if int(visibilityData[visibilityIndex]) == newTribe {
			fmt.Printf("Tile is already visible to tribe %v. No change will be made to visibility data.\n", newTribe)
			isAlreadyVisible = true
			break
		}
	}
	if !isAlreadyVisible {
		saveOutput.TileData[targetY][targetX].PlayerVisibility = append(saveOutput.TileData[targetY][targetX].PlayerVisibility, newTribe)
		fmt.Println(fmt.Sprintf("Revealed (%v, %v) for tribe %v", targetX, targetY, newTribe))
	}
	WriteTileToFile(inputFilename, saveOutput.TileData[targetY][targetX], targetX, targetY)
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

func generateRandomColor() color.RGBA {
	rand.Seed(time.Now().UnixNano())
	return color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255}
}

func BuildNewPlayerUnknownArr(oldUnknownArr1 []PlayerUnknownData, newPlayerId int) []PlayerUnknownData {
	existingLen := len(oldUnknownArr1)

	oldPlayerCount := existingLen
	oldMaximumPlayerId := oldPlayerCount - 1 // excludes player 255 nature
	if oldMaximumPlayerId >= newPlayerId {
		fmt.Println(fmt.Sprintf("Existing player count is %v, which includes players 1 to %v. No need to add player id %v.",
			oldPlayerCount, oldPlayerCount-1, newPlayerId))
		return oldUnknownArr1
	} else {
		fmt.Println(fmt.Sprintf("Existing player count is %v, which includes players 1 to %v. New player id %v needs to be included.",
			oldPlayerCount, oldPlayerCount-1, newPlayerId))
	}

	dataInsert := PlayerUnknownData{
		PlayerId: newPlayerId,
		Unknown1: 0,
		Unknown2: 0,
		Unknown3: 0,
		Unknown4: 0,
	}
	// assumes player 255 is always last
	existingPlayers := oldUnknownArr1[0 : existingLen-1]
	naturePlayer := oldUnknownArr1[existingLen-1]

	newUnknownArr1 := append(existingPlayers, dataInsert, naturePlayer)
	return newUnknownArr1
}

func convertPlayerIndexToId(playerIndex int, totalPlayers int) int {
	if playerIndex == totalPlayers-1 {
		return 255
	} else {
		return playerIndex + 1
	}
}

func ModifyAllExistingPlayerUnknownArr(inputFilename string) {
	saveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}
	newPlayerCount := len(saveOutput.PlayerData)
	fmt.Println(fmt.Sprintf("New player count: %v", newPlayerCount))

	for i := len(saveOutput.PlayerData) - 1; i >= 0; i-- {
		newPlayerId := newPlayerCount - 1
		newUnknownArr1 := BuildNewPlayerUnknownArr(saveOutput.PlayerData[i].UnknownArr1, newPlayerId)
		saveOutput.PlayerData[i].UnknownArr1 = newUnknownArr1
	}
	WritePlayersToFile(inputFilename, saveOutput.PlayerData)
}

func AddPlayer(inputFilename string) {
	saveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}
	oldPlayerCount := len(saveOutput.PlayerData)
	fmt.Println(fmt.Sprintf("Old num players: %v", oldPlayerCount))

	// existing index will be 1, 2, 3, ..., oldPlayerCount-1, 255 (size is oldPlayerCount)
	// new index list will be 1, 2, 3, ..., oldPlayerCount-1, oldPlayerCount, 255 (size is oldPlayerCount + 1)
	playerName := fmt.Sprintf("Player%v", oldPlayerCount)
	overrideColor := generateRandomColor()
	newPlayer := BuildEmptyPlayer(oldPlayerCount, playerName, overrideColor)

	newPlayerData := make([]PlayerData, 0)
	for i := 0; i < len(saveOutput.PlayerData)-1; i++ {
		newPlayerData = append(newPlayerData, saveOutput.PlayerData[i])
	}
	newPlayerData = append(newPlayerData, newPlayer)
	newPlayerData = append(newPlayerData, saveOutput.PlayerData[len(saveOutput.PlayerData)-1])
	WritePlayersToFile(inputFilename, newPlayerData)

	ModifyAllExistingPlayerUnknownArr(inputFilename)
}

func SwapPlayers(inputFilename string, playerId1 int, playerId2 int) {
	saveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}

	// assumes 254 is not used by any players
	unusedPlayerId := 254

	// Need to reassign so we don't merge the players
	for i := 0; i < saveOutput.MapHeight; i++ {
		for j := 0; j < saveOutput.MapWidth; j++ {
			if saveOutput.TileData[i][j].Owner == playerId1 {
				saveOutput.TileData[i][j].Owner = unusedPlayerId
			}

			if saveOutput.TileData[i][j].Capital == playerId1 {
				saveOutput.TileData[i][j].Capital = unusedPlayerId
			}
			if saveOutput.TileData[i][j].ImprovementData != nil && saveOutput.TileData[i][j].ImprovementData.ConnectedPlayerCapital == playerId1 {
				saveOutput.TileData[i][j].ImprovementData.ConnectedPlayerCapital = unusedPlayerId
			}

			if saveOutput.TileData[i][j].Unit != nil && saveOutput.TileData[i][j].Unit.Owner == uint8(playerId1) {
				saveOutput.TileData[i][j].Unit.Owner = uint8(unusedPlayerId)
			}

			if saveOutput.TileData[i][j].PassengerUnit != nil && saveOutput.TileData[i][j].PassengerUnit.Owner == uint8(playerId1) {
				saveOutput.TileData[i][j].PassengerUnit.Owner = uint8(unusedPlayerId)
			}
		}
	}

	// Overwrite all playerId2 tiles and units with playerId1
	for i := 0; i < saveOutput.MapHeight; i++ {
		for j := 0; j < saveOutput.MapWidth; j++ {
			if saveOutput.TileData[i][j].Owner == playerId2 {
				saveOutput.TileData[i][j].Owner = playerId1
			}

			if saveOutput.TileData[i][j].Capital == playerId2 {
				saveOutput.TileData[i][j].Capital = playerId1
			}
			if saveOutput.TileData[i][j].ImprovementData != nil && saveOutput.TileData[i][j].ImprovementData.ConnectedPlayerCapital == playerId2 {
				saveOutput.TileData[i][j].ImprovementData.ConnectedPlayerCapital = playerId1
			}

			if saveOutput.TileData[i][j].Unit != nil && saveOutput.TileData[i][j].Unit.Owner == uint8(playerId2) {
				saveOutput.TileData[i][j].Unit.Owner = uint8(playerId1)
			}

			if saveOutput.TileData[i][j].PassengerUnit != nil && saveOutput.TileData[i][j].PassengerUnit.Owner == uint8(playerId2) {
				saveOutput.TileData[i][j].PassengerUnit.Owner = uint8(playerId1)
			}
		}
	}

	// Overwrite old playerId tiles and units with playerId2
	for i := 0; i < saveOutput.MapHeight; i++ {
		for j := 0; j < saveOutput.MapWidth; j++ {
			if saveOutput.TileData[i][j].Owner == unusedPlayerId {
				saveOutput.TileData[i][j].Owner = playerId2
			}

			if saveOutput.TileData[i][j].Capital == unusedPlayerId {
				saveOutput.TileData[i][j].Capital = playerId2
			}
			if saveOutput.TileData[i][j].ImprovementData != nil && saveOutput.TileData[i][j].ImprovementData.ConnectedPlayerCapital == unusedPlayerId {
				saveOutput.TileData[i][j].ImprovementData.ConnectedPlayerCapital = playerId2
			}

			if saveOutput.TileData[i][j].Unit != nil && saveOutput.TileData[i][j].Unit.Owner == uint8(unusedPlayerId) {
				saveOutput.TileData[i][j].Unit.Owner = uint8(playerId2)
			}

			if saveOutput.TileData[i][j].PassengerUnit != nil && saveOutput.TileData[i][j].PassengerUnit.Owner == uint8(unusedPlayerId) {
				saveOutput.TileData[i][j].PassengerUnit.Owner = uint8(playerId2)
			}
		}
	}

	WriteMapToFile(inputFilename, saveOutput.TileData)

	var player1Tribe, player2Tribe int
	player1Color := make([]int, 4)
	player2Color := make([]int, 4)
	player1StartTile := [2]int{0, 0}
	player2StartTile := [2]int{0, 0}
	for i := 0; i < len(saveOutput.PlayerData); i++ {
		if saveOutput.PlayerData[i].Id == playerId1 {
			player1Tribe = saveOutput.PlayerData[i].Tribe
			copy(player1Color, saveOutput.PlayerData[i].OverrideColor)
			player1StartTile[0] = saveOutput.PlayerData[i].StartTileCoordinates[0]
			player1StartTile[1] = saveOutput.PlayerData[i].StartTileCoordinates[1]
		} else if saveOutput.PlayerData[i].Id == playerId2 {
			player2Tribe = saveOutput.PlayerData[i].Tribe
			copy(player2Color, saveOutput.PlayerData[i].OverrideColor)
			player2StartTile[0] = saveOutput.PlayerData[i].StartTileCoordinates[0]
			player2StartTile[1] = saveOutput.PlayerData[i].StartTileCoordinates[1]
		}
	}

	for i := 0; i < len(saveOutput.PlayerData); i++ {
		if saveOutput.PlayerData[i].Id == playerId1 {
			saveOutput.PlayerData[i].Tribe = player2Tribe
			saveOutput.PlayerData[i].OverrideColor = player2Color
			saveOutput.PlayerData[i].StartTileCoordinates[0] = player2StartTile[0]
			saveOutput.PlayerData[i].StartTileCoordinates[1] = player2StartTile[1]
		} else if saveOutput.PlayerData[i].Id == playerId2 {
			saveOutput.PlayerData[i].Tribe = player1Tribe
			saveOutput.PlayerData[i].OverrideColor = player1Color
			saveOutput.PlayerData[i].StartTileCoordinates[0] = player1StartTile[0]
			saveOutput.PlayerData[i].StartTileCoordinates[1] = player1StartTile[1]
		}
	}

	WritePlayersToFile(inputFilename, saveOutput.PlayerData)
}

func SetTileCapital(inputFilename string, targetX int, targetY int, newCityName string, updatedTribe int) {
	saveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}

	if updatedTribe >= 255 {
		log.Fatal("Value must be less than 255")
	}
	capitalTile := saveOutput.TileData[targetY][targetX]
	capitalTile.Capital = updatedTribe
	capitalTile.Owner = updatedTribe
	capitalTile.CapitalCoordinates[0] = capitalTile.WorldCoordinates[0]
	capitalTile.CapitalCoordinates[1] = capitalTile.WorldCoordinates[1]

	capitalTile.ImprovementExists = true
	capitalTile.ImprovementType = 1
	capitalTile.ImprovementData = &ImprovementData{
		Level:                  1,
		FoundedTurn:            0,
		CurrentPopulation:      0,
		TotalPopulation:        0,
		Production:             1,
		BaseScore:              0,
		BorderSize:             1,
		UpgradeCount:           0,
		ConnectedPlayerCapital: 0,
		HasCityName:            1,
		CityName:               newCityName,
		FoundedTribe:           0,
		CityRewards:            []int{},
		RebellionFlag:          0,
		RebellionBuffer:        []int{},
	}
	saveOutput.TileData[targetY][targetX] = capitalTile
	fmt.Println(fmt.Sprintf("Modified tile (%v, %v) to have capital %v", targetX, targetY, updatedTribe))

	for deltaX := -1; deltaX <= 1; deltaX++ {
		for deltaY := -1; deltaY <= 1; deltaY++ {
			if deltaX == 0 && deltaY == 0 {
				continue
			}

			neighborX := capitalTile.WorldCoordinates[0] + deltaX
			neighborY := capitalTile.WorldCoordinates[1] + deltaY

			if neighborX < 0 || neighborX >= saveOutput.MapWidth {
				continue
			}
			if neighborY < 0 || neighborY >= saveOutput.MapHeight {
				continue
			}

			saveOutput.TileData[neighborY][neighborX].Owner = updatedTribe
			saveOutput.TileData[neighborY][neighborX].CapitalCoordinates[0] = capitalTile.WorldCoordinates[0]
			saveOutput.TileData[neighborY][neighborX].CapitalCoordinates[1] = capitalTile.WorldCoordinates[1]
			fmt.Println(fmt.Sprintf("Set neighboring tile (%v, %v) to have owner %v", neighborX, neighborY, updatedTribe))
		}
	}
	WriteMapToFile(inputFilename, saveOutput.TileData)

	for i := 0; i < len(saveOutput.PlayerData); i++ {
		if saveOutput.PlayerData[i].Id == updatedTribe {
			saveOutput.PlayerData[i].StartTileCoordinates[0] = capitalTile.WorldCoordinates[0]
			saveOutput.PlayerData[i].StartTileCoordinates[1] = capitalTile.WorldCoordinates[1]
			WritePlayersToFile(inputFilename, saveOutput.PlayerData)
			fmt.Printf("Set player id %v start coordinates to (%v, %v)\n",
				saveOutput.PlayerData[i].Id, saveOutput.PlayerData[i].StartTileCoordinates[0], saveOutput.PlayerData[i].StartTileCoordinates[1])
			break
		}
	}
}
