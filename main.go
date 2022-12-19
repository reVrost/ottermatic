package main

import (
	"fmt"
	"math/rand"
	"time"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// TODO:
// M1
// Compile to WASM

// M2
// Rain Thunder Animation
// Rock
// Pause
// Leaderboard

const ScreenWidth = 320
const ScreenHeight = 568
const SpawnRandomness = 8
const Debug = false

type otter struct {
	Life        int32
	Survived    time.Duration
	Pos         rl.Vector2
	Hitbox      rl.Rectangle
	Texture     rl.Texture2D
	HitTexture  rl.Texture2D
	DeadTexture rl.Texture2D

	HitSound rl.Sound

	currentFrame int32
	hasCollided  bool
	frameCounter int32
	CollideFrame int32
	deadFrame    int32
}

func NewOtter() *otter {
	// Textures
	texture := rl.LoadTexture("assets/img/otter_swim.png")
	hitTexture := rl.LoadTexture("assets/img/otter_hit.png")
	deadTexture := rl.LoadTexture("assets/img/otter_dead.png")
	hitSound := rl.LoadSound("assets/sfx/hit.wav")

	startingPosition := rl.NewVector2(
		ScreenWidth/2,
		ScreenHeight+20,
	)

	return &otter{
		Life:        3,
		Texture:     texture,
		HitTexture:  hitTexture,
		DeadTexture: deadTexture,
		Pos:         startingPosition,
		HitSound:    hitSound,
	}
}

func (o *otter) Update(state GameState) {
	// hit box
	o.Hitbox = rl.NewRectangle(o.Pos.X-25, o.Pos.Y-119, 41, 100)

	// Draw
	if state == End {
		scale := float32(0.43)
		noOfFrames := int32(8)
		frameRec := rl.NewRectangle(float32(o.deadFrame*o.DeadTexture.Width/noOfFrames), 0.0, float32(o.DeadTexture.Width/noOfFrames), float32(o.DeadTexture.Height))
		destRec := rl.NewRectangle(o.Pos.X-37, o.Pos.Y-120, float32(o.DeadTexture.Width/noOfFrames)*scale, float32(o.DeadTexture.Height)*scale)
		rl.DrawTexturePro(o.DeadTexture, frameRec, destRec, rl.NewVector2(0, 0), 0.5, rl.White)

		if o.frameCounter%8 == 0 {
			o.deadFrame = o.deadFrame % noOfFrames
			o.deadFrame++
		}
	} else if o.hasCollided {

		scale := float32(0.48)
		noOfFrames := int32(5)
		frameRec := rl.NewRectangle(float32(o.CollideFrame*o.HitTexture.Width/noOfFrames), 0.0, float32(o.HitTexture.Width/noOfFrames), float32(o.HitTexture.Height))
		destRec := rl.NewRectangle(o.Pos.X-37, o.Pos.Y-130, float32(o.HitTexture.Width/noOfFrames)*scale, float32(o.HitTexture.Height)*scale)
		rl.DrawTexturePro(o.HitTexture, frameRec, destRec, rl.NewVector2(0, 0), 0.5, rl.White)

		if o.frameCounter%3 == 0 {
			o.CollideFrame = o.CollideFrame % noOfFrames
			if o.CollideFrame == 4 {
				o.hasCollided = false
				o.CollideFrame = -1
			}
			o.CollideFrame++
		}
	} else {

		if Debug {
			rl.DrawRectangleRec(o.Hitbox, rl.Red)
		}

		// Otter
		scale := float32(0.43)
		noOfFrames := int32(12)
		frameRec := rl.NewRectangle(float32(o.currentFrame*o.Texture.Width/noOfFrames), 0.0, float32(o.Texture.Width/noOfFrames), float32(o.Texture.Height))
		destRec := rl.NewRectangle(o.Pos.X-37, o.Pos.Y-120, float32(o.Texture.Width/noOfFrames)*scale, float32(o.Texture.Height)*scale)
		rl.DrawTexturePro(o.Texture, frameRec, destRec, rl.NewVector2(0, 0), 0.5, rl.White)

		if o.frameCounter%8 == 0 {
			o.currentFrame = o.currentFrame % noOfFrames
			o.currentFrame++
		}
	}
	o.frameCounter++
}

func (o *otter) Collide() {
	rl.PlaySound(o.HitSound)
	o.hasCollided = true
	o.CollideFrame = -1
	o.Life--
}

type swimLane struct {
	Pos          rl.Vector2
	ObjectHitbox rl.Rectangle
	RockTexture  rl.Texture2D
	LogTexture   rl.Texture2D

	animationFrame int32
	framesCounter  int32
	speedCounter   int32
	// delay for the next ob to spawn next
	interval    int32
	hasCollided bool
}

const swimSpeed = float32(5)
const hitBoxOffset = float32(25)

// TODO: implement classic mode
func (o *otter) SwimRight() {
	mostRight := o.Pos.X + swimSpeed + hitBoxOffset
	if mostRight > ScreenWidth {
		return
	}
	o.Pos.X += swimSpeed
}
func (o *otter) SwimLeft() {
	mostLeft := o.Pos.X - swimSpeed - hitBoxOffset
	if mostLeft < 0 {
		return
	}
	o.Pos.X -= swimSpeed
}

const swimLaneYOffset = float32(-15)

func NewSwimLane(x float32) *swimLane {
	rockTexture := rl.LoadTexture("assets/img/rock.png")
	logTexture := rl.LoadTexture("assets/img/log.png")
	return &swimLane{Pos: rl.NewVector2(x, ScreenHeight+10), RockTexture: rockTexture, LogTexture: logTexture, interval: int32(rand.Intn(5))}
}

func (s *swimLane) Collide() {
	s.Pos.Y += ScreenHeight
}

func (s *swimLane) Update(state GameState, rainCounter float32) {
	// animation
	totalFrames := int32(8)
	animationFrame := float32(s.animationFrame % totalFrames)
	imgScale := float32(0.5)
	logWidth := float32(s.LogTexture.Width / totalFrames)
	logHeight := float32(s.LogTexture.Height)

	// movement
	frameSpeed := int32(12)
	moveSpeed := float32(s.speedCounter * int32(5))

	moveSpeed = moveSpeed * (1 + (0.5 * rainCounter))

	posY := s.Pos.Y + moveSpeed

	// denotes if object is not on screen
	isEmpty := posY >= ScreenHeight

	frameRec := rl.NewRectangle(animationFrame*logWidth, 0.0, float32(s.LogTexture.Width/totalFrames), logHeight)
	destRec := rl.NewRectangle(s.Pos.X-(logWidth/5), posY, logWidth*imgScale, logHeight*imgScale)

	// hitbox
	s.ObjectHitbox = rl.NewRectangle(s.Pos.X-25, posY+23, (logWidth*imgScale)-10, logHeight*imgScale/2)
	if Debug {
		// Swimlane pos
		rl.DrawRectangleRec(rl.NewRectangle(s.Pos.X, s.Pos.Y, 8, 8), rl.Red)
		// Object pos
		rl.DrawRectangleRec(s.ObjectHitbox, rl.Red)
	}

	rl.DrawTexturePro(s.LogTexture, frameRec, destRec, rl.NewVector2(0, 0), 0.5, rl.White)

	s.framesCounter++
	if state != End && state != Pause {
		s.speedCounter++
	}
	if s.framesCounter >= 60/frameSpeed {
		s.animationFrame++
		s.framesCounter = 0
	}

	if isEmpty && s.speedCounter >= 60*s.interval {
		s.interval = int32(rand.Intn(SpawnRandomness))
		s.speedCounter = 0
		s.Pos.Y = swimLaneYOffset
	}
}

// Direction - Custom type to hold value for week day ranging from 1-4
type GameState int

// Declare related constants for each direction starting with index 1
const (
	Start GameState = iota + 1
	Rain
	Pause
	End
)

type game struct {
	State        GameState
	Life         int
	FrameCounter int
	RainCounter  int32

	gameOverSound rl.Sound
	rainSound     rl.Sound
	windSound     rl.Sound

	heartTexture     rl.Texture2D
	ElapsedInSeconds int32
}

func NewGame() *game {
	gameOverSound := rl.LoadSound("assets/sfx/game_over.wav")
	heartTexture := rl.LoadTexture("assets/img/heart.png")
	rainSound := rl.LoadSound("assets/sfx/heavyrain.wav")
	windSound := rl.LoadSound("assets/sfx/windsoar.wav")

	return &game{
		Life:          3,
		gameOverSound: gameOverSound,
		rainSound:     rainSound,
		windSound:     windSound,
		heartTexture:  heartTexture,
	}
}

func (g *game) End() {
	rl.PlaySound(g.gameOverSound)
	g.State = End
}

func (g *game) Reset() {
	g.State = Start
	g.FrameCounter = 0
	g.Life = 3
}

func (g *game) Update() {

	topMargin := int32(5)

	rl.DrawText(fmt.Sprintf("%06d m", g.FrameCounter), (ScreenWidth/2)-28, topMargin, 18, rl.Yellow)
	rl.DrawText(fmt.Sprintf("%d sec", g.ElapsedInSeconds), 10, topMargin, 18, rl.Red)

	imgScale := float32(0.18)
	frameRec := rl.NewRectangle(0.0, 0.0, float32(g.heartTexture.Width), float32(g.heartTexture.Height))
	for i := 0; i < g.Life; i++ {
		destRec := rl.NewRectangle((ScreenWidth/2)+73+float32(i*28), float32(topMargin), float32(g.heartTexture.Width)*imgScale, float32(g.heartTexture.Height)*imgScale)
		rl.DrawTexturePro(g.heartTexture, frameRec, destRec, rl.NewVector2(0, 0), 0.5, rl.White)
	}

	if g.State != End && g.State != Pause {
		g.FrameCounter++
		if g.FrameCounter%60 == 0 {
			g.ElapsedInSeconds++
			if g.ElapsedInSeconds%30 == 0 {
				g.State = Rain
				g.RainCounter++
				fmt.Println(g.ElapsedInSeconds)
				rl.PlaySound(g.windSound)
				rl.PlaySound(g.rainSound)
			}
		}
	}

}

func main() {
	rl.InitWindow(ScreenWidth, ScreenHeight, "Ottermatic")
	rl.SetTargetFPS(60)
	rl.InitAudioDevice()

	// Init Game state
	touchArea := rl.NewRectangle(0, 0, ScreenWidth, ScreenHeight)
	game := NewGame()

	// Otter
	otter := NewOtter()

	// Swimlanes
	swimLanes := []*swimLane{}
	swimlaneOffsetX := float32(50)
	for i := 0; i < 3; i++ {
		swimLanes = append(swimLanes, NewSwimLane(float32(i*ScreenWidth/3)+swimlaneOffsetX))
	}

	// Music
	// music := rl.LoadMusicStream("assets/bgm/ericanoshi_4.2.wav")
	music := rl.LoadMusicStream("assets/bgm/storm_and_sunshine_7.7.wav")
	rl.PlayMusicStream(music)

	// Game loop
	for !rl.WindowShouldClose() {
		//Music
		rl.UpdateMusicStream(music)

		// Game
		rl.ClearBackground(rl.DarkBlue)
		game.Update()

		// Control logic
		if game.State != End && game.State != Pause {
			currentGesture := rl.GetGestureDetected()
			touchPosition := rl.GetTouchPosition(0)

			touchLeft := rl.CheckCollisionPointRec(touchPosition, touchArea) &&
				(currentGesture == rl.GestureTap || currentGesture == rl.GestureHold) && touchPosition.X < ScreenWidth/2
			touchRight := rl.CheckCollisionPointRec(touchPosition, touchArea) &&
				(currentGesture == rl.GestureTap || currentGesture == rl.GestureHold) && touchPosition.X > ScreenWidth/2

			if rl.IsKeyDown(rl.KeyLeft) || touchLeft {
				otter.SwimLeft()
			} else if rl.IsKeyDown(rl.KeyRight) || touchRight {
				otter.SwimRight()
			}
		}

		rl.BeginDrawing()

		otter.Update(game.State)

		for _, sl := range swimLanes {
			// debugging
			sl.Update(game.State, float32(game.RainCounter))

			// game mechanics
			if rl.CheckCollisionRecs(otter.Hitbox, sl.ObjectHitbox) {
				sl.Collide()
				game.Life--
				if game.Life != 0 {
					otter.Collide()
				}
			}
		}

		// Check game end
		if game.State != End && game.Life <= 0 {
			game.End()
			rl.StopMusicStream(music)
		}

		// Restart menu
		if game.State == End {
			isPressed := rg.Button(rl.NewRectangle(ScreenWidth/2-75, ScreenHeight/2-12, 150, 24), "Retry")
			rl.DrawText(fmt.Sprintf("%06d", game.FrameCounter), ScreenWidth/2-28, ScreenHeight/2-50, 18, rl.Green)
			rg.LabelEx(rl.NewRectangle(ScreenWidth/2-60, ScreenHeight/2-80, 120, 24), "Your score:", rl.Green, rl.NewColor(0, 0, 0, 0), rl.NewColor(0, 0, 0, 0))
			if isPressed {
				for _, sl := range swimLanes {
					sl.Collide()
				}
				game.Reset()
				rl.PlayMusicStream(music)
				fmt.Println("Pressed")
			}
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
