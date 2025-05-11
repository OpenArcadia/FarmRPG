package scenes

import (
	"sort"
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

type Drawable struct {
	Z        int
	Y        float32
	Height   int
	DrawFunc func()
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

	var drawables []Drawable

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, tile := range g.Level.MapTextures {
			tileZ := tile.Z
			tileY := float32(tile.Y)
			tileCopy := tile

			drawables = append(drawables, Drawable{
				Z:      tileZ,
				Height: tileCopy.Height,
				Y:      tileY,
				DrawFunc: func() {
					dest := rl.NewVector2(float32(tileCopy.X), float32(tileCopy.Y))
					rl.DrawTextureRec(
						*g.Level.TextureCache[tileCopy.TextureID],
						rl.NewRectangle(float32(tileCopy.TileX), float32(tileCopy.TileY), float32(tileCopy.Width), float32(tileCopy.Height)),
						dest,
						rl.White,
					)
				},
			})
		}

		// Add player as Z = 2 drawable
		playerY := g.Player.GetRect().Y
		drawables = append(drawables, Drawable{
			Z:      2,
			Y:      playerY,
			Height: 64,
			DrawFunc: func() {
				g.Player.Draw()
			},
		})

		// First: sort by Z ascending, then by Y (only when Z > 0)
		sort.Slice(drawables, func(i, j int) bool {
			if drawables[i].Z != drawables[j].Z {
				return drawables[i].Z < drawables[j].Z
			}
			return drawables[i].Y+float32(drawables[i].Height) < drawables[j].Y+float32(drawables[j].Height)
		})
	}()

	wg.Wait()
	// Draw in sorted order
	for _, d := range drawables {
		d.DrawFunc()
	}

	rl.EndMode2D()

	g.Inventory.Draw()

	rl.DrawFPS(10, 10)
}

func (g *Game) Dispose() {
	g.Player.Dispose()
	g.Level.Dispose()
	g.Inventory.Dispose()
}
