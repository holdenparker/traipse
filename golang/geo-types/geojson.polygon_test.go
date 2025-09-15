package geotypes

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func TestPolygonGeometry(t *testing.T) {
	geoJson := []byte(`{
		"type": "Point",
		"coordinates": [[[50.0, 40, 0], [90.0, 90.0, 0]]]
	}`)

	unmarshalResult := &PolygonGeometry{}
	err := json.Unmarshal(geoJson, unmarshalResult)

	if err == nil {
		t.Fatal("We should error when parsing a Point as a Polygon!")
	}
	if !errors.Is(err, UnmarshallingTypeMismatch) {
		t.Fatalf("We should return an UnmarshallingTypeMismatch!\nActual: %v\n", err)
	}
	if !errors.Is(err, PolygonGeometryUnmarshallingError) {
		t.Fatalf("We should return a PolygonGeometryUnmarshallingError!\nActual: %v\n", err)
	}

	geoJson = []byte(`{
		"type": "Polygon",
		"coordinates": [[[50.0, 40, 0], [90.0, 90.0, 0], [40.0, 40.0, 0], [50.0, 40, 0]]]
	}`)
	polygon := &PolygonGeometry{
		Coordinates: PolygonCoords{
			LineStringCoords{
				Position{50.0, 40.0, 0},
				Position{90.0, 90.0, 0},
				Position{40, 40, 0},
				Position{50, 40, 0},
			},
		},
	}

	unmarshalResult = &PolygonGeometry{}
	err = json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling Polygon:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, polygon) {
		t.Fatalf("Unmarshalled point geometry should match!\n Expected: %v\n Actual:   %v", polygon, unmarshalResult)
	}

	expected := `{"coordinates":[[[50,40,0],[90,90,0],[40,40,0],[50,40,0]]],"type":"Polygon"}`
	marshalResult, err := polygon.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling Feature:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled polygon feature should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}

	if !polygon.IsValid() {
		t.Fatalf("The polygon should have a valid coordinate definition!\nActual: %v\n", polygon.Coordinates)
	}

	polygon.Coordinates = PolygonCoords{
		LineStringCoords{
			Position{50.0, 40.0, 0},
			Position{90.0, 90.0, 0},
			Position{40, 40, 0},
		},
	}

	if polygon.IsValid() {
		t.Fatalf("The polygon should not be valid with invalid coordinates!\nActual: %v\n", polygon.Coordinates)
	}
}

func TestPolygonFeatures(t *testing.T) {
	geoJson := []byte(`{
		"type": "Feature",
		"id": "abc",
		"bbox": [100, 100, 0, 70, 70, 0],
		"geometry": {
				"type": "Polygon",
				"coordinates": [[[90.0, 90.0, 0], [100.0, 100.0, 0], [70.0, 70.0, 0]]]
		},
		"properties": {
				"FieldA": "blue",
				"Test": "blargh"
		}
	}`)
	polygon := &Feature{
		ID:         "abc",
		Bbox:       Bbox{100, 100, 0, 70, 70, 0},
		Properties: map[string]string{"FieldA": "blue", "Test": "blargh"},
		Geometry: &PolygonGeometry{
			Coordinates: PolygonCoords{
				LineStringCoords{
					Position{90, 90, 0},
					Position{100, 100, 0},
					Position{70, 70, 0},
				},
			},
		},
	}

	unmarshalResult := &Feature{}
	err := json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling Feature:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, polygon) {
		t.Fatalf("Unmarshalled polygon feature should match!\n Expected: %v\n Actual:   %v", polygon, unmarshalResult)
	}

	expected := `{"type":"Feature","id":"abc","bbox":[100,100,0,70,70,0],"geometry":{"coordinates":[[[90,90,0],[100,100,0],[70,70,0]]],"type":"Polygon"},"properties":{"FieldA":"blue","Test":"blargh"}}`
	marshalResult, err := polygon.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling Feature:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled polygon feature should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}

	geoJson = []byte(`{
		"type": "Feature",
		"id": "def",
		"bbox": [0, 0, 0, 0, 0, 0],
		"geometry": {
				"type": "Polygon",
				"coordinates": []
		},
		"properties": {
				"FieldA": "yellow",
				"Test": "blargh"
		}
	}`)
	polygon = &Feature{
		ID:         "def",
		Bbox:       Bbox{0, 0, 0, 0, 0, 0},
		Properties: map[string]string{"FieldA": "yellow", "Test": "blargh"},
		Geometry: &PolygonGeometry{
			Coordinates: PolygonCoords{},
		},
	}

	unmarshalResult = &Feature{}
	err = json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling Feature:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, polygon) {
		t.Fatalf("Unmarshalled linestring feature should match!\n Expected: %v\n Actual:   %v", polygon, unmarshalResult)
	}

	expected = `{"type":"Feature","id":"def","bbox":[0,0,0,0,0,0],"geometry":{"coordinates":[],"type":"Polygon"},"properties":{"FieldA":"yellow","Test":"blargh"}}`
	marshalResult, err = polygon.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling Feature:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled linestring feature should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}
}
func TestPolygonGeometryCollection(t *testing.T) {
	geoJson := []byte(`{
		"type": "GeometryCollection",
		"bbox": [100, 100, 0, 70, 70, 0],
		"geometries": [{
				"type": "Polygon",
				"coordinates": [[[90.0, 90.0, 0], [100.0, 100.0, 0], [70.0, 70.0, 0]]]
		}]
	}`)
	polygon := &GeometryCollection{
		Bbox: Bbox{100, 100, 0, 70, 70, 0},
		Geometries: []GeoJSON{
			PolygonGeometry{
				Coordinates: PolygonCoords{
					LineStringCoords{
						Position{90, 90, 0},
						Position{100, 100, 0},
						Position{70, 70, 0},
					},
				},
			},
		},
	}

	unmarshalResult := &GeometryCollection{}
	err := json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling GeometryCollection:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, polygon) {
		t.Fatalf("Unmarshalled polygon geometrycollection should match!\n Expected: %v\n Actual:   %v", polygon, unmarshalResult)
	}

	expected := `{"type":"GeometryCollection","geometries":[{"coordinates":[[[90,90,0],[100,100,0],[70,70,0]]],"type":"Polygon"}],"bbox":[100,100,0,70,70,0]}`
	marshalResult, err := polygon.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling GeometryCollection:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled polygon geometrycollection should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}

	geoJson = []byte(`{
		"type": "GeometryCollection",
		"bbox": [0, 0, 0, 0, 0, 0],
		"geometries": [{
				"type": "Polygon",
				"coordinates": []
		}]
	}`)
	polygon = &GeometryCollection{
		Bbox: Bbox{0, 0, 0, 0, 0, 0},
		Geometries: []GeoJSON{
			PolygonGeometry{
				Coordinates: PolygonCoords{},
			},
		},
	}

	unmarshalResult = &GeometryCollection{}
	err = json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling GeometryCollection:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, polygon) {
		t.Fatalf("Unmarshalled linestring geometrycollection should match!\n Expected: %v\n Actual:   %v", polygon, unmarshalResult)
	}

	expected = `{"type":"GeometryCollection","geometries":[{"coordinates":[],"type":"Polygon"}],"bbox":[0,0,0,0,0,0]}`
	marshalResult, err = polygon.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling GeometryCollection:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled linestring geometrycollection should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}
}

func TestIsValidPolygonCoords(t *testing.T) {
	coords := PolygonCoords{
		LineStringCoords{
			Position{0, 0, 0},
			Position{2, 1, 0},
			Position{1, 1, 0},
			Position{0, 0, 0},
		},
	}

	if !IsValidPolygonCoords(coords) {
		t.Fatalf("Simple polygon should be considered valid!\nActual: %v\n", coords)
	}

	coords = PolygonCoords{
		LineStringCoords{
			Position{0, 0, 0},
			Position{2, 1, 0},
			Position{1, 1, 0},
			Position{0, 0, 0},
		},
		LineStringCoords{
			Position{0.5, 0.4, 0},
			Position{0.8, 0.7, 0},
			Position{0.7, 0.5, 0},
			Position{0.5, 0.4, 0},
		},
	}

	if !IsValidPolygonCoords(coords) {
		t.Fatalf("Simple polygon with a hole should be considered valid!\nActual: %v\n", coords)
	}

	coords = PolygonCoords{
		LineStringCoords{
			Position{0, 0, 0},
			Position{1, 1, 0},
			Position{0, 0, 0},
		},
	}

	if IsValidPolygonCoords(coords) {
		t.Fatalf("Simple polygons should have at least 4 coordinates!\nActual: %v\n", coords)
	}

	coords = PolygonCoords{
		LineStringCoords{
			Position{0, 0, 0},
			Position{1, 1, 0},
			Position{1, 2, 0},
			Position{2, 2, 0},
		},
	}

	if IsValidPolygonCoords(coords) {
		t.Fatalf("The first and last element of a linear ring must match!\nActuan: %v\n", coords)
	}

	coords = PolygonCoords{
		LineStringCoords{
			Position{0, 0, 0},
			Position{1, 1, 0},
			Position{2, 1, 0},
			Position{0, 0, 0},
		},
	}

	if IsValidPolygonCoords(coords) {
		t.Fatalf("The first linear ring must go counter-clockwise!\nActual: %v\n", coords)
	}

	coords = PolygonCoords{
		LineStringCoords{
			Position{0, 0, 0},
			Position{2, 1, 0},
			Position{1, 1, 0},
			Position{0, 0, 0},
		},
		LineStringCoords{
			Position{0.5, 0.4, 0},
			Position{0.7, 0.5, 0},
			Position{0.8, 0.7, 0},
			Position{0.5, 0.4, 0},
		},
	}

	if IsValidPolygonCoords(coords) {
		t.Fatalf("The second linear ring must go clockwise!\nActual: %v\n", coords)
	}
}
