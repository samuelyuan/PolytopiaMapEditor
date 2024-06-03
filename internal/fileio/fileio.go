package fileio

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	fileOffsetMap = make(map[string]int)
)

type PolytopiaSaveOutput struct {
	MapHeight         int
	MapWidth          int
	MapHeaderOutput   MapHeaderOutput
	InitialTileData   [][]TileData
	InitialPlayerData []PlayerData
	TileData          [][]TileData
	MaxTurn           int
	PlayerData        []PlayerData
	FileOffsetMap     map[string]int
	TribeCityMap      map[int][]CityLocationData
}

type CityLocationData struct {
	X        int
	Y        int
	CityName string
	Capital  int
}

type UnitLocationData struct {
	X        int
	Y        int
	UnitType int
}

func readVarString(reader *io.SectionReader, varName string) string {
	variableLength := uint8(0)
	if err := binary.Read(reader, binary.LittleEndian, &variableLength); err != nil {
		log.Fatal("Failed to load variable length: ", err)
	}

	stringValue := make([]byte, variableLength)
	if err := binary.Read(reader, binary.LittleEndian, &stringValue); err != nil {
		log.Fatal(fmt.Sprintf("Failed to load string value. Variable length: %v, name: %s. Error:", variableLength, varName), err)
	}

	return string(stringValue[:])
}

func unsafeReadUint32(reader *io.SectionReader) uint32 {
	unsignedIntValue := uint32(0)
	if err := binary.Read(reader, binary.LittleEndian, &unsignedIntValue); err != nil {
		log.Fatal("Failed to load uint32: ", err)
	}
	return unsignedIntValue
}

func unsafeReadInt32(reader *io.SectionReader) int32 {
	signedIntValue := int32(0)
	if err := binary.Read(reader, binary.LittleEndian, &signedIntValue); err != nil {
		log.Fatal("Failed to load int32: ", err)
	}
	return signedIntValue
}

func unsafeReadUint16(reader *io.SectionReader) uint16 {
	unsignedIntValue := uint16(0)
	if err := binary.Read(reader, binary.LittleEndian, &unsignedIntValue); err != nil {
		log.Fatal("Failed to load uint16: ", err)
	}
	return unsignedIntValue
}

func unsafeReadInt16(reader *io.SectionReader) int16 {
	signedIntValue := int16(0)
	if err := binary.Read(reader, binary.LittleEndian, &signedIntValue); err != nil {
		log.Fatal("Failed to load int16: ", err)
	}
	return signedIntValue
}

func unsafeReadUint8(reader *io.SectionReader) uint8 {
	unsignedIntValue := uint8(0)
	if err := binary.Read(reader, binary.LittleEndian, &unsignedIntValue); err != nil {
		log.Fatal("Failed to load uint8: ", err)
	}
	return unsignedIntValue
}

func unsafeReadFloat32(reader *io.SectionReader) float32 {
	floatValue := float32(0)
	if err := binary.Read(reader, binary.LittleEndian, &floatValue); err != nil {
		log.Fatal("Failed to load float32: ", err)
	}
	return floatValue
}

func readFixedList(streamReader *io.SectionReader, listSize int) []byte {
	buffer := make([]byte, listSize)
	if err := binary.Read(streamReader, binary.LittleEndian, &buffer); err != nil {
		log.Fatal("Failed to load buffer: ", err)
	}
	return buffer
}

func convertByteListToInt(oldArr []byte) []int {
	newArr := make([]int, len(oldArr))
	for i := 0; i < len(newArr); i++ {
		newArr[i] = int(oldArr[i])
	}
	return newArr
}

func readTileData(streamReader *io.SectionReader, tileData [][]TileData, mapWidth int, mapHeight int) {
	updateFileOffsetMap(fileOffsetMap, streamReader, buildMapStartKey())

	for i := 0; i < int(mapHeight); i++ {
		for j := 0; j < int(mapWidth); j++ {
			tileStartKey := buildTileStartKey(j, i)
			updateFileOffsetMap(fileOffsetMap, streamReader, tileStartKey)

			tileData[i][j] = DeserializeTileDataFromBytes(streamReader, i, j)

			tileEndKey := buildTileEndKey(j, i)
			updateFileOffsetMap(fileOffsetMap, streamReader, tileEndKey)
		}
	}

	updateFileOffsetMap(fileOffsetMap, streamReader, buildMapEndKey())
}

func readAllPlayerData(streamReader *io.SectionReader) []PlayerData {
	allPlayersStartKey := buildAllPlayersStartKey()
	updateFileOffsetMap(fileOffsetMap, streamReader, allPlayersStartKey)

	numPlayers := unsafeReadUint16(streamReader)
	allPlayerData := make([]PlayerData, int(numPlayers))

	for i := 0; i < int(numPlayers); i++ {
		playerStartKey := buildPlayerStartKey(i)
		updateFileOffsetMap(fileOffsetMap, streamReader, playerStartKey)
		playerData := DeserializePlayerDataFromBytes(streamReader)
		allPlayerData[i] = playerData
	}

	allPlayersEndKey := buildAllPlayersEndKey()
	updateFileOffsetMap(fileOffsetMap, streamReader, allPlayersEndKey)

	return allPlayerData
}

func buildTribeCityMap(currentMapHeaderOutput MapHeaderOutput, tileData [][]TileData) map[int][]CityLocationData {
	tribeCityMap := make(map[int][]CityLocationData)
	for i := 0; i < int(currentMapHeaderOutput.MapHeight); i++ {
		for j := 0; j < int(currentMapHeaderOutput.MapWidth); j++ {
			if tileData[i][j].ImprovementData != nil && tileData[i][j].ImprovementType == 1 {
				tribeOwner := tileData[i][j].Owner
				_, ok := tribeCityMap[tribeOwner]
				if !ok {
					tribeCityMap[tribeOwner] = make([]CityLocationData, 0)
				}

				cityName := ""
				if tileData[i][j].ImprovementData != nil {
					cityName = tileData[i][j].ImprovementData.CityName
				}

				cityLocationData := CityLocationData{
					X:        tileData[i][j].WorldCoordinates[0],
					Y:        tileData[i][j].WorldCoordinates[1],
					CityName: cityName,
					Capital:  tileData[i][j].Capital,
				}
				tribeCityMap[tribeOwner] = append(tribeCityMap[tribeOwner], cityLocationData)
			}
		}
	}
	return tribeCityMap
}

func ReadPolytopiaDecompressedFile(inputFilename string) (*PolytopiaSaveOutput, error) {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state: ", err)
		return nil, err
	}
	fi, err := inputFile.Stat()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	fileLength := fi.Size()
	streamReader := io.NewSectionReader(inputFile, int64(0), fileLength)

	fileOffsetMap = make(map[string]int)

	// Read initial map state
	updateFileOffsetMap(fileOffsetMap, streamReader, buildMapHeaderStartKey())
	initialMapHeaderOutput := DeserializeMapHeaderFromBytes(streamReader)
	updateFileOffsetMap(fileOffsetMap, streamReader, buildMapHeaderEndKey())

	initialTileData := make([][]TileData, initialMapHeaderOutput.MapHeight)
	for i := 0; i < initialMapHeaderOutput.MapHeight; i++ {
		initialTileData[i] = make([]TileData, initialMapHeaderOutput.MapWidth)
	}
	readTileData(streamReader, initialTileData, initialMapHeaderOutput.MapWidth, initialMapHeaderOutput.MapHeight)
	initialPlayerData := readAllPlayerData(streamReader)

	_ = readFixedList(streamReader, 3)

	// Read current map state
	updateFileOffsetMap(fileOffsetMap, streamReader, buildMapHeaderStartKey())
	currentMapHeaderOutput := DeserializeMapHeaderFromBytes(streamReader)
	updateFileOffsetMap(fileOffsetMap, streamReader, buildMapHeaderEndKey())

	tileData := make([][]TileData, currentMapHeaderOutput.MapHeight)
	for i := 0; i < currentMapHeaderOutput.MapHeight; i++ {
		tileData[i] = make([]TileData, currentMapHeaderOutput.MapWidth)
	}
	readTileData(streamReader, tileData, currentMapHeaderOutput.MapWidth, currentMapHeaderOutput.MapHeight)
	playerData := readAllPlayerData(streamReader)

	tribeCityMap := buildTribeCityMap(currentMapHeaderOutput, tileData)

	output := &PolytopiaSaveOutput{
		MapHeight:         currentMapHeaderOutput.MapHeight,
		MapWidth:          currentMapHeaderOutput.MapWidth,
		MapHeaderOutput:   currentMapHeaderOutput,
		InitialTileData:   initialTileData,
		InitialPlayerData: initialPlayerData,
		TileData:          tileData,
		MaxTurn:           int(currentMapHeaderOutput.MapHeaderInput.CurrentTurn),
		PlayerData:        playerData,
		FileOffsetMap:     fileOffsetMap,
		TribeCityMap:      tribeCityMap,
	}
	return output, nil
}
