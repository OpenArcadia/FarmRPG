package entity

import (
	"fmt"

	"com.openarcadia.farmrpg/utils"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Water struct {
	Textures     []*rl.Texture2D
	CurrentIndex float32
}

type WaterLocation struct {
	X int
	Y int
}

func NewWater() *Water {
	textures, err := utils.ImportFolder("water")

	if err != nil {
		fmt.Println("error while loading water textures")
		return nil
	}

	return &Water{
		Textures:     textures,
		CurrentIndex: 0,
	}
}

func (w *Water) Animate() {
	w.CurrentIndex += 4 * rl.GetFrameTime()
	if int(w.CurrentIndex) >= len(w.Textures) {
		w.CurrentIndex = 0
	}
}

func (w *Water) Draw(x, y int) {
	rl.DrawTexture(*w.Textures[int(w.CurrentIndex)], int32(x), int32(y), rl.White)
}

func (w *Water) Dispose() {

}
