package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const NUMBER_BOIDS = 200
const TARGET_FRAMERATE = 60
const RADIUS_BOID = 5
const RADIUS_AVOIDANCE = 500
const RADIUS_VIEW = 8000

var debug = false

var screenHeigh = 768
var screenWidth = 1024
var speedMax = 0.0125 * float32(screenWidth)
var speedMin = 0.00375 * float32(screenWidth)

var avoidanceFactor = float32(1) / TARGET_FRAMERATE
var alignementFactor = float32(0.5) / TARGET_FRAMERATE
var cohesionFactor = float32(0.2) / TARGET_FRAMERATE
var turnFactor = float32(0.025) * speedMax

type Boid struct {
	Pos, Vel rl.Vector2
}

func NewBoid() *Boid {
	return &Boid{
		Pos: rl.Vector2{
			X: rand.Float32() * float32(screenWidth) / 1.25,
			Y: rand.Float32() * float32(screenHeigh) / 1.25,
		},
		Vel: rl.Vector2{
			X: rand.Float32() - 0.5,
			Y: rand.Float32() - 0.5,
		},
	}
}

func (b *Boid) Draw() {
	direction := rl.Vector2{X: b.Pos.X + b.Vel.X*2, Y: b.Pos.Y + b.Vel.Y*2}
	rl.DrawCircleV(b.Pos, RADIUS_BOID, rl.Blue)
	rl.DrawLineV(b.Pos, direction, rl.Blue)
}

// Baseline implementation, adapted from:
// https://vanhunteradams.com/Pico/Animal_Movement/Boids-algorithm.html
func BoidV0() {
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.SetConfigFlags(rl.FlagWindowResizable)

	rl.InitWindow(int32(screenWidth), int32(screenHeigh), "Boids")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	boids := [NUMBER_BOIDS]Boid{}
	for i := 0; i < NUMBER_BOIDS; i++ {
		boids[i] = *NewBoid()
	}

	for !rl.WindowShouldClose() {
		if rl.IsWindowResized() {
			screenWidth = rl.GetRenderWidth()
			screenHeigh = rl.GetRenderHeight()
		}
		// Calculate update
		for i := 0; i < len(boids); i++ {
			var closeDx, closeDy, xPosAvg, yPosAvg, xVelAvg, yVelAvg float32
			var neighbours int
			var avoid int

			for j := 0; j < len(boids); j++ {
				// Do not compare a boid to itself
				if i == j {
					continue
				}

				// Compute differences in x and y coordinates
				dx := boids[i].Pos.X - boids[j].Pos.X
				dy := boids[i].Pos.Y - boids[j].Pos.Y

				// Are both those differences less than the visual range?
				if math.Abs(float64(dx)) < RADIUS_VIEW &&
					math.Abs(float64(dy)) < RADIUS_VIEW {

					// If so, calculate the squared distance
					// distance := math.Sqrt(float64(dx*dx + dy*dy))
					distance := (dx*dx + dy*dy)

					// Is squared distance less than the protected range?
					if distance < RADIUS_AVOIDANCE {
						// If so, calculate difference in x/y-coordinates to nearfield boid
						closeDx += boids[i].Pos.X - boids[j].Pos.X
						closeDy += boids[i].Pos.Y - boids[j].Pos.Y
						avoid++
						// If not in protected range, is the boid in the visual range?
					} else if distance < RADIUS_VIEW {

						// Add other boid's x/y-coord and x/y vel to accumulator variables
						xPosAvg += boids[j].Pos.X
						yPosAvg += boids[j].Pos.Y
						xVelAvg += boids[j].Vel.X
						yVelAvg += boids[j].Vel.Y

						// Increment number of boids within visual range
						neighbours++
					}
				}
			}

			// If there were any boids in the visual range
			if neighbours > 0 {
				// Divide accumulator variables by number of boids in visual range
				xPosAvg = xPosAvg / float32(neighbours)
				yPosAvg = yPosAvg / float32(neighbours)
				xVelAvg = xVelAvg / float32(neighbours)
				yVelAvg = yVelAvg / float32(neighbours)

				// Add the centering/matching contributions to velocity
				boids[i].Vel.X += (xPosAvg - boids[i].Pos.X) * cohesionFactor
				boids[i].Vel.X += (xVelAvg - boids[i].Vel.X) * alignementFactor

				boids[i].Vel.Y += (yPosAvg - boids[i].Pos.Y) * cohesionFactor
				boids[i].Vel.Y += (yVelAvg - boids[i].Vel.Y) * alignementFactor
			}

			// Add the avoidance contribution to velocity
			boids[i].Vel.X += closeDx * avoidanceFactor
			boids[i].Vel.Y += closeDy * avoidanceFactor

			// If the boid is near an edge, make it turn by turnfactor
			if boids[i].Pos.X > 0.80*float32(screenWidth) {
				boids[i].Vel.X = boids[i].Vel.X - turnFactor
			}
			if boids[i].Pos.X < 0.20*float32(screenWidth) {
				boids[i].Vel.X = boids[i].Vel.X + turnFactor
			}
			if boids[i].Pos.Y > 0.80*float32(screenHeigh) {
				boids[i].Vel.Y = boids[i].Vel.Y - turnFactor
			}
			if boids[i].Pos.Y < 0.20*float32(screenHeigh) {
				boids[i].Vel.Y = boids[i].Vel.Y + turnFactor
			}

			// Calculate the boid's speed
			// Slow step! Lookup the "alpha max plus beta min" algorithm
			speed := float32(math.Sqrt(
				float64(boids[i].Vel.X*boids[i].Vel.X + boids[i].Vel.Y*boids[i].Vel.Y),
			))

			if speed < speedMin {
				boids[i].Vel.X = (boids[i].Vel.X / speed) * speedMin
				boids[i].Vel.Y = (boids[i].Vel.Y / speed) * speedMin
			}
			if speed > speedMax {
				boids[i].Vel.X = (boids[i].Vel.X / speed) * speedMax
				boids[i].Vel.Y = (boids[i].Vel.Y / speed) * speedMax
			}

			// Update boid's position
			boids[i].Pos.X += boids[i].Vel.X
			boids[i].Pos.Y += boids[i].Vel.Y
		}

		// Draw update
		rl.BeginDrawing()

		rl.ClearBackground(color.RGBA{22, 23, 31, 255})
		for i := 0; i < len(boids); i++ {
			boids[i].Draw()
		}

		if debug {
			rl.DrawText(fmt.Sprint("FPS: ", rl.GetFPS()), 20, 20, 20, rl.LightGray)

		}
		rl.EndDrawing()
	}
}

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "--debug" {
		debug = true
	}
	BoidV0()
}
