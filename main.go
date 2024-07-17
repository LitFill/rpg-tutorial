// rpg-tutorial, LitFill <litfill at litfill dot site>
// program for me to learn buildong a game in pure go in windows.
package main

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"log/slog"
	"os"

	"github.com/LitFill/fatal"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	WIDTH  = 640
	HEIGHT = 480
)

type Vec2[N ~int | ~float64] struct {
	X, Y N
}

func NewVec2[N ~int | ~float64](x, y N) *Vec2[N] { return &Vec2[N]{x, y} }

type Sprite struct {
	Img      *ebiten.Image
	Pos      *Vec2[float64]
	Velocity *Vec2[float64]
}

type Game struct {
	Player  *Sprite
	Sprites []*Sprite
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return fmt.Errorf("game closed")
	}

	g.handleMovement()

	return nil
}

func (g *Game) handleMovement() {
	var (
		moveLeft  = ebiten.IsKeyPressed(ebiten.KeyH) || ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft)
		moveRight = ebiten.IsKeyPressed(ebiten.KeyL) || ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight)
		moveUp    = ebiten.IsKeyPressed(ebiten.KeyK) || ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp)
		moveDown  = ebiten.IsKeyPressed(ebiten.KeyJ) || ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown)
	)

	if moveLeft {
		g.Player.Velocity.X -= 2
	}
	if moveRight {
		g.Player.Velocity.X += 2
	}
	if moveUp {
		g.Player.Velocity.Y -= 2
	}
	if moveDown {
		g.Player.Velocity.Y += 2
	}

	w, h := g.Layout(WIDTH, HEIGHT)

	g.Player.Pos.X = clamp(g.Player.Pos.X+g.Player.Velocity.X, 0, float64(w-16))
	g.Player.Pos.Y = clamp(g.Player.Pos.Y+g.Player.Velocity.Y, 0, float64(h-16))

	g.Player.Velocity.X = 0
	g.Player.Velocity.Y = 0
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

	for _, sprite := range g.Sprites {
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
		Sprites: []*Sprite{
			{
				Img:      beastImg,
				Pos:      NewVec2(50.0, 100),
				Velocity: NewVec2(0.0, 0),
			},
			{
				Img:      beastImg,
				Pos:      NewVec2(100.0, 100),
				Velocity: NewVec2(0.0, 0),
			},
			{
				Img:      beastImg,
				Pos:      NewVec2(150.0, 100),
				Velocity: NewVec2(0.0, 0),
			},
		},
	}

	fatal.Log(ebiten.RunGame(game),
		errLogger, "cannot run the game",
		"game", game,
	)

	infoLogger.Info("Closing game")
}
