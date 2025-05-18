package entity

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TreeType string

const (
	SmallTree TreeType = "small"
	LargeTree TreeType = "large"
)

type Tree struct {
	LevelData
	Health        int
	shakeDuration float32 // how long to shake (in seconds)
	shakeTimer    float32 // current timer
	shakeOffsetX  float32 // current X offset for shaking
	isFalling     bool
	fallTimer     float32
	fallRotation  float32
	fadeAlpha     float32
	visible       bool
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
		Health:    health,
		visible:   true,
		fadeAlpha: 1.0,
	}
}

func (t *Tree) Draw(texture *rl.Texture2D) {
	if t.Health <= 0 && !t.visible {
		return
	}

	// Position and rotation origin
	dest := rl.NewVector2(float32(t.X)+t.shakeOffsetX+float32(t.Width)/2, float32(t.Y)+float32(t.Height))
	src := rl.NewRectangle(float32(t.TileX), float32(t.TileY), float32(t.Width), float32(t.Height))
	origin := rl.NewVector2(float32(t.Width)/2, float32(t.Height))

	tint := rl.Fade(rl.White, t.fadeAlpha)
	rl.DrawTexturePro(*texture, src, rl.NewRectangle(dest.X, dest.Y, float32(t.Width), float32(t.Height)), origin, t.fallRotation, tint)
}

func (t *Tree) Update() {
	dt := rl.GetFrameTime()

	if t.shakeTimer > 0 {
		t.shakeTimer -= dt
		t.shakeOffsetX = float32(rand.Intn(5) - 2)
		if t.shakeTimer <= 0 {
			t.shakeOffsetX = 0
		}
	}

	if t.isFalling && t.visible {
		t.fallTimer -= dt
		t.fallRotation += dt * 90 // Rotate 90 degrees over fall duration
		t.fadeAlpha -= dt * 1.2   // Fade out quickly

		if t.fadeAlpha < 0 {
			t.fadeAlpha = 0
			t.visible = false
			t.isFalling = false
		}
	}
}

func (t *Tree) Damage() {
	t.Health -= 1
	t.shakeDuration = 0.2
	t.shakeTimer = t.shakeDuration

	if t.Health <= 0 && !t.isFalling {
		t.isFalling = true
		t.fallTimer = 1.0
		t.fadeAlpha = 1.0
	}
}
