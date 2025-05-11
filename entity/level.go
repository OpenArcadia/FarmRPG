package entity

import (
	"fmt"
	"strings"

	"com.openarcadia.farmrpg/utils"
	"github.com/akashKarmakar02/tmx"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type LevelData struct {
	TextureID int
	X         int
	Y         int
	Z         int
	Height    int
	Width     int
	TileX     int
	TileY     int
}

type Level struct {
	BackgroundTexture *rl.Texture2D
	MapTextures       []*LevelData
	TextureCache      map[int]*rl.Texture2D
	Water             *Water
	WaterLocations    []*WaterLocation
}

func NewLevel() *Level {
	levelAssetPath := utils.ImportAssetPath("world/ground.png")
	levelAsset := rl.LoadTexture(levelAssetPath)

	levelTexture := &levelAsset

	tileMap, err := tmx.LoadFile("data", "map.tmx")
	if err != nil {
		fmt.Println("failed to decode filemap")
	}

	if tileMap == nil {
		fmt.Println("tilemap is empty")
	}

	textures := make([]*LevelData, 0)
	layers := []string{"Water", "HouseFloor", "HouseFurnitureBottom"}
	textureMap := map[int]*rl.Texture2D{}
	texturePathMap := map[string]int{}
	waterLocations := make([]*WaterLocation, 0)

	for _, layer := range layers {
		if houseBottom := tileMap.LayerWithName(layer); houseBottom != nil {
			tiles, _ := houseBottom.TileDefs(tileMap.TileSets)

			for i, t := range tiles {
				if t.GlobalID != 0 {
					parts := strings.Split(t.TileSet.Image.Source, "/")
					joined := strings.Join(parts[2:], "/")

					tileIndexInSet := t.ID - t.TileSet.FirstGlobalID.TileID(t.TileSet)
					columns := t.TileSet.Columns
					srcX := (int(tileIndexInSet) % columns) * t.TileSet.TileWidth
					srcY := (int(tileIndexInSet) / columns) * t.TileSet.TileHeight

					mapX := (i % tileMap.Width) * tileMap.TileWidth
					mapY := (i / tileMap.Width) * tileMap.TileHeight

					if layer == "Water" {
						waterLocations = append(waterLocations, &WaterLocation{
							X: mapX,
							Y: mapY,
						})
						continue
					}

					if texturePathMap[joined] == 0 {
						texture := rl.LoadTexture(joined)
						texturePathMap[joined] = int(texture.ID)
						textureMap[int(texture.ID)] = &texture
						textures = append(textures, &LevelData{
							TextureID: int(texture.ID),
							X:         mapX,
							Y:         mapY,
							Z:         0,
							TileX:     srcX,
							TileY:     srcY,
							Height:    64,
							Width:     64,
						})
					} else {
						textures = append(textures, &LevelData{
							TextureID: texturePathMap[joined],
							X:         mapX,
							Y:         mapY,
							Z:         0,
							Height:    64,
							Width:     64,
							TileX:     srcX,
							TileY:     srcY,
						})
					}

				}
			}
		}
	}

	layers = []string{"HouseWalls", "HouseFurnitureTop", "Fence"}
	for _, layer := range layers {
		if houseBottom := tileMap.LayerWithName(layer); houseBottom != nil {
			tiles, _ := houseBottom.TileDefs(tileMap.TileSets)

			for i, t := range tiles {
				if t.GlobalID != 0 {
					parts := strings.Split(t.TileSet.Image.Source, "/")
					joined := strings.Join(parts[2:], "/")

					tileIndexInSet := t.ID - t.TileSet.FirstGlobalID.TileID(t.TileSet)
					columns := t.TileSet.Columns
					srcX := (int(tileIndexInSet) % columns) * t.TileSet.TileWidth
					srcY := (int(tileIndexInSet) / columns) * t.TileSet.TileHeight

					mapX := (i % tileMap.Width) * tileMap.TileWidth
					mapY := (i / tileMap.Width) * tileMap.TileHeight

					if texturePathMap[joined] == 0 {
						texture := rl.LoadTexture(joined)
						texturePathMap[joined] = int(texture.ID)
						textureMap[int(texture.ID)] = &texture
						textures = append(textures, &LevelData{
							TextureID: int(texture.ID),
							X:         mapX,
							Y:         mapY,
							Z:         2,
							Height:    64,
							Width:     64,
							TileX:     srcX,
							TileY:     srcY,
						})
					} else {
						textures = append(textures, &LevelData{
							TextureID: texturePathMap[joined],
							X:         mapX,
							Y:         mapY,
							Z:         2,
							Height:    64,
							Width:     64,
							TileX:     srcX,
							TileY:     srcY,
						})
					}

				}
			}
		}
	}

	objectGroups := []string{"Decoration"}

	for _, objectGroupName := range objectGroups {
		if objectGroup := tileMap.ObjectGroupWithName(objectGroupName); objectGroup != nil {
			fmt.Println("Found", objectGroupName)

			for _, obj := range objectGroup.Objects {
				if obj.GlobalID == 0 {
					continue
				}

				// Resolve the tileset from the global ID
				var tileSet *tmx.TileSet // change this to match your types
				var localID uint32
				for _, ts := range tileMap.TileSets {
					if obj.GlobalID >= ts.FirstGlobalID && int(obj.GlobalID) < int(ts.FirstGlobalID)+ts.TileCount {
						tileSet = &ts
						localID = uint32(obj.GlobalID) - uint32(ts.FirstGlobalID)
						break
					}
				}
				if tileSet == nil {
					continue // GID not found in any tileset
				}

				// Get texture path
				parts := strings.Split(tileSet.TileWithID(tmx.TileID(localID)).Image.Source, "/")
				joined := strings.Join(parts[2:], "/")

				// Calculate sprite source position (from tileset)

				// // Use object's X/Y from map
				mapX := int(obj.X)
				mapY := int(obj.Y) // Y needs correction (Tiled anchor is bottom-left)

				if texturePathMap[joined] == 0 {
					texture := rl.LoadTexture(joined)
					texturePathMap[joined] = int(texture.ID)
					textureMap[int(texture.ID)] = &texture
					textures = append(textures, &LevelData{
						TextureID: int(texture.ID),
						X:         mapX,
						Y:         mapY - int(texture.Height),
						Height:    int(texture.Height),
						Width:     int(texture.Width),
						TileX:     0,
						TileY:     0,
					})
				} else {
					textures = append(textures, &LevelData{
						TextureID: texturePathMap[joined],
						X:         mapX,
						Y:         mapY - int(textureMap[texturePathMap[joined]].Height),
						Height:    int(textureMap[texturePathMap[joined]].Height),
						Width:     int(textureMap[texturePathMap[joined]].Width),
						Z:         2,
						TileX:     0,
						TileY:     0,
					})
				}
			}
		}
	}

	return &Level{
		BackgroundTexture: levelTexture,
		MapTextures:       textures,
		TextureCache:      textureMap,
		Water:             NewWater(),
		WaterLocations:    waterLocations,
	}
}

func (l *Level) Dispose() {
	rl.UnloadTexture(*l.BackgroundTexture)
	for _, texture := range l.TextureCache {
		rl.UnloadTexture(*texture)
	}
}
