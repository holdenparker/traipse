package algorithms

// Generates the area of a linear ring using the Shoelace formula
// Expects the ring to be at least 4 positions long, with the first
// and last position representing the same position.
func Shoelace(ring [][3]float64) float64 {
	prevX := ring[0][0]
	prevY := ring[0][1]

	var xSide float64 = 0
	var ySide float64 = 0
	for _, pos := range ring[1:] {
		xSide += prevX * pos[1]
		ySide += prevY * pos[0]

		prevX = pos[0]
		prevY = pos[1]
	}

	return (xSide - ySide) * 0.5
}
