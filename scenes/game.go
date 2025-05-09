package scenes

import (
	"com.openarcadia.farmrpg/entity"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	Player *entity.Player
}

func (g *Game) Create() {
	g.Player = entity.NewPlayer(200, 300)
}

func (g *Game) Render() {
	rl.ClearBackground(rl.White)
	g.Player.Draw()
	rl.DrawFPS(10, 10)
	g.Player.Update()
}

func (g *Game) Dispose() {

}
