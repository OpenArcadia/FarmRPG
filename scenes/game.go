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
	g.Player = entity.NewPlayer(1050, 1050, g.Inventory, g.Level.MapTextures, g.Level.Trees)
	g.Camera = &rl.Camera2D{
		Target:   rl.NewVector2(g.Player.GetRect().X, g.Player.GetRect().Y),
		Offset:   rl.NewVector2(float32(rl.GetScreenWidth()/2), float32(rl.GetScreenHeight()/2)),
		Rotation: 0.0,
		Zoom:     1.0,
	}
}

func (g *Game) Render() {
	// Update camera target
	playerPos := g.Player.GetRect()

	// Assuming level size is based on background texture
	levelWidth := float32(g.Level.BackgroundTexture.Width)
	levelHeight := float32(g.Level.BackgroundTexture.Height)

	// Get screen dimensions
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())

	halfScreenWidth := screenWidth / 2
	halfScreenHeight := screenHeight / 2

	// Clamp camera target to prevent it from moving beyond level edges
	cameraX := rl.Clamp(playerPos.X, halfScreenWidth, levelWidth-halfScreenWidth)
	cameraY := rl.Clamp(playerPos.Y, halfScreenHeight, levelHeight-halfScreenHeight)

	g.Camera.Target = rl.NewVector2(cameraX, cameraY)

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

	wg.Wait()

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
			if tile.TextureID != 0 {
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
		}

		for _, tree := range g.Level.Trees {
			if tree.TextureID != 0 {
				drawables = append(drawables, Drawable{
					Z:      tree.Z,
					Height: tree.Height,
					Y:      float32(tree.Y),
					DrawFunc: func() {
						tree.Update()
						tree.Draw(g.Level.TextureCache[tree.TextureID])
					},
				})
			}
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
		sort.SliceStable(drawables, func(i, j int) bool {
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

	// for _, ld := range g.Level.MapTextures {
	// 	rl.DrawRectangle(int32(ld.GetHitBoxRect().X), int32(ld.GetHitBoxRect().Y), int32(ld.HitBoxWidth), int32(ld.HitBoxHeight), rl.Lime)
	// }

	rl.EndMode2D()

	g.Inventory.Draw()

	rl.DrawFPS(10, 10)
}

func (g *Game) Dispose() {
	g.Player.Dispose()
	g.Level.Dispose()
	g.Inventory.Dispose()
}
