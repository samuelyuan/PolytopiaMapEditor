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

func buildMapHeaderStartKey() string {
	return "MapHeaderStart"
}

func buildMapHeaderEndKey() string {
	return "MapHeaderEnd"
}

func updateFileOffsetMap(fileOffsetMap map[string]int, streamReader *io.SectionReader, unitLocationKey string) {
	fileOffset, err := streamReader.Seek(0, io.SeekCurrent)
	if err != nil {
		log.Fatal(err)
	}
	fileOffsetMap[unitLocationKey] = int(fileOffset)
}
