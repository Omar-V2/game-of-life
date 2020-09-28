package main

import (
	"flag"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/hajimehoshi/ebiten"
)

const (
	width  = 640
	height = 480
)

var (
	density     = flag.Int("density", 10, "How densely populated the board will be on the first generation, lower is more dense")
	interactive = flag.Bool("interactive", false, "Starts the game in an interactive mode where you can set the starting state via mouse interaction.")
	begin       = false
)

type Game struct {
	cells  [][]bool
	width  int
	height int
}

// creates a new Game with all cells dead as the starting state
func newGame(width, height int) *Game {
	cells := make([][]bool, height)
	for i := range cells {
		cells[i] = make([]bool, width)
	}
	return &Game{cells: cells, width: width, height: height}
}

// randomly set some of the cells in the game to the alive state
func (g *Game) newLife() {
	rand.Seed(time.Now().UnixNano())
	cells := g.cells
	for i := range cells {
		for j := range cells[i] {
			if rand.Intn(*density) == 1 {
				cells[i][j] = true
			}
		}
	}
}

// sets all cells in the game to the dead state
func (g *Game) clearLife() {
	for i := range g.cells {
		for j := range g.cells[i] {
			g.cells[i][j] = false
		}
	}
}

// transitions the game to the next state using the game's rules:
// https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life
func (g *Game) nextState() {
	nextCells := newGame(width, height).cells
	for y := range g.cells {
		for x := range g.cells[y] {
			count := g.countNeighbours(x, y)
			alive := g.cells[y][x]
			if alive && count < 2 {
				nextCells[y][x] = false
			}
			if alive && (count == 2 || count == 3) {
				nextCells[y][x] = true
			}
			if alive && count > 3 {
				nextCells[y][x] = false
			}
			if !alive && count == 3 {
				nextCells[y][x] = true
			}
		}
	}
	g.cells = nextCells
}

// counts the number of neighbours of the cell at position x, y on the board
func (g *Game) countNeighbours(x, y int) int {
	directions := []int{0, 1, -1}
	count := 0
	for _, dX := range directions {
		for _, dY := range directions {
			// this is the current cell, we only care about neighbours
			if dX == 0 && dY == 0 {
				continue
			}
			if inBounds(x+dX, y+dY) && g.cells[y+dY][x+dX] == true {
				count++
			}
		}
	}
	return count
}

// checks if a given co-ordinate is in the bounds of the board's dimensions
func inBounds(x, y int) bool {
	xOk := x >= 0 && x < width
	yOk := y >= 0 && y < height
	return xOk && yOk
}

// Sets the cell at the position that the cursor was left clicked to alive
func (g *Game) interact() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if inBounds(y, x) {
			g.cells[x][y] = true
		}
	}
}

// run the game in interactive mode
func (g *Game) interactiveMode(begin bool) {
	if begin {
		g.nextState()
	} else {
		g.interact()
	}
}

func (g *Game) Update(screen *ebiten.Image) error {
	if *interactive {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			begin = true
		} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
			g.clearLife()
			begin = false
		}
		g.interactiveMode(begin)
	} else {
		g.nextState()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := range g.cells {
		for j := range g.cells[i] {
			if g.cells[i][j] {
				ebitenutil.DrawRect(screen, float64(i), float64(j), 1, 1, color.White)
			}
		}
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	flag.Parse()

	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Welcome to the Game of Life!")

	game := newGame(width, height)
	if !*interactive {
		game.newLife()
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
