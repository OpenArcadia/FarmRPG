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
	x                float32
	y                float32
	height           int
	width            int
	status           PlayerStatus
	frameIndex       float32
	timer            *map[string]*timer.Timer
	animations       *map[string][]*rl.Texture2D
	inventory        *ui.Inventory
	colidableObjects []*LevelData
	trees            []*Tree
	HitBox           *rl.Rectangle
}

type CollisionInfo struct {
	Collided bool
	Top      bool
	Bottom   bool
	Left     bool
	Right    bool
	Object   *LevelData
	Tree     *Tree
}

func NewPlayer(x, y float32, inventory *ui.Inventory, levelData []*LevelData, trees []*Tree) *Player {

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
		x:                x,
		y:                y,
		height:           72,
		width:            60,
		animations:       animationsDict,
		status:           DownIdle,
		frameIndex:       0,
		inventory:        inventory,
		colidableObjects: levelData,
		trees:            trees,
	}

	timerMap := map[string]*timer.Timer{
		"use tool":    timer.NewTimer(500, p.useTool),
		"tool switch": timer.NewTimer(200, nil),
	}

	hitRect := rl.Rectangle{
		X:      x,
		Y:      y,
		Width:  10,
		Height: 10,
	}

	p.timer = &timerMap
	p.HitBox = &hitRect

	return p
}

func (p *Player) Draw() {
	rl.DrawEllipse(int32(p.x+32), int32(p.y+64), 20, 8, rl.Fade(rl.Black, 0.3))
	rl.DrawTexture(*(*p.animations)[p.status.ToString()][int(p.frameIndex)], int32(p.x-55), int32(p.y-32), rl.White)
}

func (p *Player) UpdateHitBox() {
	const offset = 10.0

	switch p.status {
	case Up, UpAxe, UpHoe, UpIdle, UpWater:
		p.HitBox.X = p.x + float32(p.width)/2 // center horizontally
		p.HitBox.Y = p.y - offset - 5         // above the player
	case Down, DownAxe, DownHoe, DownIdle, DownWater:
		p.HitBox.X = p.x + float32(p.width)/2
		p.HitBox.Y = p.y + float32(p.height) + 20
	case Left, LeftAxe, LeftHoe, LeftIdle, LeftWater:
		p.HitBox.X = p.x - offset - 20
		p.HitBox.Y = p.y + float32(p.height)/2 + 20
	case Right, RightAxe, RightHoe, RightIdle, RightWater:
		p.HitBox.X = p.x + float32(p.width) + 20
		p.HitBox.Y = p.y + float32(p.height)/2 + 20
	}
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

	p.UpdateHitBox()
}

func (p *Player) isColliding() []CollisionInfo {
	playerRect := *p.GetHitBoxRect()
	var collisions []CollisionInfo

	for _, obj := range p.colidableObjects {
		if obj.Z != 2 {
			continue
		}

		objRect := *obj.GetHitBoxRect()
		if objRect.Height > 0 && objRect.Width > 0 && rl.CheckCollisionRecs(playerRect, objRect) {
			info := CollisionInfo{
				Collided: true,
				Object:   obj,
			}

			// Compute the edges
			playerRight := playerRect.X + playerRect.Width
			playerBottom := playerRect.Y + playerRect.Height
			objRight := objRect.X + objRect.Width
			objBottom := objRect.Y + objRect.Height

			// Determine which sides are overlapping
			if playerBottom > objRect.Y && playerRect.Y < objRect.Y {
				info.Bottom = true // Player is hitting top of the object
			}
			if playerRect.Y < objBottom && playerBottom > objBottom {
				info.Top = true // Player is hitting bottom of the object
			}
			if playerRight > objRect.X && playerRect.X < objRect.X {
				info.Right = true // Player is hitting left side of the object
			}
			if playerRect.X < objRight && playerRight > objRight {
				info.Left = true // Player is hitting right side of the object
			}

			collisions = append(collisions, info)
		}
	}

	for _, obj := range p.trees {
		if obj.Z != 2 {
			continue
		}
		if obj.Health < 1 {
			continue
		}

		objRect := *obj.GetHitBoxRect()
		if objRect.Height > 0 && objRect.Width > 0 && rl.CheckCollisionRecs(playerRect, objRect) {
			info := CollisionInfo{
				Collided: true,
				Tree:     obj,
			}

			// Compute the edges
			playerRight := playerRect.X + playerRect.Width
			playerBottom := playerRect.Y + playerRect.Height
			objRight := objRect.X + objRect.Width
			objBottom := objRect.Y + objRect.Height

			// Determine which sides are overlapping
			if playerBottom > objRect.Y && playerRect.Y < objRect.Y {
				info.Bottom = true // Player is hitting top of the object
			}
			if playerRect.Y < objBottom && playerBottom > objBottom {
				info.Top = true // Player is hitting bottom of the object
			}
			if playerRight > objRect.X && playerRect.X < objRect.X {
				info.Right = true // Player is hitting left side of the object
			}
			if playerRect.X < objRight && playerRight > objRight {
				info.Left = true // Player is hitting right side of the object
			}

			collisions = append(collisions, info)
		}
	}

	return collisions
}

func getOverlap(player, obj rl.Rectangle) (float32, float32) {
	pxRight := player.X + player.Width
	pxBottom := player.Y + player.Height
	oxRight := obj.X + obj.Width
	oyBottom := obj.Y + obj.Height

	overlapX := float32(0)
	overlapY := float32(0)

	if player.X < oxRight && pxRight > obj.X {
		if pxRight-obj.X < oxRight-player.X {
			overlapX = pxRight - obj.X
		} else {
			overlapX = -(oxRight - player.X)
		}
	}

	if player.Y < oyBottom && pxBottom > obj.Y {
		if pxBottom-obj.Y < oyBottom-player.Y {
			overlapY = pxBottom - obj.Y
		} else {
			overlapY = -(oyBottom - player.Y)
		}
	}

	return overlapX, overlapY
}

func (p *Player) handleInput() {
	dt := rl.GetFrameTime()
	collisions := p.isColliding()

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

	// Normalize movement
	length := float32(math.Sqrt(float64(moveX*moveX + moveY*moveY)))
	if length != 0 {
		moveX /= length
		moveY /= length
	}

	// Adjust movement based on collisions
	for _, c := range collisions {
		if !c.Collided {
			continue
		}

		if c.Top && moveY < 0 {
			moveY = 0
		}
		if c.Bottom && moveY > 0 {
			moveY = 0
		}
		if c.Left && moveX < 0 {
			moveX = 0
		}
		if c.Right && moveX > 0 {
			moveX = 0
		}
	}

	// Apply movement
	speed := dt * 250
	p.x += moveX * speed
	p.y += moveY * speed

	playerRect := *p.GetHitBoxRect()
	for _, c := range p.isColliding() {
		if !c.Collided {
			continue
		}
		var objRect rl.Rectangle
		if c.Object != nil {
			objRect = *c.Object.GetHitBoxRect()
		} else if c.Tree != nil {
			objRect = *c.Tree.GetHitBoxRect()
		}

		if rl.CheckCollisionRecs(playerRect, objRect) {
			overlapX, overlapY := getOverlap(playerRect, objRect)

			if math.Abs(float64(overlapX)) < math.Abs(float64(overlapY)) {
				// Push in X direction
				p.x -= overlapX
			} else {
				// Push in Y direction
				p.y -= overlapY
			}

			// Update player rect after pushing
			playerRect = *p.GetHitBoxRect()
		}
	}

	// Tool use
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
	if p.inventory.Tools[p.inventory.SelectedIndex].Tool == ui.Axe {
		for _, t := range p.trees {
			if rl.CheckCollisionRecs(*p.HitBox, *t.GetHitBoxRect()) {
				t.Damage()
			}
		}
	}
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

func (p *Player) GetHitBoxRect() *rl.Rectangle {
	return &rl.Rectangle{
		X:      float32(p.x),
		Y:      float32(p.y + (float32(p.height)/2 - 10)),
		Width:  float32(p.width),
		Height: float32(p.height/2) + 10,
	}
}
