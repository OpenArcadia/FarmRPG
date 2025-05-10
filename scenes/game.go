package scenes

import (
	"com.openarcadia.farmrpg/entity"
	"com.openarcadia.farmrpg/scenes/ui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	Player    *entity.Player
	Camera    *rl.Camera2D
	Inventory *ui.Inventory
	Level     *entity.Level
}

func (g *Game) Create() {
	g.Level = entity.NewLevel()
	g.Inventory = ui.NewInventory()
	g.Player = entity.NewPlayer(200, 300, g.Inventory)
	g.Camera = &rl.Camera2D{
		Target:   rl.NewVector2(g.Player.GetRect().X, g.Player.GetRect().Y),
		Offset:   rl.NewVector2(float32(rl.GetScreenWidth()/2), float32(rl.GetScreenHeight()/2)),
		Rotation: 0.0,
		Zoom:     1.0,
	}

}

func (g *Game) Render() {
	rl.ClearBackground(rl.White)

	g.Camera.Target = rl.NewVector2(g.Player.GetRect().X, g.Player.GetRect().Y)

	rl.BeginMode2D(*g.Camera)

	// Draw the level background
	if g.Level.BackgroundTexture != nil {
		rl.DrawTexture(*g.Level.BackgroundTexture, 0, 0, rl.White)
	}

	for _, tile := range g.Level.MapTextures {
		dest := rl.NewVector2(float32(tile.X), float32(tile.Y))
		rl.DrawTextureRec(*g.Level.TextureCache[tile.TextureID], rl.NewRectangle(float32(tile.TileX), float32(tile.TileY), float32(64), float32(64)), dest, rl.White)
	}

	// Draw player and UI
	g.Player.Draw()
	// camera will not affect other element drawing
	rl.EndMode2D()

	g.Inventory.Draw()

	// // Draw FPS and update logic
	rl.DrawFPS(10, 10)
	g.Player.Update()
	g.Inventory.Update()
}

func (g *Game) Dispose() {
	g.Player.Dispose()
	g.Level.Dispose()
	g.Inventory.Dispose()
}
