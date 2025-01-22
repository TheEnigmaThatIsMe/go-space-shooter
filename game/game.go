package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"log"
	"math/rand"
	"time"
)

const (
	screenWidth       = 800
	screenHeight      = 600
	playerSpeed       = 7
	bulletSpeed       = 7
	asteroidSpeed     = 2
	spriteFrameWidth  = 170 // Example: Adjust based on sprite sheet
	spriteFrameHeight = 158
	spriteCols        = 5   // Number of frames per row
	bulletCooldown    = 0.2 // Cooldown period in seconds
)

var playerSpriteSheet *ebiten.Image
var asteroidSpriteSheet *ebiten.Image
var bulletSpriteSheet *ebiten.Image

// Game holds the state of the game.
type Game struct {
	playerImage    *ebiten.Image
	playerX        float64
	playerY        float64
	bullets        []Bullet
	asteroids      []Asteroid
	score          int
	lastBulletTime float64
}

// Bullet represents a bullet shot by the player.
type Bullet struct {
	bulletImage *ebiten.Image
	x, y        float64
}

// Asteroid represents an incoming asteroid.
type Asteroid struct {
	asteroidImage *ebiten.Image
	x, y          float64
}

func initSpriteSheets() {
	var err error
	playerSpriteSheet, _, err = ebitenutil.NewImageFromFile("assets/spaceship.png")
	if err != nil {
		log.Fatalf("Failed to load player sprite sheet: %v", err)
	}

	asteroidSpriteSheet, _, err = ebitenutil.NewImageFromFile("assets/asteroids_sprite.png")
	if err != nil {
		log.Fatalf("Failed to load asteroid sprite sheet: %v", err)
	}

	bulletSpriteSheet, _, err = ebitenutil.NewImageFromFile("assets/laser.png")
	if err != nil {
		log.Fatalf("Failed to load laser sprite sheet: %v", err)
	}
}

// NewGame initializes and returns a new game.
func NewGame() *Game {
	initSpriteSheets()
	playerImg := getSpriteImage(playerSpriteSheet, 0) // 0 is the frame index for idle

	return &Game{
		playerImage: playerImg,
		playerX:     screenWidth / 2,
		playerY:     screenHeight - 70,
		bullets:     []Bullet{},
		asteroids:   []Asteroid{},
	}
}

func getSpriteImage(spriteSheet *ebiten.Image, frameIndex int) *ebiten.Image {
	frameX := (frameIndex % spriteCols) * spriteFrameWidth
	frameY := (frameIndex / spriteCols) * spriteFrameHeight
	frameRect := image.Rect(frameX, frameY, frameX+spriteFrameWidth, frameY+spriteFrameHeight)

	// Extract the frame as a subimage
	return spriteSheet.SubImage(frameRect).(*ebiten.Image)
}

// Update updates the game state.
func (g *Game) Update() error {
	currentTime := float64(time.Now().UnixMilli()) / 1000

	// Move player
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && g.playerX > 0 {
		g.playerX -= playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) && g.playerX < screenWidth-50 {
		g.playerX += playerSpeed
	}

	// Shoot bullets with cooldown
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if currentTime-g.lastBulletTime >= bulletCooldown {
			g.bullets = append(g.bullets, Bullet{bulletImage: bulletSpriteSheet, x: g.playerX + 20, y: g.playerY})
			g.lastBulletTime = currentTime
		}
	}

	// Update bullets
	for i := 0; i < len(g.bullets); i++ {
		g.bullets[i].y -= bulletSpeed
		if g.bullets[i].y < 0 {
			g.bullets = append(g.bullets[:i], g.bullets[i+1:]...)
			i--
		}
	}

	// Generate asteroids randomly
	if rand.Intn(100) < 2 { // Adjust frequency
		g.asteroids = append(g.asteroids, Asteroid{asteroidImage: getSpriteImage(asteroidSpriteSheet, rand.Intn(5)), x: float64(rand.Intn(screenWidth)), y: 0})
	}

	// Update asteroids
	for i := 0; i < len(g.asteroids); i++ {
		g.asteroids[i].y += asteroidSpeed
		if g.asteroids[i].y > screenHeight {
			g.asteroids = append(g.asteroids[:i], g.asteroids[i+1:]...)
			i--
		}
	}

	// Check for collisions
	for i := 0; i < len(g.asteroids); i++ {
		asteroidWidth := float64(g.asteroids[i].asteroidImage.Bounds().Dx()) * (float64(screenWidth) * 0.1 / float64(g.asteroids[i].asteroidImage.Bounds().Dx()))
		asteroidHeight := float64(g.asteroids[i].asteroidImage.Bounds().Dy()) * (float64(screenHeight) * 0.1 / float64(g.asteroids[i].asteroidImage.Bounds().Dy()))

		for j := 0; j < len(g.bullets); j++ {
			bulletWidth := float64(g.bullets[j].bulletImage.Bounds().Dx()) * (float64(screenWidth) * 0.15 / float64(g.bullets[j].bulletImage.Bounds().Dx()))
			bulletHeight := float64(g.bullets[j].bulletImage.Bounds().Dy()) * (float64(screenHeight) * 0.15 / float64(g.bullets[j].bulletImage.Bounds().Dy()))

			if isColliding(g.asteroids[i].x, g.asteroids[i].y, asteroidWidth, asteroidHeight, g.bullets[j].x, g.bullets[j].y, bulletWidth, bulletHeight) {
				g.asteroids = append(g.asteroids[:i], g.asteroids[i+1:]...)
				g.bullets = append(g.bullets[:j], g.bullets[j+1:]...)
				g.score++
				i--
				break
			}
		}
	}

	return nil
}

// Draw renders the game screen.
func (g *Game) Draw(screen *ebiten.Image) {

	// Example: Draw the spaceship idle frame at position (100, 100)
	//drawSprite(screen, 100, 100, 0) // 0 is the frame index for idle

	DrawPlayer(g, screen)

	// Draw bullets
	for _, b := range g.bullets {
		DrawBullet(b, screen)
	}

	// Draw asteroids
	for _, a := range g.asteroids {
		DrawAsteroid(a, screen)
		//vector.DrawFilledRect(screen, float32(a.x), float32(a.y), 40, 40, color.Gray{Y: 0x80}, true)
	}

	// Display score
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d", g.score))
}

func DrawPlayer(g *Game, screen *ebiten.Image) {
	// Get the current size of the spaceship image
	imageWidth := g.playerImage.Bounds().Dx()
	imageHeight := g.playerImage.Bounds().Dy()

	// Desired size relative to the window
	desiredWidth := float64(screenWidth) * 0.15 // 15% of the window width
	desiredHeight := float64(screenHeight) * 0.2

	// Calculate scaling factors
	scaleX := desiredWidth / float64(imageWidth)
	scaleY := desiredHeight / float64(imageHeight)

	// Apply scaling and translation
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(scaleX, scaleY)
	options.GeoM.Translate(g.playerX, g.playerY)

	// Draw the spaceship
	screen.DrawImage(g.playerImage, options)
}

func DrawBullet(b Bullet, screen *ebiten.Image) {
	// Get the current size of the bullet image
	imageWidth := b.bulletImage.Bounds().Dx()
	imageHeight := b.bulletImage.Bounds().Dy()

	// Desired size relative to the window
	desiredWidth := float64(screenWidth) * 0.15 // 15% of the window width
	desiredHeight := float64(screenHeight) * 0.15

	// Calculate scaling factors
	scaleX := desiredWidth / float64(imageWidth)
	scaleY := desiredHeight / float64(imageHeight)

	// Apply scaling and translating
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(scaleX, scaleY)
	options.GeoM.Translate(b.x, b.y)

	// Draw the bullet
	screen.DrawImage(b.bulletImage, options)
}

func DrawAsteroid(a Asteroid, screen *ebiten.Image) {
	// Get the current size of the bullet image
	imageWidth := a.asteroidImage.Bounds().Dx()
	imageHeight := a.asteroidImage.Bounds().Dy()

	// Desired size relative to the window
	desiredWidth := float64(screenWidth) * 0.1
	desiredHeight := float64(screenHeight) * 0.1

	// Calculate scaling factors
	scaleX := desiredWidth / float64(imageWidth)
	scaleY := desiredHeight / float64(imageHeight)

	// Apply scaling and translating
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(scaleX, scaleY)
	options.GeoM.Translate(a.x, a.y)

	// Draw the bullet
	screen.DrawImage(a.asteroidImage, options)
}

// Layout sets the screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// Helper function to detect collision
func isColliding(ax, ay, aw, ah, bx, by, bw, bh float64) bool {
	return ax < bx+bw && ax+aw > bx && ay < by+bh && ay+ah > by
}
