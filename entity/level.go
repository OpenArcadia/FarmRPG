package entity

import (
	"fmt"
	"strings"

	"com.openarcadia.farmrpg/utils"
	"github.com/akashKarmakar02/tmx"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type LevelData struct {
	TextureID    int
	X            int
	Y            int
	Z            int
	Height       int
	Width        int
	HitBoxHeight int
	HitBoxWidth  int
	TileX        int
	TileY        int
}

type Level struct {
	BackgroundTexture *rl.Texture2D
	MapTextures       []*LevelData
	Water             *Water
	WaterLocations    []*WaterLocation
	Trees             []*Tree
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

					textureID := utils.GetTextureID(joined)

					textures = append(textures, &LevelData{
						TextureID: textureID,
						X:         mapX,
						Y:         mapY,
						Z:         0,
						TileX:     srcX,
						TileY:     srcY,
						Height:    64,
						Width:     64,
					})

				}
			}
		}
	}

	layers = []string{"HouseWalls", "Fence", "HouseFurnitureTop"}
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

					textureID := utils.GetTextureID(joined)

					textures = append(textures, &LevelData{
						TextureID: textureID,
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

	layers = []string{"Collision"}
	for _, layer := range layers {
		if houseBottom := tileMap.LayerWithName(layer); houseBottom != nil {
			tiles, _ := houseBottom.TileDefs(tileMap.TileSets)

			for i, t := range tiles {
				if t.GlobalID != 0 {
					tileIndexInSet := t.ID - t.TileSet.FirstGlobalID.TileID(t.TileSet)
					columns := t.TileSet.Columns
					srcX := (int(tileIndexInSet) % columns) * t.TileSet.TileWidth
					srcY := (int(tileIndexInSet) / columns) * t.TileSet.TileHeight

					mapX := (i % tileMap.Width) * tileMap.TileWidth
					mapY := (i / tileMap.Width) * tileMap.TileHeight

					textures = append(textures, &LevelData{
						X:            mapX,
						Y:            mapY,
						Z:            2,
						Height:       64,
						Width:        64,
						HitBoxHeight: 64,
						HitBoxWidth:  64,
						TileX:        srcX,
						TileY:        srcY,
					})

				}
			}
		}
	}

	objectGroups := []string{"Decoration", "Trees"}
	trees := make([]*Tree, 0)

	for _, objectGroupName := range objectGroups {
		if objectGroup := tileMap.ObjectGroupWithName(objectGroupName); objectGroup != nil {
			if objectGroupName == "Collision" {
				fmt.Println("Found collision thinks")
			}
			for _, obj := range objectGroup.Objects {
				if obj.GlobalID == 0 {
					continue
				}

				// // Use object's X/Y from map
				mapX := int(obj.X)
				mapY := int(obj.Y) // Y needs correction (Tiled anchor is bottom-left)

				var hitboxHeight int
				var hitboxWidth int

				var texture *rl.Texture2D
				var textureID int
				// var stumpTextureID int

				// find stump path for tree
				// if obj.Name == "Large" {
				// 	stumpTextureID = utils.GetTextureID(utils.ImportAssetPath("objects/stump_medium.png"))
				// } else if obj.Name == "Small" {
				// 	stumpTextureID = utils.GetTextureID(utils.ImportAssetPath("objects/stump_small.png"))
				// }

				textureID = utils.GetTextureID(obj.Image.Source)

				texture = utils.GetTextureFromID(textureID)

				// tree need different collision hitbox
				if objectGroupName == "Trees" {
					hitboxHeight = int(float32(texture.Height) * 0.30)
					hitboxWidth = int(texture.Width)
				} else {
					hitboxHeight = int(float32(texture.Height) * 0.1)
					hitboxWidth = int(texture.Width)
				}

				if objectGroupName == "Trees" {
					treeType := LargeTree
					if obj.Name == "Small" {
						treeType = SmallTree
					}
					trees = append(trees, NewTree(mapX, mapY-int(texture.Height), 2, 0, 0, int(texture.Width), int(texture.Height), hitboxWidth, hitboxHeight, 5, textureID, treeType))
				} else {
					textures = append(textures, &LevelData{
						TextureID:    textureID,
						X:            mapX,
						Y:            mapY - int(texture.Height),
						Height:       int(texture.Height),
						Width:        int(texture.Width),
						HitBoxHeight: hitboxHeight,
						HitBoxWidth:  hitboxWidth,
						Z:            2,
						TileX:        0,
						TileY:        0,
					})
				}

			}
		}
	}

	return &Level{
		BackgroundTexture: levelTexture,
		MapTextures:       textures,
		Water:             NewWater(),
		WaterLocations:    waterLocations,
		Trees:             trees,
	}
}

func (ld *LevelData) GetHitBoxRect() *rl.Rectangle {
	return &rl.Rectangle{
		X:      float32(ld.X),
		Y:      float32(ld.Y + (ld.Height - ld.HitBoxHeight)),
		Width:  float32(ld.HitBoxWidth),
		Height: float32(ld.HitBoxHeight),
	}
}

func (l *Level) Dispose() {
	rl.UnloadTexture(*l.BackgroundTexture)
}
