package utils

import (
	"os"
	"runtime"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func ImportAssetPath(path string) string {
	_, isFlatpak := os.LookupEnv("container")

	var basePath string
	if isFlatpak && runtime.GOOS == "linux" {
		basePath = "/app/bin/assets/"
	} else {
		basePath = "assets/"
	}

	return basePath + path
}

func ImportFolder(path string) ([]*rl.Texture2D, error) {
	path = ImportAssetPath(path)
	files, err := os.ReadDir(path)
	if err != nil {
		return make([]*rl.Texture2D, 0), nil
	}

	textures := make([]*rl.Texture2D, 0)

	for _, file := range files {
		if !file.IsDir() {
			texture := rl.LoadTexture(path + "/" + file.Name())
			textures = append(textures, &texture)
		}
	}

	return textures, nil
}
