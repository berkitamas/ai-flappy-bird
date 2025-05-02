package main

import (
	"math"
	"math/rand"
)

// updateGameplay updates the gameplay elements (bird, pipes, collisions)
func (g *Game) updateGameplay(demoMode bool) {
	// Apply gravity to bird
	g.bird.velocity += gravity
	g.bird.y += g.bird.velocity

	// Update wing flapping animation
	g.bird.wingAngle = math.Sin(float64(g.tick)*0.3) * 0.3

	// Generate new pipes
	if g.tick%pipeSpacing == 0 {
		gapStart := rand.Intn(screenHeight-pipeGap-100) + 50
		g.pipes = append(g.pipes, Pipe{
			x:        screenWidth,
			gapStart: gapStart,
		})
	}

	// Update pipes and check for score
	for i := 0; i < len(g.pipes); i++ {
		// Move pipe
		g.pipes[i].x -= pipeSpeed

		// Check if bird passed the pipe
		birdCenterX := g.bird.x + float64(g.bird.width)/2
		if g.pipes[i].x+pipeWidth <= int(birdCenterX) && g.pipes[i].x+pipeWidth > int(birdCenterX)-pipeSpeed {
			// Bird passed a pipe, increment score
			g.score++
			g.scoreSound.Rewind()
			g.scoreSound.Play()
		}

		// Remove pipes that have moved off-screen
		if g.pipes[i].x+pipeWidth < 0 {
			g.pipes = append(g.pipes[:i], g.pipes[i+1:]...)
			i--
		}
	}

	// Check for collisions
	if demoMode {
		// In demo mode, check for collisions but don't end the game
		if g.checkCollisions() {
			// If collision detected in demo mode, restart after a short delay
			if g.tick%60 == 0 { // Wait about 1 second (60 frames)
				g.Reset()
			}
			g.isGameOver = true // Stop background movement
		}

		// Control bird in demo mode
		g.controlBirdAI()
	} else {
		// In play mode, check for collisions and end the game if needed
		if g.checkCollisions() && !g.isGameOver {
			g.isGameOver = true
			g.collisionSound.Rewind()
			g.collisionSound.Play()
		}
	}
}

// checkCollisions checks if the bird has collided with pipes or boundaries
func (g *Game) checkCollisions() bool {
	// Check if bird hit the ground or ceiling
	if g.bird.y <= 0 || g.bird.y+float64(g.bird.height) >= float64(screenHeight-30) {
		return true
	}

	// Check if bird hit any pipes
	birdCenterX := g.bird.x + float64(g.bird.width)/2
	birdCenterY := g.bird.y + float64(g.bird.height)/2
	collisionRadius := float64(g.bird.width) * 0.6 // Use a slightly larger collision radius

	for _, pipe := range g.pipes {
		pipeX := float64(pipe.x)
		pipeRightX := pipeX + float64(pipeWidth)

		// Only check pipes that are near the bird horizontally
		if pipeRightX < g.bird.x-collisionRadius || pipeX > g.bird.x+float64(g.bird.width)+collisionRadius {
			continue
		}

		// Check collision with top pipe
		topPipeBottom := float64(pipe.gapStart)
		if birdCenterY-collisionRadius < topPipeBottom && birdCenterX+collisionRadius > pipeX && birdCenterX-collisionRadius < pipeRightX {
			return true
		}

		// Check collision with bottom pipe
		bottomPipeTop := float64(pipe.gapStart + pipeGap)
		if birdCenterY+collisionRadius > bottomPipeTop && birdCenterX+collisionRadius > pipeX && birdCenterX-collisionRadius < pipeRightX {
			return true
		}
	}

	return false
}

// wouldCollide checks if the bird would collide with pipes or boundaries without setting game over
func (g *Game) wouldCollide() bool {
	// Check if bird would hit the ground or ceiling
	if g.bird.y <= 0 || g.bird.y+float64(g.bird.height) >= float64(screenHeight-30) {
		return true
	}

	// Check if bird would hit any pipes
	birdCenterX := g.bird.x + float64(g.bird.width)/2
	birdCenterY := g.bird.y + float64(g.bird.height)/2
	collisionRadius := float64(g.bird.width) * 0.6 // Use a slightly larger collision radius

	for _, pipe := range g.pipes {
		pipeX := float64(pipe.x)
		pipeRightX := pipeX + float64(pipeWidth)

		// Only check pipes that are near the bird horizontally
		if pipeRightX < g.bird.x-collisionRadius || pipeX > g.bird.x+float64(g.bird.width)+collisionRadius {
			continue
		}

		// Check collision with top pipe
		topPipeBottom := float64(pipe.gapStart)
		if birdCenterY-collisionRadius < topPipeBottom && birdCenterX+collisionRadius > pipeX && birdCenterX-collisionRadius < pipeRightX {
			return true
		}

		// Check collision with bottom pipe
		bottomPipeTop := float64(pipe.gapStart + pipeGap)
		if birdCenterY+collisionRadius > bottomPipeTop && birdCenterX+collisionRadius > pipeX && birdCenterX-collisionRadius < pipeRightX {
			return true
		}
	}

	return false
}

// wouldCollideAfterJump simulates a jump and checks if it would lead to a collision
func (g *Game) wouldCollideAfterJump() bool {
	// Save current state
	originalY := g.bird.y
	originalVelocity := g.bird.velocity

	// Simulate jump
	g.bird.velocity = jumpForce

	// Simulate several frames ahead
	for i := 0; i < 5; i++ {
		g.bird.velocity += gravity
		g.bird.y += g.bird.velocity

		if g.wouldCollide() {
			// Restore original state
			g.bird.y = originalY
			g.bird.velocity = originalVelocity
			return true
		}
	}

	// Restore original state
	g.bird.y = originalY
	g.bird.velocity = originalVelocity
	return false
}

// wouldCollideWithoutJump simulates not jumping and checks if it would lead to a collision
func (g *Game) wouldCollideWithoutJump() bool {
	// Save current state
	originalY := g.bird.y
	originalVelocity := g.bird.velocity

	// Simulate several frames ahead without jumping
	for i := 0; i < 5; i++ {
		g.bird.velocity += gravity
		g.bird.y += g.bird.velocity

		if g.wouldCollide() {
			// Restore original state
			g.bird.y = originalY
			g.bird.velocity = originalVelocity
			return true
		}
	}

	// Restore original state
	g.bird.y = originalY
	g.bird.velocity = originalVelocity
	return false
}
