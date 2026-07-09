package geotypes

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func TestMultiLineStringGeometry(t *testing.T) {
	geoJson := []byte(`{
		"type": "Point",
		"coordinates": [[[50.0, 40, 0], [90.0, 90.0, 0]]]
	}`)

	unmarshalResult := &MultiLineStringGeometry{}
	err := json.Unmarshal(geoJson, unmarshalResult)

	if err == nil {
		t.Fatal("We should error when parsing a Point as a MultiLineString!")
	}
	if !errors.Is(err, UnmarshallingTypeMismatch) {
		t.Fatalf("We should return an UnmarshallingTypeMismatch!\nActual: %v\n", err)
	}
	if !errors.Is(err, MultiLineStringGeometryUnmarshallingError) {
		t.Fatalf("We should return a MultiLineStringGeometryUnmarshallingError!\nActual: %v\n", err)
	}

	geoJson = []byte(`{
		"type": "MultiLineString",
		"coordinates": [[[50.0, 40, 0], [90.0, 90.0, 0]]]
	}`)
	multiLineString := &MultiLineStringGeometry{
		Coordinates: MultiLineStringCoords{LineStringCoords{Position{50.0, 40.0, 0}, Position{90.0, 90.0, 0}}},
	}

	unmarshalResult = &MultiLineStringGeometry{}
	err = json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling MultiLineStringGeometry:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, multiLineString) {
		t.Fatalf("Unmarshalled point geometry should match!\n Expected: %v\n Actual:   %v", multiLineString, unmarshalResult)
	}

	expected := `{"coordinates":[[[50,40,0],[90,90,0]]],"type":"MultiLineString"}`
	marshalResult, err := multiLineString.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling Feature:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled multilinestring feature should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}

	if !multiLineString.IsValid() {
		t.Fatalf("MultiLineString with valid coordinates should be valid!\nActual: %v\n", multiLineString.Coordinates)
	}

	multiLineString.Coordinates = MultiLineStringCoords{LineStringCoords{Position{50.0, 40.0, 0}}}

	if multiLineString.IsValid() {
		t.Fatalf("MultiLineString with invalid coordinates should not be valid!\nActual: %v\n", multiLineString.Coordinates)
	}
}

func TestMultiLineStringFeatures(t *testing.T) {
	geoJson := []byte(`{
		"type": "Feature",
		"id": "abc",
		"bbox": [100, 100, 0, 70, 70, 0],
		"geometry": {
				"type": "MultiLineString",
				"coordinates": [[[90.0, 90.0, 0], [100.0, 100.0, 0], [70.0, 70.0, 0]]]
		},
		"properties": {
				"FieldA": "blue",
				"Test": "blargh"
		}
	}`)
	multiLineString := &Feature{
		ID:         "abc",
		Bbox:       Bbox{100, 100, 0, 70, 70, 0},
		Properties: map[string]string{"FieldA": "blue", "Test": "blargh"},
		Geometry: &MultiLineStringGeometry{
			Coordinates: MultiLineStringCoords{
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

	if !reflect.DeepEqual(unmarshalResult, multiLineString) {
		t.Fatalf("Unmarshalled multilinestring feature should match!\n Expected: %v\n Actual:   %v", multiLineString, unmarshalResult)
	}

	expected := `{"type":"Feature","id":"abc","bbox":[100,100,0,70,70,0],"geometry":{"coordinates":[[[90,90,0],[100,100,0],[70,70,0]]],"type":"MultiLineString"},"properties":{"FieldA":"blue","Test":"blargh"}}`
	marshalResult, err := multiLineString.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling Feature:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled multilinestring feature should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}

	geoJson = []byte(`{
		"type": "Feature",
		"id": "def",
		"bbox": [0, 0, 0, 0, 0, 0],
		"geometry": {
				"type": "MultiLineString",
				"coordinates": []
		},
		"properties": {
				"FieldA": "yellow",
				"Test": "blargh"
		}
	}`)
	multiLineString = &Feature{
		ID:         "def",
		Bbox:       Bbox{0, 0, 0, 0, 0, 0},
		Properties: map[string]string{"FieldA": "yellow", "Test": "blargh"},
		Geometry: &MultiLineStringGeometry{
			Coordinates: MultiLineStringCoords{},
		},
	}

	unmarshalResult = &Feature{}
	err = json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling Feature:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, multiLineString) {
		t.Fatalf("Unmarshalled linestring feature should match!\n Expected: %v\n Actual:   %v", multiLineString, unmarshalResult)
	}

	expected = `{"type":"Feature","id":"def","bbox":[0,0,0,0,0,0],"geometry":{"coordinates":[],"type":"MultiLineString"},"properties":{"FieldA":"yellow","Test":"blargh"}}`
	marshalResult, err = multiLineString.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling Feature:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled linestring feature should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}
}

func TestMultiLineStringGeometryCollection(t *testing.T) {
	geoJson := []byte(`{
		"type": "GeometryCollection",
		"bbox": [100, 100, 0, 70, 70, 0],
		"geometries": [{
				"type": "MultiLineString",
				"coordinates": [[[90.0, 90.0, 0], [100.0, 100.0, 0], [70.0, 70.0, 0]]]
		}]
	}`)
	multiLineString := &GeometryCollection{
		Bbox: Bbox{100, 100, 0, 70, 70, 0},
		Geometries: []GeoJSON{
			MultiLineStringGeometry{
				Coordinates: MultiLineStringCoords{
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

	if !reflect.DeepEqual(unmarshalResult, multiLineString) {
		t.Fatalf("Unmarshalled multilinestring geometrycollection should match!\n Expected: %v\n Actual:   %v", multiLineString, unmarshalResult)
	}

	expected := `{"type":"GeometryCollection","geometries":[{"coordinates":[[[90,90,0],[100,100,0],[70,70,0]]],"type":"MultiLineString"}],"bbox":[100,100,0,70,70,0]}`
	marshalResult, err := multiLineString.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling geometrycollection:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled multilinestring geometrycollection should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}

	geoJson = []byte(`{
		"type": "GeometryCollection",
		"bbox": [0, 0, 0, 0, 0, 0],
		"geometries": [{
				"type": "MultiLineString",
				"coordinates": []
		}]
	}`)
	multiLineString = &GeometryCollection{
		Bbox: Bbox{0, 0, 0, 0, 0, 0},
		Geometries: []GeoJSON{
			MultiLineStringGeometry{
				Coordinates: MultiLineStringCoords{},
			},
		},
	}

	unmarshalResult = &GeometryCollection{}
	err = json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling GeometryCollection:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, multiLineString) {
		t.Fatalf("Unmarshalled linestring geometrycollection should match!\n Expected: %v\n Actual:   %v", multiLineString, unmarshalResult)
	}

	expected = `{"type":"GeometryCollection","geometries":[{"coordinates":[],"type":"MultiLineString"}],"bbox":[0,0,0,0,0,0]}`
	marshalResult, err = multiLineString.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling GeometryCollection:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled linestring geometrycollection should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}
}
