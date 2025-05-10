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
	TileX     int
	TileY     int
}

type Level struct {
	BackgroundTexture *rl.Texture2D
	MapTextures       []*LevelData
	TextureCache      map[int]*rl.Texture2D
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
	layers := []string{"Forest Grass", "HouseFloor", "HouseFurnitureBottom", "HouseWalls", "HouseFurnitureTop", "Fence"}
	textureMap := map[int]*rl.Texture2D{}
	texturePathMap := map[string]int{}

	for _, layer := range layers {
		if houseBottom := tileMap.LayerWithName(layer); houseBottom != nil {
			tiles, _ := houseBottom.TileDefs(tileMap.TileSets)
			for i, t := range tiles {
				if t.ID != 0 {
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
							TileX:     srcX,
							TileY:     srcY,
						})
					} else {
						textures = append(textures, &LevelData{
							TextureID: texturePathMap[joined],
							X:         mapX,
							Y:         mapY,
							TileX:     srcX,
							TileY:     srcY,
						})
					}

				}
			}
		}
	}

	return &Level{
		BackgroundTexture: levelTexture,
		MapTextures:       textures,
		TextureCache:      textureMap,
	}
}

func (l *Level) Dispose() {
	rl.UnloadTexture(*l.BackgroundTexture)
	for _, texture := range l.TextureCache {
		rl.UnloadTexture(*texture)
	}
}
