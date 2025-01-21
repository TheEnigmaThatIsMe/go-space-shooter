package main

import (
	"go-space-shooter/game"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := game.NewGame()

	// Start the game
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
