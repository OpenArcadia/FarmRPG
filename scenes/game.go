package scenes

import (
	"sync"

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
	// Update camera target
	g.Camera.Target = rl.NewVector2(g.Player.GetRect().X, g.Player.GetRect().Y)

	// Step 1: Run updates in parallel
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		g.Level.Water.Animate()
	}()

	go func() {
		defer wg.Done()
		g.Player.Update()
	}()

	go func() {
		defer wg.Done()
		g.Inventory.Update()
	}()

	wg.Wait() // Ensure updates are complete before rendering

	// Step 2: Render everything (on main thread only)
	rl.ClearBackground(rl.White)
	rl.BeginMode2D(*g.Camera)

	for _, loc := range g.Level.WaterLocations {
		g.Level.Water.Draw(loc.X, loc.Y)
	}

	if g.Level.BackgroundTexture != nil {
		rl.DrawTexture(*g.Level.BackgroundTexture, 0, 0, rl.White)
	}

	for _, tile := range g.Level.MapTextures {
		dest := rl.NewVector2(float32(tile.X), float32(tile.Y))
		rl.DrawTextureRec(
			*g.Level.TextureCache[tile.TextureID],
			rl.NewRectangle(float32(tile.TileX), float32(tile.TileY), 64, 64),
			dest,
			rl.White,
		)
	}

	g.Player.Draw()
	rl.EndMode2D()

	g.Inventory.Draw()

	rl.DrawFPS(10, 10)
}

func (g *Game) Dispose() {
	g.Player.Dispose()
	g.Level.Dispose()
	g.Inventory.Dispose()
}
