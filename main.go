package main

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const ScreenWidth = 320
const ScreenHeight = 568

type game struct {
}

type otter struct {
	Life         int
	Survived     time.Duration
	Box          rl.Rectangle
	Pos          rl.Vector2
	Texture      rl.Texture2D
	currentFrame int32
}

func NewOtter() *otter {
	// otterImg := rl.LoadImage("assets/img/otter_swim.png") // Load image in CPU memory (RAM)
	// rl.ImageResize(otterImg, otterImg.Width/2, otterImg.Height/2)
	// texture := rl.LoadTextureFromImage(otterImg)
	texture := rl.LoadTexture("assets/img/otter_swim.png")
	startingPosition := rl.NewVector2(
		ScreenWidth/2,
		ScreenHeight,
	)

	return &otter{
		Life:    3,
		Box:     rl.NewRectangle(startingPosition.X-25, startingPosition.Y-115, 50, 115),
		Texture: texture,
		Pos:     startingPosition,
	}
}

func (o *otter) Draw(frameCounter int) {
	rl.DrawText(fmt.Sprintf("%d", frameCounter), 0, 0, 12, rl.Black)
	if frameCounter%8 == 0 {
		o.currentFrame = o.currentFrame + 1%12
	}

	scale := float32(0.5)
	frameRec := rl.NewRectangle(float32(o.currentFrame*o.Texture.Width/11), 0.0, float32(o.Texture.Width/12), float32(o.Texture.Height))
	destRec := rl.NewRectangle(o.Pos.X-37, o.Pos.Y-120, float32(o.Texture.Width/12)*scale, float32(o.Texture.Height)*scale)

	rl.DrawRectangleRec(o.Box, rl.Red)
	rl.DrawTexturePro(o.Texture, frameRec, destRec, rl.NewVector2(0, 0), 0.5, rl.White)
}

type swimLane struct {
	Pos          rl.Vector2
	RockTexture  rl.Texture2D
	LogTexture   rl.Texture2D
	currentFrame int32
}

func NewSwimLane(Pos rl.Vector2) *swimLane {
	rockTexture := rl.LoadTexture("assets/img/rock.png")
	logTexture := rl.LoadTexture("assets/img/log.png")
	return &swimLane{Pos: Pos, RockTexture: rockTexture, LogTexture: logTexture}
}

func (s *swimLane) Draw() {
	rl.DrawRectangleRec(rl.NewRectangle(s.Pos.X, s.Pos.Y, 25, 25), rl.Red)
}

func (s *swimLane) SpawnLog(frameCounter int) {
	totalFrame := int32(8)
	if frameCounter%8 == 0 {
		s.currentFrame = s.currentFrame + 1%totalFrame
	}

	scale := float32(0.5)
	frameRec := rl.NewRectangle(float32(s.currentFrame*s.LogTexture.Width/totalFrame), 0.0, float32(s.LogTexture.Width/totalFrame), float32(s.LogTexture.Height))
	destRec := rl.NewRectangle(s.Pos.X-37, s.Pos.Y-120, float32(s.LogTexture.Width/12)*scale, float32(s.LogTexture.Height)*scale)

	rl.DrawTexturePro(s.LogTexture, frameRec, destRec, rl.NewVector2(0, 0), 0.5, rl.White)
}

func (s *swimLane) SpawnRock(frameCounter int) {
	// if frameCounter%8 == 0 {
	// 	s.currentFrame = s.currentFrame + 1%4
	// }
	// scale := float32(0.5)
	// frameRec := rl.NewRectangle(float32(o.currentFrame*o.Texture.Width/11), 0.0, float32(o.Texture.Width/12), float32(o.Texture.Height))
	// destRec := rl.NewRectangle(o.Pos.X-37, o.Pos.Y-120, float32(o.Texture.Width/12)*scale, float32(o.Texture.Height)*scale)
	//
	// rl.DrawRectangleRec(o.Box, rl.Red)
	// rl.DrawTexturePro(o.Texture, frameRec, destRec, rl.NewVector2(0, 0), 0.5, rl.White)

}

// type log struct {
// 	Texture rl.Texture2D
// 	Pos     rl.Vector2
// 	// Lane is the position
// 	Lane int
// }
//
// func NewLog() *log {
// 	texture := rl.LoadTexture("assets/img/otter_swim.png")
// 	startingPosition := rl.NewVector2(
// 		ScreenWidth/2,
// 		ScreenHeight,
// 	)
// 	return &log{
// 		Texture: texture,
// 		Pos:     startingPosition,
// 	}
// }
// func (l *log) Spawn(frameCounter int) {
//
// }

func main() {
	rl.InitWindow(ScreenWidth, ScreenHeight, "Ottermatic")
	rl.SetTargetFPS(60)

	// Init Game state
	o := NewOtter()
	swimlanes := []*swimLane{}
	swimlaneOffset := float32(42)
	for i := 0; i < 3; i++ {
		swimlanes = append(swimlanes, NewSwimLane(rl.NewVector2(float32(i*ScreenWidth/3)+swimlaneOffset, 0)))
	}
	frameCounter := 0

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		o.Draw(frameCounter)
		for i := 0; i < 3; i++ {
			swimlanes[i].Draw()
		}

		rl.ClearBackground(rl.DarkBlue)

		rl.EndDrawing()
		frameCounter++
	}

	rl.CloseWindow()
}

func Game() {

}
