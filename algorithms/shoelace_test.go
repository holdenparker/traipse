package algorithms

import "testing"

func TestShoelace(t *testing.T) {
	expectedArea := 1.0
	linearRing := [][3]float64{
		{100.0, 0.0, 0},
		{101.0, 0.0, 0},
		{101.0, 1.0, 0},
		{100.0, 1.0, 0},
		{100.0, 0.0, 0},
	}
	actualArea := Shoelace(linearRing)

	if expectedArea != actualArea {
		t.Fatalf("Shoelace did not generate the expected area!\n Expected: %v\n Actual: %v\n", expectedArea, actualArea)
	}

	expectedArea = 26.5
	linearRing = [][3]float64{
		{2, 1, 0},
		{8, 3, 0},
		{6, 8, 0},
		{3, 4, 0},
		{-1, 6, 0},
		{2, 1, 0},
	}
	actualArea = Shoelace(linearRing)

	if expectedArea != actualArea {
		t.Fatalf("Shoelace did not generate the expected area!\n Expected: %v\n Actual: %v\n", expectedArea, actualArea)
	}

	// This is a reverse-wound ring, so the area should be negative
	expectedArea = -26.5
	linearRing = [][3]float64{
		{2, 1, 0},
		{-1, 6, 0},
		{3, 4, 0},
		{6, 8, 0},
		{8, 3, 0},
		{2, 1, 0},
	}
	actualArea = Shoelace(linearRing)

	if expectedArea != actualArea {
		t.Fatalf("Shoelace did not generate the expected area!\n Expected: %v\n Actual: %v\n", expectedArea, actualArea)
	}
}
