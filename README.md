# Flappy Bird in Go

A complete implementation of the classic Flappy Bird game using Go and the Ebiten 2D game library. This version features a graphical UI with procedurally generated images and sounds.

## Description

Flappy Bird is a side-scrolling game where the player controls a bird, attempting to fly between columns of green pipes without hitting them. The game is designed to be simple yet challenging.

## Features

- Graphical UI using the Ebiten library with detailed procedurally generated images
- Main menu with auto-playing demo mode
- Sound effects for jumping, scoring, and collisions
- Background music in the menu
- Simple controls (press space to jump)
- Score tracking
- Game over and restart functionality
- No external assets required (images and sounds are generated at runtime)

## Requirements

- Go 1.16 or higher
- Ebiten v2 library (with audio support)

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/flappy-bird.git
   cd flappy-bird
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

## How to Run

To run the game, simply execute:

```
go run main.go
```

Or build and run the executable:

```
go build
./flappy-bird
```

## Game Controls

### Menu Controls
- **Up/Down Arrow Keys**: Navigate menu options
- **Space/Enter**: Select menu option
- **Escape**: Return to menu from gameplay

### Gameplay Controls
- **Space**: Make the bird jump/flap
- **R**: Restart the game after game over
- **Escape**: Return to main menu

## Game Rules

1. The game starts with a main menu featuring an auto-playing demo mode
2. Select "Start Game" to begin playing
3. The bird automatically moves forward and is affected by gravity
4. Press space to make the bird jump upward
5. Navigate through the gaps between pipes
6. Each successfully passed pipe earns one point
7. The game ends if the bird hits a pipe or the top/bottom of the screen
8. Press R to restart after game over or Escape to return to the main menu
9. Sound effects play when jumping, scoring, and colliding
10. Background music plays in the menu

## Game Mechanics

- The game gets progressively more challenging as you play
- Pipes are randomly generated with varying gap positions
- The bird's movement is affected by gravity, requiring precise timing of jumps
- All images are procedurally generated at runtime:
  - Detailed bird with body, eye, beak, and wing
  - Textured pipes with caps
  - Scrolling ground with grass
  - Sky background with distant mountains
- All sounds are procedurally generated at runtime:
  - Jump sound (upward frequency sweep)
  - Score sound (pleasant "ding")
  - Collision sound (crash effect)
  - Menu music (looping melody)

## License

This project is open source and available under the [MIT License](LICENSE).
