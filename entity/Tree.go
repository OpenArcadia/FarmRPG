package entity

import rl "github.com/gen2brain/raylib-go/raylib"

type TreeType string

const (
	SmallTree TreeType = "small"
	LargeTree TreeType = "large"
)

type Tree struct {
	LevelData
	Health int
}

func NewTree(x, y, z, tileX, tileY, width, height, hitboxWidth, hitboxHeight, health, textureID int) *Tree {
	return &Tree{
		LevelData: LevelData{
			TextureID:    textureID,
			X:            x,
			Y:            y,
			Z:            z,
			TileX:        tileX,
			TileY:        tileY,
			Width:        width,
			Height:       height,
			HitBoxWidth:  hitboxWidth,
			HitBoxHeight: hitboxHeight,
		},
		Health: health,
	}
}

func (t *Tree) Draw(texture *rl.Texture2D) {
	dest := rl.NewVector2(float32(t.X), float32(t.Y))
	rl.DrawTextureRec(
		*texture,
		rl.NewRectangle(float32(t.TileX), float32(t.TileY), float32(t.Width), float32(t.Height)),
		dest,
		rl.White,
	)
}

func (t *Tree) Damage() {}
