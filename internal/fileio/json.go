package fileio

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type PolytopiaSaveJson struct {
	GameName        string
	FileFormat      string
	TileData        [][]TileData
	PlayerData      []PlayerData
	MapHeaderOutput MapHeaderOutput
}

func ImportPolytopiaDataFromJson(inputFilename string) *PolytopiaSaveJson {
	jsonFile, err := os.Open(inputFilename)
	if err != nil {
		log.Fatal("Failed to open json file", err)
	}
	defer jsonFile.Close()

	jsonContents, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	var polytopiaSaveJson *PolytopiaSaveJson
	json.Unmarshal(jsonContents, &polytopiaSaveJson)

	if polytopiaSaveJson == nil {
		log.Fatal("The json data in " + inputFilename + " is missing or incorrect")
	}

	return polytopiaSaveJson
}

func ExportPolytopiaJsonFile(saveOutput *PolytopiaSaveOutput, outputFilename string) {
	polytopiaJson := &PolytopiaSaveJson{
		GameName:        "Battle of Polytopia",
		FileFormat:      "Polytopia Save State",
		TileData:        saveOutput.TileData,
		PlayerData:      saveOutput.PlayerData,
		MapHeaderOutput: saveOutput.MapHeaderOutput,
	}

	file, err := json.MarshalIndent(polytopiaJson, "", " ")
	if err != nil {
		log.Fatal("Failed to marshal data: ", err)
	}

	err = ioutil.WriteFile(outputFilename, file, 0644)
	if err != nil {
		log.Fatal("Error writing to ", outputFilename)
	}
}
