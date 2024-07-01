package geo

import (
	"math"
)

// DegreesToDirectionCcw converts degrees to a direction
// vector, where 0 = right and it continues clockwise, so 90 is down.
func DegreesToDirectionCw(degrees float64) PtF {
	theta := math.Pi * degrees / 180.0 // Convert degrees to radians
	return Pt(math.Cos(theta), math.Sin(theta))
}
