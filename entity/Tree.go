package entity

import (
	"math/rand"

	"com.openarcadia.farmrpg/utils"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type TreeType string

const (
	SmallTree TreeType = "small"
	LargeTree TreeType = "large"
)

type Tree struct {
	LevelData
	Health         int
	StumpTextureID int
	Type           TreeType
	shakeDuration  float32
	shakeTimer     float32
	shakeOffsetX   float32
	isFalling      bool
	fallTimer      float32
	fallRotation   float32
	fadeAlpha      float32
	visible        bool
	Apples         []*Apple
}

type Apple struct {
	Offset rl.Vector2
}

var AppleTexture *rl.Texture2D
var applePositionsByType = map[TreeType][]rl.Vector2{
	SmallTree: {
		{X: 18, Y: 17},
		{X: 30, Y: 37},
		{X: 12, Y: 50},
		{X: 30, Y: 45},
		{X: 20, Y: 30},
		{X: 30, Y: 10},
	},
	LargeTree: {
		{X: 30, Y: 24},
		{X: 60, Y: 65},
		{X: 50, Y: 50},
		{X: 16, Y: 40},
		{X: 45, Y: 50},
		{X: 42, Y: 70},
	},
}

func LoadAssets() {
	if AppleTexture == nil {
		texture := rl.LoadTexture(utils.ImportAssetPath("fruit/apple.png"))
		AppleTexture = &texture
	}
}

func NewTree(x, y, z, tileX, tileY, width, height, hitboxWidth, hitboxHeight, health, textureID int, treeType TreeType) *Tree {
	LoadAssets()
	tree := &Tree{
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
		Type:      treeType,
		Health:    health,
		visible:   true,
		fadeAlpha: 1.0,
	}

	if rand.Float32() < 0.2 {
		applePositions := applePositionsByType[treeType]

		rand.Shuffle(len(applePositions), func(i, j int) {
			applePositions[i], applePositions[j] = applePositions[j], applePositions[i]
		})

		maxApples := 3
		count := 0

		for i := 0; i < len(applePositions) && count < maxApples; i++ {
			if rand.Float32() <= 0.2 {
				pos := applePositions[i]
				apple := &Apple{
					Offset: rl.Vector2{
						X: pos.X,
						Y: pos.Y,
					},
				}
				tree.Apples = append(tree.Apples, apple)
				count++
			}
		}

	}
	return tree
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

	for _, apple := range t.Apples {
		ax := float32(t.X) + apple.Offset.X
		ay := float32(t.Y) + apple.Offset.Y
		rl.DrawTexture(*AppleTexture, int32(ax), int32(ay), rl.White)
	}
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

	if len(t.Apples) > 0 {
		index := rand.Intn(len(t.Apples))
		t.Apples = append(t.Apples[:index], t.Apples[index+1:]...)
	}

	if t.Health <= 0 && !t.isFalling {
		t.isFalling = true
		t.fallTimer = 1.0
		t.fadeAlpha = 1.0
	}
}
