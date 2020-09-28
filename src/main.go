package main

import (
	"flag"
	"log"

	"github.com/Omar-V2/game-of-life/pkg/gameoflife"

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

func main() {
	flag.Parse()

	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Welcome to the Game of Life!")

	game := gameoflife.NewGame(width, height, *density, *interactive, false)
	if !*interactive {
		game.NewLife()
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
