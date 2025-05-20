package utils

import rl "github.com/gen2brain/raylib-go/raylib"

var textureMap = map[int]*rl.Texture2D{}
var texturePathMap = map[string]int{}

func GetTextureID(path string) int {
	if texturePathMap[path] == 0 {
		texture := rl.LoadTexture(path)
		textureMap[int(texture.ID)] = &texture
		texturePathMap[path] = int(texture.ID)
		return int(texture.ID)
	}
	return texturePathMap[path]
}

func GetTextureFromID(id int) *rl.Texture2D {
	return textureMap[id]
}
