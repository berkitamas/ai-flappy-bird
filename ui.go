package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"math"
)

// Draw draws the game state to the screen
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw background
	g.drawBackground(screen)

	// Draw game elements (bird, pipes)
	g.drawGameElements(screen)

	// Draw UI based on game state
	switch g.state {
	case stateMenu:
		g.drawMenu(screen)
	case statePlay:
		g.drawGameUI(screen)
	}
}

// drawBackground draws the sky, clouds, and ground
func (g *Game) drawBackground(screen *ebiten.Image) {
	// Draw sky background
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(g.backgroundImg, op)

	// Draw clouds
	for _, cloud := range g.clouds {
		op := &ebiten.DrawImageOptions{}

		// Draw cloud as a white ellipse
		cloudImg := ebiten.NewImage(int(cloud.width), int(cloud.height))

		// Draw elliptical cloud
		for y := 0; y < int(cloud.height); y++ {
			for x := 0; x < int(cloud.width); x++ {
				// Calculate normalized coordinates
				nx := float64(x)/cloud.width*2 - 1
				ny := float64(y)/cloud.height*2 - 1

				// Ellipse equation: (x/a)² + (y/b)² <= 1
				if nx*nx+ny*ny*4 <= 1 {
					cloudImg.Set(x, y, cloudColor)
				}
			}
		}

		op.GeoM.Translate(cloud.x, cloud.y)
		screen.DrawImage(cloudImg, op)
	}

	// Draw ground
	groundY := float64(screenHeight - 30)
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.groundOffset, groundY)
	screen.DrawImage(g.groundImg, op)

	// Draw a second ground image to create a seamless scrolling effect
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.groundOffset+float64(screenWidth), groundY)
	screen.DrawImage(g.groundImg, op)
}

// drawGameElements draws the bird and pipes
func (g *Game) drawGameElements(screen *ebiten.Image) {
	// Draw bird
	op := &ebiten.DrawImageOptions{}

	// Apply rotation based on velocity for more natural movement
	rotationAngle := math.Atan2(g.bird.velocity, pipeSpeed) // Use velocity to determine angle

	// Rotate around the center of the bird
	op.GeoM.Translate(-float64(g.bird.width)/2, -float64(g.bird.height)/2) // Move to center
	op.GeoM.Rotate(rotationAngle)                                          // Apply rotation
	op.GeoM.Translate(float64(g.bird.width)/2, float64(g.bird.height)/2)   // Move back

	// Set bird position
	op.GeoM.Translate(g.bird.x, g.bird.y)

	// Draw the bird image
	screen.DrawImage(g.birdImg, op)

	// Draw pipes
	for _, pipe := range g.pipes {
		pipeX := float64(pipe.x)

		// Draw top pipe (upside down)
		topPipeHeight := float64(pipe.gapStart)

		// Pipe body
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Scale(1, topPipeHeight/100) // Scale to correct height (pipe image is 100px tall)
		op.GeoM.Translate(pipeX, 0)
		screen.DrawImage(g.pipeImg, op)

		// Pipe cap
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(pipeX-5, topPipeHeight-20) // Position cap at bottom of top pipe
		screen.DrawImage(g.pipeCapImg, op)

		// Draw bottom pipe
		bottomPipeY := float64(pipe.gapStart + pipeGap)
		bottomPipeHeight := float64(screenHeight - pipe.gapStart - pipeGap)

		// Pipe body
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Scale(1, bottomPipeHeight/100) // Scale to correct height
		op.GeoM.Translate(pipeX, bottomPipeY)
		screen.DrawImage(g.pipeImg, op)

		// Pipe cap
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(pipeX-5, bottomPipeY) // Position cap at top of bottom pipe
		screen.DrawImage(g.pipeCapImg, op)
	}
}

// drawMenu draws the main menu UI
func (g *Game) drawMenu(screen *ebiten.Image) {
	// Draw semi-transparent overlay to make menu text more visible
	overlayColor := color.RGBA{0, 0, 0, 128} // Semi-transparent black
	ebitenutil.DrawRect(screen, 0, 0, float64(screenWidth), float64(screenHeight), overlayColor)

	// Draw game title
	titleText := "FLAPPY BIRD"
	ebitenutil.DebugPrintAt(screen, titleText, screenWidth/2-50, screenHeight/4)

	// Draw menu options
	startText := "Start Game"
	exitText := "Exit"

	// Highlight selected option
	if g.menuOption == 0 {
		startText = "> " + startText + " <"
	} else {
		exitText = "> " + exitText + " <"
	}

	ebitenutil.DebugPrintAt(screen, startText, screenWidth/2-50, screenHeight/2)
	ebitenutil.DebugPrintAt(screen, exitText, screenWidth/2-50, screenHeight/2+30)

	// Draw controls help
	controlsText := "Use UP/DOWN to navigate, SPACE/ENTER to select"
	ebitenutil.DebugPrintAt(screen, controlsText, screenWidth/2-150, screenHeight*3/4)

	// Draw demo indicator
	demoText := "DEMO PLAYING"
	ebitenutil.DebugPrintAt(screen, demoText, screenWidth-100, 30)

	// Draw score in demo
	scoreText := fmt.Sprintf("Score: %d", g.score)
	ebitenutil.DebugPrintAt(screen, scoreText, 10, 30)
}

// drawGameUI draws the gameplay UI
func (g *Game) drawGameUI(screen *ebiten.Image) {
	// Draw score
	scoreText := fmt.Sprintf("Score: %d", g.score)
	ebitenutil.DebugPrintAt(screen, scoreText, 10, 30)

	// Draw game over message
	if g.isGameOver {
		gameOverText := "Game Over! Press 'R' to restart."
		ebitenutil.DebugPrintAt(screen, gameOverText, screenWidth/2-150, screenHeight/2)
	}

	// Draw controls help
	controlsText := "Press ESC to return to menu"
	ebitenutil.DebugPrintAt(screen, controlsText, screenWidth-200, screenHeight-20)
}
