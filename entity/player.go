package entity

import (
	"fmt"
	"math"

	"com.openarcadia.farmrpg/scenes/ui"
	"com.openarcadia.farmrpg/utils"
	"com.openarcadia.farmrpg/utils/timer"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type PlayerStatus string

func (ps *PlayerStatus) ToString() string {
	return string(*ps)
}

const (
	Down       PlayerStatus = "down"
	DownAxe    PlayerStatus = "down_axe"
	DownHoe    PlayerStatus = "down_hoe"
	DownIdle   PlayerStatus = "down_idle"
	DownWater  PlayerStatus = "down_water"
	Up         PlayerStatus = "up"
	UpAxe      PlayerStatus = "up_axe"
	UpHoe      PlayerStatus = "up_hoe"
	UpIdle     PlayerStatus = "up_idle"
	UpWater    PlayerStatus = "up_water"
	Left       PlayerStatus = "left"
	LeftAxe    PlayerStatus = "left_axe"
	LeftHoe    PlayerStatus = "left_hoe"
	LeftIdle   PlayerStatus = "left_idle"
	LeftWater  PlayerStatus = "left_water"
	Right      PlayerStatus = "right"
	RightAxe   PlayerStatus = "right_axe"
	RightHoe   PlayerStatus = "right_hoe"
	RightIdle  PlayerStatus = "right_idle"
	RightWater PlayerStatus = "right_water"
)

type Player struct {
	x          float32
	y          float32
	height     int
	width      int
	status     PlayerStatus
	frameIndex float32
	timer      *map[string]*timer.Timer
	animations *map[string][]*rl.Texture2D
	inventory  *ui.Inventory
}

func NewPlayer(x, y float32, inventory *ui.Inventory) *Player {

	animationsDict := &map[string][]*rl.Texture2D{
		"down":        make([]*rl.Texture2D, 0),
		"down_axe":    make([]*rl.Texture2D, 0),
		"down_hoe":    make([]*rl.Texture2D, 0),
		"down_idle":   make([]*rl.Texture2D, 0),
		"down_water":  make([]*rl.Texture2D, 0),
		"up":          make([]*rl.Texture2D, 0),
		"up_axe":      make([]*rl.Texture2D, 0),
		"up_hoe":      make([]*rl.Texture2D, 0),
		"up_idle":     make([]*rl.Texture2D, 0),
		"up_water":    make([]*rl.Texture2D, 0),
		"left":        make([]*rl.Texture2D, 0),
		"left_axe":    make([]*rl.Texture2D, 0),
		"left_hoe":    make([]*rl.Texture2D, 0),
		"left_idle":   make([]*rl.Texture2D, 0),
		"left_water":  make([]*rl.Texture2D, 0),
		"right":       make([]*rl.Texture2D, 0),
		"right_axe":   make([]*rl.Texture2D, 0),
		"right_hoe":   make([]*rl.Texture2D, 0),
		"right_idle":  make([]*rl.Texture2D, 0),
		"right_water": make([]*rl.Texture2D, 0),
	}

	for key, _ := range *animationsDict {
		textures, err := utils.ImportFolder("character/" + key)
		if err != nil {
			fmt.Printf("error occured while loading texture %s", key)
			continue
		}
		(*animationsDict)[key] = textures
	}

	p := &Player{
		x:          x,
		y:          y,
		height:     72,
		width:      60,
		animations: animationsDict,
		status:     DownIdle,
		frameIndex: 0,
		inventory:  inventory,
	}

	timerMap := map[string]*timer.Timer{
		"use tool":    timer.NewTimer(750, p.useTool),
		"tool switch": timer.NewTimer(200, nil),
	}

	p.timer = &timerMap

	return p
}

func (p *Player) Draw() {
	rl.DrawTexture(*(*p.animations)[p.status.ToString()][int(p.frameIndex)], int32(p.x-55), int32(p.y-32), rl.White)
}

func (p *Player) Update() {
	if !(*p.timer)["use tool"].IsActive() {
		p.handleInput()
		// if player is not moving then setting the state based on movement direction
		if !rl.IsKeyDown(rl.KeyW) &&
			!rl.IsKeyDown(rl.KeyA) &&
			!rl.IsKeyDown(rl.KeyS) &&
			!rl.IsKeyDown(rl.KeyD) {

			switch p.status {
			case Left, LeftAxe, LeftHoe, LeftWater:
				p.status = LeftIdle
			case Up, UpAxe, UpHoe, UpWater:
				p.status = UpIdle
			case Right, RightAxe, RightHoe, RightWater:
				p.status = RightIdle
			case Down, DownAxe, DownHoe, DownWater:
				p.status = DownIdle
			}
		}
	} else {
		(*p.timer)["use tool"].Update()
	}

	p.animate()

	if (*p.timer)["use tool"].IsActive() {
		item := p.inventory.Tools[p.inventory.SelectedIndex]
		if item == nil {
			(*p.timer)["use tool"].Deactivate()
			return
		}

		switch p.status {
		case LeftIdle, Left:
			if item.Tool == ui.Axe {
				p.status = LeftAxe
			} else if item.Tool == ui.Water {
				p.status = LeftWater
			} else if item.Tool == ui.Hoe {
				p.status = LeftHoe
			}
		case Right, RightIdle:
			if item.Tool == ui.Axe {
				p.status = RightAxe
			} else if item.Tool == ui.Water {
				p.status = RightWater
			} else if item.Tool == ui.Hoe {
				p.status = RightHoe
			}
		case Up, UpIdle:
			if item.Tool == ui.Axe {
				p.status = UpAxe
			} else if item.Tool == ui.Water {
				p.status = UpWater
			} else if item.Tool == ui.Hoe {
				p.status = UpHoe
			}
		case Down, DownIdle:
			if item.Tool == ui.Axe {
				p.status = DownAxe
			} else if item.Tool == ui.Water {
				p.status = DownWater
			} else if item.Tool == ui.Hoe {
				p.status = DownHoe
			}
		}

	}
}

func (p *Player) handleInput() {
	dt := rl.GetFrameTime()

	var moveX, moveY float32

	if rl.IsKeyDown(rl.KeyW) {
		moveY -= 1
		p.status = Up
	}
	if rl.IsKeyDown(rl.KeyS) {
		moveY += 1
		p.status = Down
	}
	if rl.IsKeyDown(rl.KeyD) {
		moveX += 1
		p.status = Right
	}
	if rl.IsKeyDown(rl.KeyA) {
		moveX -= 1
		p.status = Left
	}

	// Normalize the direction vector
	length := float32(math.Sqrt(float64(moveX*moveX + moveY*moveY)))
	if length != 0 {
		moveX /= length
		moveY /= length
	}

	speed := dt * 250
	p.x += moveX * speed
	p.y += moveY * speed

	// tools
	if rl.IsMouseButtonDown(rl.MouseButtonLeft) && !(*p.timer)["use tool"].IsActive() {
		if p.inventory.SelectedIndex < len(p.inventory.Tools) {
			(*p.timer)["use tool"].Activate()
			p.frameIndex = 0
		}
	}
}

func (p *Player) animate() {
	p.frameIndex += 4 * rl.GetFrameTime()
	if int(p.frameIndex) >= len((*p.animations)[p.status.ToString()]) {
		p.frameIndex = 0
	}
}

func (p *Player) useTool() {
	// fmt.Println(p.tool)
}

func (p *Player) Dispose() {
	for _, textures := range *p.animations {
		for _, texture := range textures {
			rl.UnloadTexture(*texture)
		}
	}
}

func (p *Player) GetRect() *rl.Rectangle {
	return &rl.Rectangle{
		X:      float32(p.x),
		Y:      float32(p.y),
		Width:  float32(p.width),
		Height: float32(p.height),
	}
}
