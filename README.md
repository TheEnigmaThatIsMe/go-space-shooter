# Go Space Shooter

**Go Space Shooter** is a 2D top-down arcade-style game built using the Go programming language and the Ebiten game library. Players control a spaceship to shoot asteroids, avoid collisions, and earn points as the game progresses.

## Features

- **Player-Controlled Spaceship**: Navigate a spaceship using keyboard controls.
- **Shooting Mechanism**: Shoot lasers at asteroids to destroy them.
- **Asteroid Generation**: Randomly spawning asteroids that move toward the player.
- **Score System**: Gain points for every asteroid destroyed.
- **Dynamic Scaling**: Sprites dynamically scale relative to the window size.

---

## Requirements

- **Go**: Version 1.18 or later
- **Ebiten**: Go game library for rendering and input handling

To install Ebiten:

```bash
go get -u github.com/hajimehoshi/ebiten/v2
```
---
## Installation
1. Clone the repository:
```bash
git clone https://github.com/TheEnigmaThatIsMe/go-space-shooter.git
cd go-space-shooter
```

2. Install dependencies:
```bash
go mod tidy
```
3. Run the game:
```bash
go run .
```
---
## Controls
- **Arrow Keys**: Move the spaceship
- **Spacebar**: Shoot lasers
---
## How to Play
1.	Launch the game by running the main.go file.
2.	Use arrow keys to move your spaceship left and right.
3.	Press the spacebar to shoot lasers at approaching asteroids.
4.	Destroy as many asteroids as possible to score points.
5.	Avoid collisions with asteroids to keep playing.

Enjoy the game!
---
## Future Improvements
- Add power-ups (e.g., faster lasers, shield boosts).
- Implement levels with increasing difficulty.
- Introduce more asteroid types and behaviors.
- Add a leaderboard to save high scores.
- Include sound effects and background music.