package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Create and run the game
	game := NewGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Flappy Bird")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
