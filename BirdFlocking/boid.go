package main

import (
	"math"
	"math/rand"
	"time"
)

type Boid struct {
	position Vector2D
	velocity Vector2D
	id       int
}

func (b *Boid) calcAcceleration() Vector2D {
	upper, lower := b.position.AddV(viewRadius), b.position.AddV(-viewRadius)
	avgPosition, avgVelocity, avgSeperation := Vector2D{0, 0}, Vector2D{0, 0}, Vector2D{0, 0}
	count := 0.0
	rWlock.RLock()
	for i := math.Max(lower.x, 0); i <= math.Min(upper.x, screenWidth); i++ {
		for j := math.Max(lower.y, 0); j <= math.Min(upper.y, screenHeight); j++ {
			if otherBoidId := boidMap[int(i)][int(j)]; otherBoidId != -1 && otherBoidId != b.id {
				if dist := boids[otherBoidId].position.Distance(b.position); dist < viewRadius {
					count++
					avgVelocity = avgVelocity.Add(boids[otherBoidId].velocity)
					avgPosition = avgPosition.Add(boids[otherBoidId].position)
					avgSeperation = avgSeperation.Add(b.position.Subtract(boids[otherBoidId].position).DivisionV(dist))
					// separation = separation.Add(b.position.Subtract(boids[otherBoidId].position).DivisionV(dist))
				}
			}
		}
	}
	rWlock.RUnlock()

	accel := Vector2D{b.borderBounce(b.position.x, screenWidth), b.borderBounce(b.position.y, screenHeight)}
	if count > 0 {
		avgPosition, avgVelocity = avgPosition.DivisionV(count), avgVelocity.DivisionV(count)
		accelAlignment := avgVelocity.Subtract(b.velocity).MultiplyV(adjRate)
		accelCohesion := avgPosition.Subtract(b.position).MultiplyV(adjRate)
		accelSeperation := avgSeperation.MultiplyV(adjRate)
		accel = accel.Add(accelAlignment).Add(accelCohesion).Add(accelSeperation)
	}
	return accel
}

func (b *Boid) borderBounce(pos, maxBorderPos float64) float64 {
	if pos < viewRadius {
		return 1 / pos
	} else if pos > maxBorderPos-viewRadius {
		return 1 / (pos - maxBorderPos)
	}
	return 0
}

// func (b *Boid) calcAcceleration() Vector2D {
// 	upper, lower := b.position.AddV(viewRadius), b.position.AddV(-viewRadius)
// 	avgPosition, avgVelocity := Vector2D{0, 0}, Vector2D{0, 0}
// 	count := 0.0
// 	rWlock.Lock()
// 	for i := math.Max(lower.x, 0); i <= math.Min(upper.x, screenWidth); i++ {
// 		for j := math.Max(lower.y, 0); j <= math.Min(upper.y, screenHeight); j++ {
// 			if otherBoidId := boidMap[int(i)][int(j)]; otherBoidId != -1 && otherBoidId != b.id {
// 				if dist := boids[otherBoidId].position.Distance(b.position); dist < viewRadius {
// 					count++
// 					avgVelocity = avgVelocity.Add(boids[otherBoidId].velocity)
// 					avgPosition = avgVelocity.Add(boids[otherBoidId].position)

// 				}
// 			}
// 		}
// 	}
// 	rWlock.Unlock()

// 	accel := Vector2D{0, 0}
// 	if count > 0 {
// 		avgVelocity = avgVelocity.DivisionV(count)
// 		avgPosition = avgPosition.DivisionV(count)
// 		accelAlignment := avgVelocity.Subtract(b.velocity).MultiplyV(adjRate)
// 		accelCohesion := avgVelocity.Subtract(b.position).MultiplyV(adjRate)
// 		accel = accel.Add(accelAlignment).Add(accelCohesion)
// 	}

// 	return accel
// }

func (b *Boid) moveOne() {
	acceleration := b.calcAcceleration()
	rWlock.Lock()
	b.velocity = b.velocity.Add(acceleration).limit(-1, 1)
	boidMap[int(b.position.x)][int(b.position.y)] = -1
	b.position = b.position.Add(b.velocity)
	boidMap[int(b.position.x)][int(b.position.y)] = b.id

	rWlock.Unlock()
}

func (b *Boid) start() {
	for {
		b.moveOne()
		time.Sleep(5 * time.Millisecond)
	}
}

func createBoid(bid int) {
	b := Boid{
		position: Vector2D{x: rand.Float64() * screenWidth, y: rand.Float64()},
		velocity: Vector2D{x: (rand.Float64() * 2) - 1, y: (rand.Float64() * 2) - 1},
		id:       bid,
	}
	boids[bid] = &b
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	go b.start()
}
