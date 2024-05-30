package fileio

import (
	"io"
)

type ImprovementData struct {
	Level                  int
	FoundedTurn            int
	CurrentPopulation      int
	TotalPopulation        int
	Production             int
	BaseScore              int
	BorderSize             int // For cities, 1 is default, 2 is expanded border
	UpgradeCount           int // For cities, seems to be -1 * (level - 1). Level 1 is starting point and no upgrades.
	ConnectedPlayerCapital int
	HasCityName            int
	CityName               string
	FoundedTribe           int
	CityRewards            []int
	RebellionFlag          int
	RebellionBuffer        []int
}

func DeserializeImprovementDataFromBytes(streamReader *io.SectionReader) ImprovementData {
	level := unsafeReadUint16(streamReader)
	foundedTurn := unsafeReadUint16(streamReader)
	currentPopulation := unsafeReadInt16(streamReader)
	totalPopulation := unsafeReadUint16(streamReader)
	production := unsafeReadInt16(streamReader)
	baseScore := unsafeReadInt16(streamReader)
	borderSize := unsafeReadInt16(streamReader)
	upgradeCount := unsafeReadInt16(streamReader)
	connectedPlayerCapital := unsafeReadUint8(streamReader)
	hasCityName := unsafeReadUint8(streamReader)
	cityName := ""
	if hasCityName == 1 {
		cityName = readVarString(streamReader, "CityName")
	}

	foundedTribe := unsafeReadUint8(streamReader)

	cityRewardsSize := unsafeReadUint16(streamReader)
	cityRewards := make([]int, cityRewardsSize)
	for i := 0; i < int(cityRewardsSize); i++ {
		cityReward := unsafeReadUint16(streamReader)
		cityRewards[i] = int(cityReward)
	}

	rebellionFlag := unsafeReadUint16(streamReader)
	rebellionBuffer := []int{}
	if rebellionFlag != 0 {
		rebellionBuffer = convertByteListToInt(readFixedList(streamReader, 2))
	}

	return ImprovementData{
		Level:                  int(level),
		FoundedTurn:            int(foundedTurn),
		CurrentPopulation:      int(currentPopulation),
		TotalPopulation:        int(totalPopulation),
		Production:             int(production),
		BaseScore:              int(baseScore),
		BorderSize:             int(borderSize),
		UpgradeCount:           int(upgradeCount),
		ConnectedPlayerCapital: int(connectedPlayerCapital),
		HasCityName:            int(hasCityName),
		CityName:               cityName,
		FoundedTribe:           int(foundedTribe),
		CityRewards:            cityRewards,
		RebellionFlag:          int(rebellionFlag),
		RebellionBuffer:        rebellionBuffer,
	}
}

func SerializeImprovementDataToBytes(improvementData ImprovementData) []byte {
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
