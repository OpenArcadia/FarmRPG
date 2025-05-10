package main

import (
	"runtime"

	"com.openarcadia.farmrpg/scenes"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func init() {
	rl.SetCallbackFunc(main)
}

func main() {
	rl.InitWindow(0, 0, "FarmRPG")

	rl.InitAudioDevice()

	rl.SetExitKey(0)

	rl.ToggleFullscreen()

	scenes.ChangeScreen(&scenes.Game{})

	windowShouldClose := false

	rl.TraceLog(rl.LogAll, "Entering Main Game Loog")
	for !windowShouldClose {
		rl.BeginDrawing()

		if runtime.GOOS == "android" && rl.IsKeyDown(rl.KeyBack) || rl.WindowShouldClose() {
			windowShouldClose = true
		}
		rl.TraceLog(rl.LogAll, "Updating Scene")
		scenes.Update()

		rl.EndDrawing()
	}

	scenes.Dispose()

	rl.CloseWindow()
	rl.CloseAudioDevice()
}
