package main

import (
	"testing"
)

// TestControlBirdAI tests the AI control logic
func TestControlBirdAI(t *testing.T) {
	// Create a game instance with a bird and pipes for testing
	g := &Game{
		bird: Bird{
			x:        float64(screenWidth) / 4,
			y:        float64(screenHeight) / 2,
			width:    birdSize,
			height:   birdSize,
			velocity: 0,
		},
		pipes: []Pipe{},
	}

	// Test 1: No pipes - bird should maintain a middle position
	g.controlBirdAI()
	// We can't make specific assertions about the bird's velocity since it depends on the bird's position
	// But we can verify that the function doesn't panic

	// Test 2: Bird too low with no pipes - should jump
	g.bird.y = float64(screenHeight) * 0.8
	g.bird.velocity = 0
	initialVelocity := g.bird.velocity
	g.controlBirdAI()
	if g.bird.velocity >= initialVelocity {
		t.Error("Expected bird to jump when too low with no pipes")
	}

	// Test 3: Bird too high with no pipes and moving up - should not jump
	g.bird.y = float64(screenHeight) * 0.2
	g.bird.velocity = -2
	initialVelocity = g.bird.velocity
	g.controlBirdAI()
	if g.bird.velocity < initialVelocity {
		t.Error("Expected bird not to jump when too high and moving up with no pipes")
	}

	// Test 4: Bird moving down too fast with no pipes - should jump
	g.bird.y = float64(screenHeight) / 2
	g.bird.velocity = 5
	g.controlBirdAI()
	if g.bird.velocity >= 5 {
		t.Error("Expected bird to jump when moving down too fast with no pipes")
	}

	// Test 5: With a pipe far ahead - bird should maintain a good position
	g.pipes = []Pipe{{
		x:        screenWidth,
		gapStart: int(float64(screenHeight)/2) - pipeGap/2,
	}}
	g.bird.y = float64(screenHeight) / 2
	g.bird.velocity = 0
	g.controlBirdAI()
	// Again, we can't make specific assertions about the bird's velocity
	// But we can verify that the function doesn't panic

	// Test 6: With a pipe close ahead and bird below gap - should jump
	g.pipes = []Pipe{{
		x:        int(g.bird.x) + 100,
		gapStart: int(g.bird.y) - pipeGap - 20, // Gap is above the bird
	}}
	g.bird.velocity = 0
	g.controlBirdAI()
	if g.bird.velocity >= 0 {
		t.Error("Expected bird to jump when below gap with pipe close ahead")
	}

	// Test 7: With a pipe close ahead and bird above gap - should not jump
	g.pipes = []Pipe{{
		x:        int(g.bird.x) + 100,
		gapStart: int(g.bird.y) + 20, // Gap is below the bird
	}}
	g.bird.velocity = 0
	g.controlBirdAI()
	if g.bird.velocity < 0 {
		t.Error("Expected bird not to jump when above gap with pipe close ahead")
	}

	// Test 8: With two pipes ahead - should consider both
	g.pipes = []Pipe{
		{
			x:        int(g.bird.x) + 150,
			gapStart: int(g.bird.y) - pipeGap/2, // First gap centered on bird
		},
		{
			x:        int(g.bird.x) + 300,
			gapStart: int(g.bird.y) - pipeGap - 50, // Second gap above bird
		},
	}
	g.bird.velocity = 0
	g.controlBirdAI()
	// We can't make specific assertions about the bird's velocity
	// But we can verify that the function doesn't panic
}

// TestFindNextPipe tests the logic for finding the next pipe
func TestFindNextPipe(t *testing.T) {
	// Create a game instance with a bird and pipes for testing
	g := &Game{
		bird: Bird{
			x:      float64(screenWidth) / 4,
			y:      float64(screenHeight) / 2,
			width:  birdSize,
			height: birdSize,
		},
		pipes: []Pipe{},
	}

	// Test the logic for finding the next pipe by examining the controlBirdAI function
	// This is an indirect test since the function doesn't return the next pipe

	// Test 1: No pipes
	g.controlBirdAI()
	// Verify the function doesn't panic with no pipes

	// Test 2: One pipe ahead
	g.pipes = []Pipe{{
		x:        int(g.bird.x) + 100,
		gapStart: int(g.bird.y) - pipeGap/2,
	}}
	g.controlBirdAI()
	// Verify the function doesn't panic with one pipe

	// Test 3: One pipe behind
	g.pipes = []Pipe{{
		x:        int(g.bird.x) - 100,
		gapStart: int(g.bird.y) - pipeGap/2,
	}}
	g.controlBirdAI()
	// Verify the function doesn't panic with one pipe behind

	// Test 4: Multiple pipes ahead
	g.pipes = []Pipe{
		{
			x:        int(g.bird.x) + 200,
			gapStart: int(g.bird.y) - pipeGap/2,
		},
		{
			x:        int(g.bird.x) + 100,
			gapStart: int(g.bird.y) - pipeGap/2,
		},
	}
	g.controlBirdAI()
	// Verify the function doesn't panic with multiple pipes
}

// TestPredictBirdPosition tests the logic for predicting the bird's position
func TestPredictBirdPosition(t *testing.T) {
	// Create a game instance with a bird for testing
	g := &Game{
		bird: Bird{
			x:        float64(screenWidth) / 4,
			y:        float64(screenHeight) / 2,
			width:    birdSize,
			height:   birdSize,
			velocity: 0,
		},
	}

	// Test the prediction logic by examining the controlBirdAI function
	// This is an indirect test since the function doesn't return the prediction

	// Test 1: Bird with zero velocity
	initialY := g.bird.y
	g.bird.velocity = 0
	g.controlBirdAI()
	// Verify the function doesn't panic

	// Test 2: Bird with positive velocity (falling)
	g.bird.y = initialY
	g.bird.velocity = 2
	g.controlBirdAI()
	// Verify the function doesn't panic

	// Test 3: Bird with negative velocity (rising)
	g.bird.y = initialY
	g.bird.velocity = -2
	g.controlBirdAI()
	// Verify the function doesn't panic
}
