// rpg-tutorial, LitFill <litfill at litfill dot site>
// program for me to learn buildong a game in pure go in windows.
package main

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"log/slog"
	"math"
	"os"

	"github.com/LitFill/fatal"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Number interface {
	int | float64
}

const (
	WIDTH  = 640
	HEIGHT = 480
)

type Vec2[N Number] struct {
	X, Y N
}

func (v *Vec2[N]) Add(v2 *Vec2[N]) *Vec2[N] {
	return &Vec2[N]{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
	}
}

func (v *Vec2[N]) Sub(v2 *Vec2[N]) *Vec2[N] {
	return &Vec2[N]{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
	}
}

func (v *Vec2[N]) Mul(scalar N) *Vec2[N] {
	return &Vec2[N]{
		X: v.X * scalar,
		Y: v.Y * scalar,
	}
}

func (v *Vec2[N]) Div(scalar N) *Vec2[N] {
	if scalar != 0 {
		return &Vec2[N]{
			X: v.X / scalar,
			Y: v.Y / scalar,
		}
	}
	return &Vec2[N]{
		X: N(math.Inf(0)),
		Y: N(math.Inf(0)),
	}
}

func (v *Vec2[N]) Add_nr(v2 *Vec2[N]) { *v = *v.Add(v2) }
func (v *Vec2[N]) Sub_nr(v2 *Vec2[N]) { *v = *v.Sub(v2) }
func (v *Vec2[N]) Zero()              { v.Sub_nr(v) }
func (v *Vec2[N]) Len() N             { return N(math.Sqrt(float64(v.X*v.X + v.Y*v.Y))) }

func (v *Vec2[N]) Normalize() *Vec2[N] {
	length := v.Len()
	if length != 0 {
		return v.Div(length)
	}
	return &Vec2[N]{0, 0}
}

func (v *Vec2[N]) Sign() *Vec2[N] {
	return &Vec2[N]{
		X: sign(v.X),
		Y: sign(v.Y),
	}
}

func NewVec2[N Number](x, y N) *Vec2[N] { return &Vec2[N]{x, y} }
func NewVec2Zero[N Number]() *Vec2[N]   { return NewVec2[N](0, 0) }

var MOVE_LEFT = NewVec2(-2.0, 0.0)
var MOVE_RIGHT = NewVec2(2.0, 0.0)
var MOVE_UP = NewVec2(0.0, -2.0)
var MOVE_DOWN = NewVec2(0.0, 2.0)

type Sprite struct {
	Img      *ebiten.Image
	Pos      *Vec2[float64]
	Velocity *Vec2[float64]
}

type Enemy struct {
	*Sprite
	IsFollowPlayer bool
}

type Game struct {
	Player  *Sprite
	Enemies []*Enemy
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return fmt.Errorf("game closed")
	}

	g.handlePlayerMovement()

	for _, enemy := range g.Enemies {
		if enemy.IsFollowPlayer {
			dist := g.Player.Pos.Sub(enemy.Pos).Sign()
			enemy.Velocity = enemy.Velocity.Add(dist).Sign()

			enemy.Pos.Add_nr(enemy.Velocity)
		}
	}

	return nil
}

func (g *Game) handlePlayerMovement() {
	var (
		w, h      = g.Layout(WIDTH, HEIGHT)
		moveLeft  = ebiten.IsKeyPressed(ebiten.KeyH) || ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft)
		moveRight = ebiten.IsKeyPressed(ebiten.KeyL) || ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight)
		moveUp    = ebiten.IsKeyPressed(ebiten.KeyK) || ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp)
		moveDown  = ebiten.IsKeyPressed(ebiten.KeyJ) || ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown)
		arah      = NewVec2Zero[float64]()
	)

	if moveLeft && g.Player.Pos.X > 0 {
		arah.Add_nr(MOVE_LEFT)
	}
	if moveRight && g.Player.Pos.X < float64(w-16) {
		arah.Add_nr(MOVE_RIGHT)
	}
	if moveUp && g.Player.Pos.Y > 0 {
		arah.Add_nr(MOVE_UP)
	}
	if moveDown && g.Player.Pos.Y < float64(h-16) {
		arah.Add_nr(MOVE_DOWN)
	}

	g.Player.Velocity.Add_nr(arah)
	g.Player.Pos.Add_nr(g.Player.Velocity)
	g.Player.Velocity.Zero()
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return 320, 240
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(g.Player.Pos.X, g.Player.Pos.Y)

	screen.DrawImage(
		g.Player.Img.SubImage(image.Rect(
			0, 0,
			16, 16,
		)).(*ebiten.Image),
		opts,
	)
	opts.GeoM.Reset()

	for _, sprite := range g.Enemies {
		opts.GeoM.Translate(sprite.Pos.X, sprite.Pos.Y)

		screen.DrawImage(
			sprite.Img.SubImage(image.Rect(
				0, 0,
				16, 16,
			)).(*ebiten.Image),
			opts,
		)
		opts.GeoM.Reset()
	}
}

func main() {
	logFile := fatal.CreateLogFile("log.json")
	defer logFile.Close()
	errLogger := fatal.CreateLogger(io.MultiWriter(logFile, os.Stderr), slog.LevelError)
	infoLogger := fatal.CreateLogger(io.MultiWriter(logFile, os.Stdout), slog.LevelInfo)

	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("LitFill's rpg")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	playerImg, _, err := ebitenutil.NewImageFromFile("./assets/images/NinjaBlue2/SpriteSheet.png")
	fatal.Log(err, errLogger, "cannot load playerImg")

	beastImg, _, err := ebitenutil.NewImageFromFile("./assets/images/Beast/Beast.png")
	fatal.Log(err, errLogger, "cannot load playerImg")

	infoLogger.Info("Running game")
	game := &Game{
		Player: &Sprite{
			Img:      playerImg,
			Pos:      NewVec2(50.0, 50),
			Velocity: NewVec2(0.0, 0),
		},
		Enemies: []*Enemy{
			{
				Sprite: &Sprite{
					Img:      beastImg,
					Pos:      NewVec2(50.0, 100),
					Velocity: NewVec2(0.0, 0),
				},
				IsFollowPlayer: true,
			},
			{
				Sprite: &Sprite{
					Img:      beastImg,
					Pos:      NewVec2(100.0, 100),
					Velocity: NewVec2(0.0, 0),
				},
				IsFollowPlayer: false,
			},
			{
				Sprite: &Sprite{
					Img:      beastImg,
					Pos:      NewVec2(150.0, 100),
					Velocity: NewVec2(0.0, 0),
				},
				IsFollowPlayer: true,
			},
		},
	}

	fatal.Log(ebiten.RunGame(game),
		errLogger, "cannot run the game",
		"game", game,
	)

	infoLogger.Info("Closing game", "game state", game)
}
