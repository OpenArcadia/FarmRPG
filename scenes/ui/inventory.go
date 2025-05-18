package ui

import (
	"math"

	"com.openarcadia.farmrpg/utils"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Item string

func (pt *Item) ToString() string {
	return string(*pt)
}

const (
	Non   Item = ""
	Axe   Item = "axe"
	Hoe   Item = "hoe"
	Water Item = "water"
)

type InventoryItem struct {
	Name  string
	Asset *rl.Texture2D
	Tool  Item
}

type Inventory struct {
	MaxSize       int
	SelectedIndex int
	Tools         []*InventoryItem
}

var itemSize = 80

func NewInventory() *Inventory {
	axeAssetPath := utils.ImportAssetPath("overlay/axe.png")
	hoeAssetPath := utils.ImportAssetPath("overlay/hoe.png")
	waterAssetPath := utils.ImportAssetPath("overlay/water.png")

	axeAsset := rl.LoadTexture(axeAssetPath)
	hoeAsset := rl.LoadTexture(hoeAssetPath)
	waterAsset := rl.LoadTexture(waterAssetPath)

	defaultTools := []*InventoryItem{
		{
			Name:  "basic axe",
			Asset: &axeAsset,
			Tool:  Axe,
		},
		{
			Name:  "basic hoe",
			Asset: &hoeAsset,
			Tool:  Hoe,
		},
		{
			Name:  "basic water",
			Asset: &waterAsset,
			Tool:  Water,
		},
	}
	return &Inventory{
		MaxSize:       9,
		Tools:         defaultTools,
		SelectedIndex: 0,
	}
}

func (i *Inventory) AddItem(inventoryItem *InventoryItem) {
	i.Tools = append(i.Tools, inventoryItem)
}

func (in *Inventory) Update() {
	if rl.IsKeyPressed(rl.KeyOne) {
		in.SelectedIndex = 0
	} else if rl.IsKeyPressed(rl.KeyTwo) {
		in.SelectedIndex = 1
	} else if rl.IsKeyPressed(rl.KeyThree) {
		in.SelectedIndex = 2
	} else if rl.IsKeyPressed(rl.KeyFour) {
		in.SelectedIndex = 3
	} else if rl.IsKeyPressed(rl.KeyFive) {
		in.SelectedIndex = 4
	} else if rl.IsKeyPressed(rl.KeySix) {
		in.SelectedIndex = 5
	} else if rl.IsKeyPressed(rl.KeySeven) {
		in.SelectedIndex = 6
	} else if rl.IsKeyPressed(rl.KeyEight) {
		in.SelectedIndex = 7
	} else if rl.IsKeyPressed(rl.KeyNine) {
		in.SelectedIndex = 8
	}

	// Mouse scroll to cycle through inventory
	scroll := rl.GetMouseWheelMove()
	if scroll != 0 {
		in.SelectedIndex -= int(scroll)
		if in.SelectedIndex < 0 {
			in.SelectedIndex = in.MaxSize - 1
		} else if in.SelectedIndex >= in.MaxSize {
			in.SelectedIndex = 0
		}
	}
}

func (in *Inventory) Draw() {
	screenWidth := rl.GetScreenWidth()
	screenHeight := rl.GetScreenHeight()

	spacing := 0
	borderWidth := 5

	totalWidth := in.MaxSize*itemSize + (in.MaxSize-1)*spacing
	startX := (screenWidth - totalWidth) / 2
	startY := screenHeight - itemSize - 20

	for idx := 0; idx < in.MaxSize; idx++ {
		x := startX + idx*(itemSize+spacing)
		y := startY

		rl.DrawRectangle(int32(x), int32(y), int32(itemSize), int32(itemSize), rl.Color{R: 238, G: 222, B: 224, A: 200})

		for b := 0; b < borderWidth; b++ {
			rl.DrawRectangleLines(
				int32(x+b), int32(y+b),
				int32(itemSize-2*b), int32(itemSize-2*b),
				rl.Color{R: 193, G: 125, B: 99, A: 255},
			)
		}

		if idx < len(in.Tools) && in.Tools[idx] != nil && in.Tools[idx].Asset != nil {
			tex := in.Tools[idx].Asset

			texWidth := tex.Width
			texHeight := tex.Height

			scale := float32(itemSize) * 0.6 / float32(math.Max(float64(texWidth), float64(texHeight))) // Fit texture in 60% of slot
			drawWidth := float32(texWidth) * scale
			drawHeight := float32(texHeight) * scale
			posX := float32(x) + (float32(itemSize)-drawWidth)/2
			posY := float32(y) + (float32(itemSize)-drawHeight)/2

			rl.DrawTextureEx(*tex, rl.Vector2{X: posX, Y: posY}, 0, scale, rl.White)
		}
	}
	x := startX + in.SelectedIndex*(itemSize+spacing)

	rl.DrawRectangleLinesEx(rl.Rectangle{
		X:      float32(x - 2),
		Y:      float32(startY - 2),
		Width:  float32(itemSize + 4),
		Height: float32(itemSize + 4),
	}, 3, rl.Gold)
}

func (in *Inventory) Dispose() {
	for _, tool := range in.Tools {
		rl.UnloadTexture(*tool.Asset)
	}
}
