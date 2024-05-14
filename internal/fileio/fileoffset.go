package fileio

import (
	"fmt"
	"io"
	"log"
)

func buildMapStartKey() string {
	return "MapStart"
}

func buildMapEndKey() string {
	return "MapEnd"
}

func buildTileStartKey(x int, y int) string {
	return fmt.Sprintf("TileStart%v,%v", x, y)
}

func buildTileEndKey(x int, y int) string {
	return fmt.Sprintf("TileEnd%v,%v", x, y)
}

func buildAllPlayersStartKey() string {
	return "AllPlayersStart"
}

func buildAllPlayersEndKey() string {
	return "AllPlayersEnd"
}

func buildPlayerStartKey(index int) string {
	return fmt.Sprintf("PlayerStart%v", index)
}

func buildPlayerArr1Key(playerId int) string {
	return fmt.Sprintf("PlayerArr1-Id%v", playerId)
}

func buildPlayerCurrencyKey(playerId int) string {
	return fmt.Sprintf("PlayerCurrency-Id%v", playerId)
}

func updateFileOffsetMap(fileOffsetMap map[string]int, streamReader *io.SectionReader, unitLocationKey string) {
	fileOffset, err := streamReader.Seek(0, io.SeekCurrent)
	if err != nil {
		log.Fatal(err)
	}
	fileOffsetMap[unitLocationKey] = int(fileOffset)
}
