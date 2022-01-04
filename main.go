package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	spaceship  = NewObject("textures/spaceship/spaceship_idle.png")
	is_spawned = false
)

type Object struct {
	image   *ebiten.Image
	options *ebiten.DrawImageOptions
	x, y    float64
}

func (obj *Object) SetCoord(x, y float64) {
	obj.x, obj.y = x, y
}

type Game struct{}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(1)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !is_spawned {
		Spawn(screen)
	}
	// SPACESHIP UPDATE
	x, y, ok := CheckCursorPosition(screen.Bounds(), spaceship, 0, (spaceship.image.Bounds().Max.X / 10), 0, spaceship.image.Bounds().Max.Y/10)
	ebitenutil.DebugPrint(screen, fmt.Sprint(x, y, ebiten.CurrentFPS()))
	if ok {
		rx, ry := spaceship.TranslateTo(float64(x), float64(y))
		spaceship.options.GeoM.Translate(rx, ry)
		spaceship.SetCoord(float64(x), float64(y))
	}
	ebitenutil.DebugPrintAt(screen, "SPACEGO v0.1\nExit -> ESC", 500, 500)
	screen.DrawImage(spaceship.image, spaceship.options)
	// END
}

func Spawn(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "abcd")
	spaceship.options.GeoM.Scale(0.1, 0.1)
	spaceship.options.GeoM.Translate(0, 0)
	spaceship.x = 0
	spaceship.y = 0
	screen.DrawImage(spaceship.image, spaceship.options)
	is_spawned = true
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720
}

func main() {
	ebiten.SetFullscreen(true)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Render an image")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func LoadTexture(path string) *ebiten.Image {
	ret, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		panic(err)
	}
	return ret
}

func CheckCursorPosition(bounds image.Rectangle, obj *Object, right, left, bottom, top int) (int, int, bool) {
	x, y := ebiten.CursorPosition()
	err := 0
	if !(x >= bounds.Min.X+right && x <= bounds.Max.X-left) {
		x = bounds.Max.X - (obj.image.Bounds().Max.X / 10)
		err++
	}
	if !(y >= bounds.Min.Y+bottom && y <= bounds.Max.Y-top) {
		y = bounds.Max.Y - (spaceship.image.Bounds().Max.Y / 10)
		err++
	}
	return x, y, err != 2
}

func NewObject(path string) *Object {
	return &Object{image: LoadTexture(path), options: &ebiten.DrawImageOptions{}}
}

func (obj *Object) TranslateTo(x, y float64) (rx, ry float64) {
	if x < obj.x {
		rx = -1 * (obj.x - x)
	}
	if x > obj.x {
		rx = x - obj.x
	}
	if y < obj.y {
		ry = -1 * (obj.y - y)
	}
	if y > obj.y {
		ry = y - obj.y
	}
	return
}
