package mapdraw

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"unicode/utf8"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/samuelyuan/PolytopiaMapEditor/internal/fileio"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	radius = 30.0
)

var (
	NeighborOffset = [4][2]int{{0, 1}, {-1, 0}, {0, -1}, {1, 0}}
)

func getImagePosition(i int, j int) (float64, float64) {
	x := float64(j) * (radius)
	y := float64(i) * (radius)
	return x, y
}

func getNeighbors(x int, y int) [4][2]int {
	offset := NeighborOffset

	neighbors := [4][2]int{}
	for i := 0; i < 4; i++ {
		newX := x + offset[i][0]
		newY := y + offset[i][1]
		neighbors[i][0] = newX
		neighbors[i][1] = newY
	}
	return neighbors
}

func getPhysicalMapTileColor(terrain int) color.RGBA {
	switch terrain {
	case 1: // Water
		return color.RGBA{95, 149, 149, 255}
	case 2: // Ocean
		return color.RGBA{47, 74, 93, 255}
	case 3, 4, 5: // flat land
		return color.RGBA{105, 125, 54, 255}
	case 6: // ice
		return color.RGBA{238, 249, 255, 255}
	}

	// default
	return color.RGBA{0, 0, 0, 255}
}

func getPoliticalMapTileColor(saveData *fileio.PolytopiaSaveOutput, row int, column int) color.RGBA {
	tileOwner := saveData.TileData[row][column].Owner
	tribe, ok := saveData.OwnerTribeMap[tileOwner]
	if !ok {
		// default
		return color.RGBA{0, 0, 0, 255}
	}

	for i := 0; i < len(saveData.PlayerData); i++ {
		playerData := saveData.PlayerData[i]
		if playerData.Id == tileOwner {
			// override color
			if playerData.OverrideColor[3] != 255 {
				return color.RGBA{uint8(playerData.OverrideColor[2]), uint8(playerData.OverrideColor[1]), uint8(playerData.OverrideColor[0]), 255}
			}
		}
	}

	switch tribe {
	case 2: // Ai-Mo
		return color.RGBA{54, 226, 170, 255}
	case 3: // Aquarion
		return color.RGBA{243, 131, 129, 255}
	case 4: // Bardur
		return color.RGBA{53, 37, 20, 255}
	case 5: // Elyrion
		return color.RGBA{255, 0, 153, 255}
	case 6: // Hoodrick
		return color.RGBA{153, 102, 0, 255}
	case 7: // Imperius
		return color.RGBA{0, 0, 255, 255}
	case 8: // Kickoo
		return color.RGBA{0, 255, 0, 255}
	case 9: // Luxidoor
		return color.RGBA{171, 59, 214, 255}
	case 10: // Oumaji
		return color.RGBA{255, 255, 0, 255}
	case 11: // Quetzali
		return color.RGBA{39, 92, 74, 255}
	case 12: // Vengir
		return color.RGBA{255, 255, 255, 255}
	case 13: // Xin-xi
		return color.RGBA{204, 0, 0, 255}
	case 14: // Yadakk
		return color.RGBA{125, 35, 28, 255}
	case 15: // Zebasi
		return color.RGBA{255, 153, 0, 255}
	case 16: // Polaris
		return color.RGBA{182, 161, 133, 255}
	case 17: // Cymanti
		return color.RGBA{194, 253, 0, 255}
	}

	return color.RGBA{128, 128, 128, 255}
}

func drawCityIcon(dc *gg.Context, imageX float64, imageY float64, cityColor color.RGBA) {
	iconColor := cityColor
	dc.DrawRectangle(imageX+(radius/4), imageY+(radius/4), radius/2, radius/2)
	dc.SetRGB255(int(iconColor.R), int(iconColor.G), int(iconColor.B))
	dc.Fill()
}

func drawMountain(dc *gg.Context, imageX float64, imageY float64) {
	// draw base
	dc.DrawRegularPolygon(3, imageX+(radius/2), imageY+(radius/3), radius/2, math.Pi)
	dc.SetRGB255(89, 90, 86) // gray
	dc.Fill()

	// draw mountain peak
	dc.DrawRegularPolygon(3, imageX+(radius/2), imageY+(radius*2/3), radius/4, math.Pi)
	dc.SetRGB255(234, 244, 253) // white
	dc.Fill()
}

func drawForest(dc *gg.Context, imageX float64, imageY float64) {
	dc.DrawRegularPolygon(3, imageX+(radius/2), imageY+(radius/3), radius/2, math.Pi)
	dc.SetRGB255(53, 72, 44) // dark green
	dc.Fill()
}

func drawTerritoryTiles(dc *gg.Context, saveData *fileio.PolytopiaSaveOutput, mapHeight int, mapWidth int) {
	for i := 0; i < mapHeight; i++ {
		for j := 0; j < mapWidth; j++ {
			x, y := getImagePosition(i, j)

			dc.DrawRectangle(x, y, radius, radius)
			tileData := saveData.TileData[i][j]
			terrain := tileData.Terrain

			terrainTileColor := getPhysicalMapTileColor(terrain)
			dc.SetRGB255(int(terrainTileColor.R), int(terrainTileColor.G), int(terrainTileColor.B))
			dc.Fill()

			if terrain == 4 {
				drawMountain(dc, x, y)
			} else if terrain == 5 {
				drawForest(dc, x, y)
			}

			// Draw cities
			if tileData.ImprovementData != nil && tileData.ImprovementType == 1 {
				if tileData.Owner > 0 {
					// Capital city
					cityColor := getPoliticalMapTileColor(saveData, i, j)
					drawCityIcon(dc, x, y, cityColor)
				} else {
					// Village
					drawCityIcon(dc, x, y, color.RGBA{255, 255, 255, 255})
				}
			}
		}
	}
}

func drawBorders(dc *gg.Context, saveData *fileio.PolytopiaSaveOutput, mapHeight int, mapWidth int) {
	for i := 0; i < mapHeight; i++ {
		for j := 0; j < mapWidth; j++ {
			x1, y1 := getImagePosition(i, j)
			neighbors := getNeighbors(j, i)
			currentTileOwner := saveData.TileData[i][j].Owner
			if currentTileOwner == 0 {
				continue
			}

			tileColor := getPoliticalMapTileColor(saveData, i, j)
			lineWidth := 1.5
			for n := 0; n < len(neighbors); n++ {
				newX := neighbors[n][0]
				newY := neighbors[n][1]
				if newX >= 0 && newY >= 0 && newX < mapWidth && newY < mapHeight {
					otherTileOwner := saveData.TileData[newY][newX].Owner
					if currentTileOwner != otherTileOwner {
						angle1 := (math.Pi / 4) + float64(n)*(math.Pi/2)
						angle2 := (math.Pi / 4) + float64(n+1)*(math.Pi/2)

						centerX := x1 + (radius / 2)
						centerY := y1 + (radius / 2)

						edgeX1 := centerX + ((radius-1)*math.Sqrt2/2)*math.Cos(angle1)
						edgeY1 := centerY + ((radius-1)*math.Sqrt2/2)*math.Sin(angle1)
						edgeX2 := centerX + ((radius-1)*math.Sqrt2/2)*math.Cos(angle2)
						edgeY2 := centerY + ((radius-1)*math.Sqrt2/2)*math.Sin(angle2)

						dc.SetRGB255(int(tileColor.R), int(tileColor.G), int(tileColor.B))
						dc.SetLineWidth(lineWidth)
						dc.DrawLine(edgeX1, edgeY1, edgeX2, edgeY2)
						dc.Stroke()
					}
				}
			}
		}
	}
	dc.SetLineWidth(1.0)
}

func drawCityNames(dc *gg.Context, saveData *fileio.PolytopiaSaveOutput, mapHeight int, mapWidth int) {
	dc.SetRGB255(255, 255, 255)
	for i := 0; i < mapHeight; i++ {
		for j := 0; j < mapWidth; j++ {
			// Invert depth because the map is inverted
			x, y := getImagePosition(mapHeight-i, j)

			tile := saveData.TileData[i][j]
			cityName := ""
			if tile.ImprovementData != nil && tile.ImprovementType == 1 {
				cityName = tile.ImprovementData.CityName
			}
			if len(cityName) == 0 {
				continue
			}

			if utf8.RuneCountInString(cityName) <= 4 {
				dc.DrawString(cityName, x, y-radius)
			} else {
				dc.DrawString(cityName, x-(2.0*float64(len(cityName))/2.0), y-radius)
			}
		}
	}
}

func DrawMap(saveData *fileio.PolytopiaSaveOutput, highlightedTileX int, highlightedTileY int) image.Image {
	mapHeight := saveData.MapHeight
	mapWidth := saveData.MapWidth

	maxImageWidth, maxImageHeight := getImagePosition(mapHeight, mapWidth)
	dc := gg.NewContext(int(maxImageWidth), int(maxImageHeight))

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	face := truetype.NewFace(font, &truetype.Options{Size: 14})
	dc.SetFontFace(face)

	// Need to invert image because the map format is inverted
	dc.InvertY()

	drawTerritoryTiles(dc, saveData, mapHeight, mapWidth)
	drawBorders(dc, saveData, mapHeight, mapWidth)

	// draw highlighted tile
	if highlightedTileX != -1 && highlightedTileY != -1 {
		dc.SetRGB255(0, 0, 0)
		x, y := getImagePosition(highlightedTileY, highlightedTileX)
		dc.SetLineWidth(2.0)
		dc.DrawRectangle(x, y, radius, radius)
		dc.Stroke()
	}

	dc.InvertY()

	// Draw city names after inversion
	drawCityNames(dc, saveData, mapHeight, mapWidth)

	return dc.Image()
}

func SaveImage(outputFilename string, im image.Image) {
	gg.SavePNG(outputFilename, im)
	fmt.Println("Saved image to", outputFilename)
}

func GetTileCoordinates(pixelX int, pixelY int, mapHeight int, mapWidth int) (int, int) {
	tileX := -1
	tileY := -1
	for i := 0; i < mapHeight; i++ {
		for j := 0; j < mapWidth; j++ {
			tileTopLeftCornerX, tileTopLeftCornerY := getImagePosition(i, j)
			tileBottomRightCornerX, tileBottomRightCornerY := getImagePosition(i+1, j+1)

			if float64(pixelX) >= tileTopLeftCornerX && float64(pixelX) <= tileBottomRightCornerX &&
				float64(pixelY) >= tileTopLeftCornerY && float64(pixelY) <= tileBottomRightCornerY {
				tileX = j
				tileY = (mapHeight - 1) - i
				break
			}
		}
	}

	return tileX, tileY
}
