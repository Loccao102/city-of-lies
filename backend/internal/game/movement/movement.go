package movement

import (
	"math"
)

type AgentTransit struct {
	AgentID          string
	FromLocationID   string
	ToLocationID     string
	DepartGameSecond int64
	ArriveGameSecond int64
}

// CalculateTransitDuration estimates travel time in seconds between two points.
func CalculateTransitDuration(fromX, fromZ, toX, toZ float64) int64 {
	dx := toX - fromX
	dz := toZ - fromZ
	dist := math.Sqrt(dx*dx + dz*dz)
	// Base walking speed: ~1.0 units per second, clamped between 15s and 45s
	duration := int64(dist)
	if duration < 15 {
		duration = 15
	}
	if duration > 45 {
		duration = 45
	}
	return duration
}
