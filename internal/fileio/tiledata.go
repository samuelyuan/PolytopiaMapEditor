package fileio

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

type TileDataHeader struct {
	WorldCoordinates   [2]uint32
	Terrain            uint16
	Climate            uint16
	Altitude           int16
	Owner              uint8
	Capital            uint8
	CapitalCoordinates [2]int32
}

type TileData struct {
	WorldCoordinates           [2]int
	Terrain                    int
	Climate                    int
	Altitude                   int
	Owner                      int
	Capital                    int
	CapitalCoordinates         [2]int
	ResourceExists             bool
	ResourceType               int
	ImprovementExists          bool
	ImprovementType            int
	ImprovementData            *ImprovementData
	Unit                       *UnitData
	PassengerUnit              *UnitData
	UnitEffectData             []int // flags: 0 - ice, 1 - poison, 2 - boost, 3 - invisible
	UnitDirectionData          []int // contains direction flag (0 - southwest, 1 - west, 2 - northwest, 3 - north, 4 - northeast, 5 - east, 6 - southwest, 7 - south)
	PassengerUnitEffectData    []int
	PassengerUnitDirectionData []int
	PlayerVisibility           []int
	HasRoad                    bool
	HasWaterRoute              bool
	TileSkin                   int
	Unknown                    []int
	FloodedFlag                int // introduced in new aquarion update (version 105)
	FloodedValue               int // introduced in new aquarion update (version 105)
}

type UnitData struct {
	Id                 uint32
	Owner              uint8
	UnitType           uint16
	FollowerUnitId     uint32 // only initialized for cymanti centipedes and segments
	LeaderUnitId       uint32 // only initialized for cymanti centipedes and segments
	CurrentCoordinates [2]int32
	HomeCoordinates    [2]int32
	Health             uint16 // should be divided by 10 to get value ingame
	PromotionLevel     uint16
	Experience         uint16
	Moved              bool
	Attacked           bool
	Flipped            bool
	CreatedTurn        uint16
}

func DeserializeTileDataFromBytes(streamReader *io.SectionReader, expectedRow int, expectedCol int, gameVersion int) TileData {
	tileDataHeader := TileDataHeader{}
	if err := binary.Read(streamReader, binary.LittleEndian, &tileDataHeader); err != nil {
		log.Fatal("Failed to load tileDataHeader: ", err)
	}

	// Sanity check
	if int(tileDataHeader.WorldCoordinates[0]) != expectedCol || int(tileDataHeader.WorldCoordinates[1]) != expectedRow {
		log.Fatal(fmt.Sprintf("File reached unexpected location. Iteration (%v, %v) isn't equal to world coordinates (%v, %v)",
			expectedRow, expectedCol, tileDataHeader.WorldCoordinates[0], tileDataHeader.WorldCoordinates[1]))
	}

	resourceExistsFlag := unsafeReadUint8(streamReader)
	resourceType := -1
	if resourceExistsFlag == 1 {
		resourceType = int(unsafeReadUint16(streamReader))
	}

	improvementExistsFlag := unsafeReadUint8(streamReader)
	improvementType := -1
	if improvementExistsFlag == 1 {
		improvementType = int(unsafeReadUint16(streamReader))
	}

	var improvementDataPtr *ImprovementData
	if improvementType != -1 {
		improvementData := DeserializeImprovementDataFromBytes(streamReader)
		improvementDataPtr = &improvementData
	}

	// Read unit data
	hasUnitFlag := unsafeReadUint8(streamReader)
	var unitDataPtr *UnitData
	var passengerUnitDataPtr *UnitData

	unitEffectData := make([]int, 0)
	unitDirectionData := make([]byte, 0)
	passengerUnitEffectData := make([]int, 0)
	passengerUnitDirectionData := make([]byte, 0)
	if hasUnitFlag == 1 {
		unitData := UnitData{}
		if err := binary.Read(streamReader, binary.LittleEndian, &unitData); err != nil {
			log.Fatal("Failed to load buffer: ", err)
		}
		unitDataPtr = &unitData

		hasOtherUnitFlag := unsafeReadUint8(streamReader)
		if hasOtherUnitFlag == 1 {
			// If unit embarks or disembarks, a new unit is created in the backend, but it's still the same unit in the game
			passengerUnitData := UnitData{}
			if err := binary.Read(streamReader, binary.LittleEndian, &passengerUnitData); err != nil {
				log.Fatal("Failed to load buffer: ", err)
			}
			passengerUnitDataPtr = &passengerUnitData

			// might be other unit flag for passenger unit
			// should always be zero because passenger unit can't carry another unit
			unknownFlag := int(unsafeReadUint8(streamReader))
			if unknownFlag != 0 {
				log.Fatal("Passenger unit's other unit flag isn't zero")
			}

			passengerUnitEffectCount := int(unsafeReadUint16(streamReader))
			passengerUnitEffectData = make([]int, 0)
			for statusIndex := 0; statusIndex < passengerUnitEffectCount; statusIndex++ {
				passengerUnitEffectData = append(passengerUnitEffectData, int(unsafeReadUint16(streamReader)))
			}
			passengerUnitDirectionData = readFixedList(streamReader, 5)

			unitEffectCount := int(unsafeReadUint16(streamReader))
			unitEffectData = make([]int, 0)
			for statusIndex := 0; statusIndex < unitEffectCount; statusIndex++ {
				unitEffectData = append(unitEffectData, int(unsafeReadUint16(streamReader)))
			}
			unitDirectionData = readFixedList(streamReader, 5)
		} else {
			unitEffectCount := int(unsafeReadUint16(streamReader))
			unitEffectData = make([]int, 0)
			for statusIndex := 0; statusIndex < unitEffectCount; statusIndex++ {
				unitEffectData = append(unitEffectData, int(unsafeReadUint16(streamReader)))
			}
			unitDirectionData = readFixedList(streamReader, 5)
		}
	}

	playerVisibilityListSize := unsafeReadUint8(streamReader)
	playerVisibilityList := convertByteListToInt(readFixedList(streamReader, int(playerVisibilityListSize)))
	hasRoad := unsafeReadUint8(streamReader)
	hasWaterRoute := unsafeReadUint8(streamReader)
	tileSkin := unsafeReadUint16(streamReader)
	unknown := convertByteListToInt(readFixedList(streamReader, 2))
	var floodedFlag int
	var floodedValue int
	if gameVersion >= 105 {
		floodedFlag = int(unsafeReadUint8(streamReader))
		if floodedFlag == 1 {
			floodedValue = int(unsafeReadUint32(streamReader))
		}
	}

	return TileData{
		WorldCoordinates:           [2]int{int(tileDataHeader.WorldCoordinates[0]), int(tileDataHeader.WorldCoordinates[1])},
		Terrain:                    int(tileDataHeader.Terrain),
		Climate:                    int(tileDataHeader.Climate),
		Altitude:                   int(tileDataHeader.Altitude),
		Owner:                      int(tileDataHeader.Owner),
		Capital:                    int(tileDataHeader.Capital),
		CapitalCoordinates:         [2]int{int(tileDataHeader.CapitalCoordinates[0]), int(tileDataHeader.CapitalCoordinates[1])},
		ResourceExists:             resourceExistsFlag != 0,
		ResourceType:               resourceType,
		ImprovementExists:          improvementExistsFlag != 0,
		ImprovementType:            improvementType,
		ImprovementData:            improvementDataPtr,
		Unit:                       unitDataPtr,
		PassengerUnit:              passengerUnitDataPtr,
		UnitEffectData:             unitEffectData,
		UnitDirectionData:          convertByteListToInt(unitDirectionData),
		PassengerUnitEffectData:    passengerUnitEffectData,
		PassengerUnitDirectionData: convertByteListToInt(passengerUnitDirectionData),
		PlayerVisibility:           playerVisibilityList,
		HasRoad:                    hasRoad != 0,
		HasWaterRoute:              hasWaterRoute != 0,
		TileSkin:                   int(tileSkin),
		Unknown:                    unknown,
		FloodedFlag:                floodedFlag,
		FloodedValue:               floodedValue,
	}
}

func SerializeTileToBytes(tileData TileData, gameVersion int) []byte {
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
		tileBytes = append(tileBytes, SerializeImprovementDataToBytes(*tileData.ImprovementData)...)
	}

	// no unit
	if tileData.Unit != nil {
		tileBytes = append(tileBytes, 1)
		tileBytes = append(tileBytes, SerializeUnitDataToBytes(*tileData.Unit)...)

		if tileData.PassengerUnit != nil {
			tileBytes = append(tileBytes, 1) // has other unit flag is 1
			tileBytes = append(tileBytes, SerializeUnitDataToBytes(*tileData.PassengerUnit)...)

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
	tileBytes = append(tileBytes, ConvertUint16Bytes(tileData.TileSkin)...)
	tileBytes = append(tileBytes, ConvertByteList(tileData.Unknown)...)
	if gameVersion >= 105 {
		tileBytes = append(tileBytes, byte(tileData.FloodedFlag))
		if tileData.FloodedFlag == 1 {
			tileBytes = append(tileBytes, ConvertUint32Bytes(tileData.FloodedValue)...)
		}
	}
	return tileBytes
}

func SerializeUnitDataToBytes(unitData UnitData) []byte {
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
