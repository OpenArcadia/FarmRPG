package main

import (
	"com.openarcadia.farmrpg/scenes"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(0, 0, "FarmRPG")
	defer rl.CloseWindow()
	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	rl.SetExitKey(0)

	rl.ToggleFullscreen()

	scenes.ChangeScreen(&scenes.Game{})

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		scenes.Update()

		rl.EndDrawing()
	}
}
