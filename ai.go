package main

import (
	"math"
)

// controlBirdAI controls the bird in demo mode
func (g *Game) controlBirdAI() {
	// Find the next pipe to navigate through
	var nextPipe *Pipe
	var secondPipe *Pipe
	birdCenterX := g.bird.x + float64(g.bird.width)/2
	birdCenterY := g.bird.y + float64(g.bird.height)/2

	// Find the next pipe and the second pipe (if any)
	for i, pipe := range g.pipes {
		pipeX := float64(pipe.x)

		// If pipe is ahead of the bird, it's a candidate for next pipe
		if pipeX > birdCenterX {
			if nextPipe == nil {
				nextPipe = &g.pipes[i]
			} else if pipeX < float64(nextPipe.x) {
				secondPipe = nextPipe
				nextPipe = &g.pipes[i]
			} else if secondPipe == nil || pipeX < float64(secondPipe.x) {
				secondPipe = &g.pipes[i]
			}
		}
	}

	// If there's a pipe to navigate through
	if nextPipe != nil {
		pipeX := float64(nextPipe.x)
		gapCenter := float64(nextPipe.gapStart) + float64(pipeGap)/2
		distanceToPipe := pipeX - birdCenterX

		// Calculate how many frames it will take to reach the pipe
		framesToPipe := distanceToPipe / pipeSpeed

		// Predict where the bird will be when it reaches the pipe
		predictedY := birdCenterY
		predictedVelocity := g.bird.velocity

		for i := 0; i < int(framesToPipe) && i < 30; i++ {
			predictedVelocity += gravity
			predictedY += predictedVelocity
		}

		// Different strategies based on distance to pipe
		if distanceToPipe < 120 {
			// Very close to pipe - precise positioning
			verticalDistToGap := predictedY - gapCenter

			// Look ahead to the second pipe if it's close enough
			if secondPipe != nil {
				secondPipeX := float64(secondPipe.x)
				secondGapCenter := float64(secondPipe.gapStart) + float64(pipeGap)/2
				distanceToSecondPipe := secondPipeX - birdCenterX

				if distanceToSecondPipe < 250 {
					// Start adjusting for the second pipe if the gap is significantly higher or lower
					if math.Abs(secondGapCenter-gapCenter) > 40 {
						// Blend between current gap and next gap based on distance
						blendFactor := 1.0 - (distanceToPipe / distanceToSecondPipe)
						gapCenter = gapCenter*(1-blendFactor) + secondGapCenter*blendFactor
						verticalDistToGap = predictedY - gapCenter
					}
				}
			}

			// Fine-tuning based on current position and velocity
			if verticalDistToGap < -20 && g.bird.velocity > 0 {
				// Above gap and moving down - let gravity work
				return
			} else if verticalDistToGap > 20 && g.bird.velocity < 0 {
				// Below gap and moving up - let momentum carry
				return
			} else if verticalDistToGap > 20 {
				// Below gap and not moving up enough - jump
				g.bird.velocity = jumpForce
				return
			}

			// Additional check for velocity - allow more speed variation
			if g.bird.velocity > 4 {
				// Moving down too fast
				g.bird.velocity = jumpForce
				return
			} else if g.bird.velocity < -4 {
				// Moving up too fast
				return
			}
		} else if distanceToPipe < 250 {
			// Medium distance - start aligning with gap
			// Use predicted position for better decision making
			verticalDistToPredicted := predictedY - gapCenter

			if verticalDistToPredicted < -35 {
				// Predicted to be too high, let gravity work
				return
			} else if verticalDistToPredicted > 35 {
				// Predicted to be too low, jump
				g.bird.velocity = jumpForce
				return
			} else if math.Abs(verticalDistToPredicted) < 25 && g.bird.velocity > 2.0 {
				// Near center but moving down too fast
				g.bird.velocity = jumpForce
				return
			}
		} else {
			// Far from pipe - maintain a good position
			// Use a more central target position
			targetY := float64(screenHeight) / 2

			// Blend between middle of screen and gap center based on distance
			blendFactor := math.Min(1, distanceToPipe/300)
			targetY = gapCenter*(1-blendFactor) + targetY*blendFactor

			if birdCenterY > targetY+40 {
				// Too low, jump
				g.bird.velocity = jumpForce
				return
			} else if birdCenterY < targetY-60 && g.bird.velocity < 0 {
				// Too high and still going up
				return
			}
		}
	} else {
		// If no pipe is visible, maintain a middle position
		birdCenterY := g.bird.y + float64(g.bird.height)/2
		middleY := float64(screenHeight) / 2.2 // Slightly lower than center to prepare for pipes

		// Keep bird in the middle third of the screen with more flexibility
		if birdCenterY > middleY+40 {
			// Too low, jump
			g.bird.velocity = jumpForce
		} else if birdCenterY < middleY-60 && g.bird.velocity < 0 {
			// Too high and still going up
			return
		}

		// Don't let velocity get too extreme, but allow more variation
		if g.bird.velocity > 4 {
			// Moving down too fast
			g.bird.velocity = jumpForce
		} else if g.bird.velocity < -4 {
			// Moving up too fast
			return
		}
	}
}
