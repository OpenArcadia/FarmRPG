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

var itemSize = 60

func NewInventory() *Inventory {
	axeAssetPath := utils.ImportAssetPath("overlay/axe.png")
	hoeAssetPath := utils.ImportAssetPath("overlay/hoe.png")
	waterAssetPath := utils.ImportAssetPath("overlay/water.png")

	axeAsset := rl.LoadTexture(axeAssetPath)
	hoeAsset := rl.LoadTexture(hoeAssetPath)
	waterAsset := rl.LoadTexture(waterAssetPath)

	defaultTools := []*InventoryItem{
		&InventoryItem{
			Name:  "basic axe",
			Asset: &axeAsset,
			Tool:  Axe,
		},
		&InventoryItem{
			Name:  "basic hoe",
			Asset: &hoeAsset,
			Tool:  Hoe,
		},
		&InventoryItem{
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
	if rl.IsKeyDown(rl.KeyOne) {
		in.SelectedIndex = 0
	}
	if rl.IsKeyDown(rl.KeyTwo) {
		in.SelectedIndex = 1
	}
	if rl.IsKeyDown(rl.KeyThree) {
		in.SelectedIndex = 2
	}
}

func (in *Inventory) Draw() {
	screenWidth := rl.GetScreenWidth()
	screenHeight := rl.GetScreenHeight()

	spacing := 0
	borderWidth := 3

	totalWidth := in.MaxSize*itemSize + (in.MaxSize-1)*spacing
	startX := (screenWidth - totalWidth) / 2
	startY := screenHeight - itemSize - 20

	for idx := 0; idx < in.MaxSize; idx++ {
		x := startX + idx*(itemSize+spacing)
		y := startY

		// Slot background (blur-like solid color)
		rl.DrawRectangle(int32(x), int32(y), int32(itemSize), int32(itemSize), rl.Color{R: 238, G: 222, B: 224, A: 200})

		// Thick border
		for b := 0; b < borderWidth; b++ {
			rl.DrawRectangleLines(
				int32(x+b), int32(y+b),
				int32(itemSize-2*b), int32(itemSize-2*b),
				rl.Color{R: 193, G: 125, B: 99, A: 255},
			)
		}

		// Draw item texture if present in the tools list
		if idx < len(in.Tools) && in.Tools[idx] != nil && in.Tools[idx].Asset != nil {
			tex := in.Tools[idx].Asset
			// Center the texture inside the slot
			texWidth := tex.Width
			texHeight := tex.Height

			scale := float32(itemSize) * 0.6 / float32(math.Max(float64(texWidth), float64(texHeight))) // Fit texture in 60% of slot
			drawWidth := float32(texWidth) * scale
			drawHeight := float32(texHeight) * scale
			posX := float32(x) + (float32(itemSize)-drawWidth)/2
			posY := float32(y) + (float32(itemSize)-drawHeight)/2

			rl.DrawTextureEx(*tex, rl.Vector2{X: posX, Y: posY}, 0, scale, rl.White)
		}

		// Highlight selected slot
		if idx == in.SelectedIndex {

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
