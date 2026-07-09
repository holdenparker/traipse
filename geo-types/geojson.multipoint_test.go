package geotypes

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func TestMultiPointGeometry(t *testing.T) {
	geoJson := []byte(`{
		"type": "Point",
		"coordinates": [[80.0, 80, 0], [90.0, 90.0, 0]]
	}`)

	unmarshalResult := &MultiPointGeometry{}
	err := json.Unmarshal(geoJson, unmarshalResult)

	if err == nil {
		t.Fatal("We should error when parsing a Point as a MultiPoint!")
	}
	if !errors.Is(err, UnmarshallingTypeMismatch) {
		t.Fatalf("We should return an UnmarshallingTypeMismatch!\nActual: %v\n", err)
	}
	if !errors.Is(err, MultiPointGeometryUnmarshallingError) {
		t.Fatalf("We should return a MultiPointGeometryUnmarshallingError!\nActual: %v\n", err)
	}

	geoJson = []byte(`{
		"type": "MultiPoint",
		"coordinates": [[80.0, 80, 0], [90.0, 90.0, 0]]
	}`)
	multiPoint := &MultiPointGeometry{
		Coordinates: MultiPointCoords{Position{80.0, 80.0, 0}, Position{90.0, 90.0, 0}},
	}

	unmarshalResult = &MultiPointGeometry{}
	err = json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling MultiPointGeometry:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, multiPoint) {
		t.Fatalf("Unmarshalled point geometry should match!\n Expected: %v\n Actual:   %v", multiPoint, unmarshalResult)
	}

	expected := `{"coordinates":[[80,80,0],[90,90,0]],"type":"MultiPoint"}`
	marshalResult, err := multiPoint.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling Feature:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled multipoint feature should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}

	if !multiPoint.IsValid() {
		t.Fatalf("Points should always be valid since the condition is enforced by the value type!\nActual: %v\n", multiPoint.Coordinates)
	}
}

func TestMultiPointFeatures(t *testing.T) {
	geoJson := []byte(`{
		"type": "Feature",
		"id": "abc",
		"bbox": [100, 100, 0, 70, 70, 0],
		"geometry": {
				"type": "MultiPoint",
				"coordinates": [[90.0, 90.0, 0], [100.0, 100.0, 0], [70.0, 70.0, 0]]
		},
		"properties": {
				"FieldA": "blue",
				"Test": "blargh"
		}
	}`)
	multiPoint := &Feature{
		ID:         "abc",
		Bbox:       Bbox{100, 100, 0, 70, 70, 0},
		Properties: map[string]string{"FieldA": "blue", "Test": "blargh"},
		Geometry: &MultiPointGeometry{
			Coordinates: MultiPointCoords{
				Position{90, 90, 0},
				Position{100, 100, 0},
				Position{70, 70, 0},
			},
		},
	}

	unmarshalResult := &Feature{}
	err := json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling Feature:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, multiPoint) {
		t.Fatalf("Unmarshalled multipoint feature should match!\n Expected: %v\n Actual:   %v", multiPoint, unmarshalResult)
	}

	expected := `{"type":"Feature","id":"abc","bbox":[100,100,0,70,70,0],"geometry":{"coordinates":[[90,90,0],[100,100,0],[70,70,0]],"type":"MultiPoint"},"properties":{"FieldA":"blue","Test":"blargh"}}`
	marshalResult, err := multiPoint.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling Feature:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled multipoint feature should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}

	geoJson = []byte(`{
		"type": "Feature",
		"id": "def",
		"bbox": [0, 0, 0, 0, 0, 0],
		"geometry": {
				"type": "MultiPoint",
				"coordinates": []
		},
		"properties": {
				"FieldA": "yellow",
				"Test": "blargh"
		}
	}`)
	multiPoint = &Feature{
		ID:         "def",
		Bbox:       Bbox{0, 0, 0, 0, 0, 0},
		Properties: map[string]string{"FieldA": "yellow", "Test": "blargh"},
		Geometry: &MultiPointGeometry{
			Coordinates: MultiPointCoords{},
		},
	}

	unmarshalResult = &Feature{}
	err = json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling Feature:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, multiPoint) {
		t.Fatalf("Unmarshalled multipoint feature should match!\n Expected: %v\n Actual:   %v", multiPoint, unmarshalResult)
	}

	expected = `{"type":"Feature","id":"def","bbox":[0,0,0,0,0,0],"geometry":{"coordinates":[],"type":"MultiPoint"},"properties":{"FieldA":"yellow","Test":"blargh"}}`
	marshalResult, err = multiPoint.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling Feature:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled multipoint feature should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}
}

func TestMultiPointGeometryCollection(t *testing.T) {
	geoJson := []byte(`{
		"type": "GeometryCollection",
		"bbox": [100, 100, 0, 70, 70, 0],
		"geometries": [{
				"type": "MultiPoint",
				"coordinates": [[90.0, 90.0, 0], [100.0, 100.0, 0], [70.0, 70.0, 0]]
		}]
	}`)
	multiPoint := &GeometryCollection{
		Bbox: Bbox{100, 100, 0, 70, 70, 0},
		Geometries: []GeoJSON{
			MultiPointGeometry{
				Coordinates: MultiPointCoords{
					Position{90, 90, 0},
					Position{100, 100, 0},
					Position{70, 70, 0},
				},
			},
		},
	}

	unmarshalResult := &GeometryCollection{}
	err := json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling GeometryCollection:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, multiPoint) {
		t.Fatalf("Unmarshalled multipoint geometrycollection should match!\n Expected: %v\n Actual:   %v", multiPoint, unmarshalResult)
	}

	expected := `{"type":"GeometryCollection","geometries":[{"coordinates":[[90,90,0],[100,100,0],[70,70,0]],"type":"MultiPoint"}],"bbox":[100,100,0,70,70,0]}`
	marshalResult, err := multiPoint.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling GeometryCollection:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled multipoint geometrycollection should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}

	geoJson = []byte(`{
		"type": "GeometryCollection",
		"bbox": [0, 0, 0, 0, 0, 0],
		"geometries": [{
				"type": "MultiPoint",
				"coordinates": []
		}]
	}`)
	multiPoint = &GeometryCollection{
		Bbox: Bbox{0, 0, 0, 0, 0, 0},
		Geometries: []GeoJSON{
			MultiPointGeometry{
				Coordinates: MultiPointCoords{},
			},
		},
	}

	unmarshalResult = &GeometryCollection{}
	err = json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling GeometryCollection:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, multiPoint) {
		t.Fatalf("Unmarshalled multipoint geometrycollection should match!\n Expected: %v\n Actual:   %v", multiPoint, unmarshalResult)
	}

	expected = `{"type":"GeometryCollection","geometries":[{"coordinates":[],"type":"MultiPoint"}],"bbox":[0,0,0,0,0,0]}`
	marshalResult, err = multiPoint.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling GeometryCollection:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled multipoint geometrycollection should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}
}
