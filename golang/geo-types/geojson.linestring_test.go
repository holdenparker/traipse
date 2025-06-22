package geotypes

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestLineStringGeometry(t *testing.T) {
	geoJson := []byte(`{
		"type": "LineString",
		"coordinates": [[50.0, 40, 0], [90.0, 90.0, 0]]
	}`)
	lineString := &LineStringGeometry{
		Coordinates: LineStringCoords{Position{50.0, 40.0, 0}, Position{90.0, 90.0, 0}},
	}

	unmarshalResult := &LineStringGeometry{}
	err := json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling LineStringGeometry:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, lineString) {
		t.Fatalf("Unmarshalled point geometry should match!\n Expected: %v\n Actual:   %v", lineString, unmarshalResult)
	}

	expected := `{"coordinates":[[50,40,0],[90,90,0]],"type":"LineString"}`
	marshalResult, err := lineString.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling Feature:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled linestring feature should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}
}

func TestLineStringFeatures(t *testing.T) {
	geoJson := []byte(`{
		"type": "Feature",
		"id": "abc",
		"bbox": [100, 100, 0, 70, 70, 0],
		"geometry": {
				"type": "LineString",
				"coordinates": [[90.0, 90.0, 0], [100.0, 100.0, 0], [70.0, 70.0, 0]]
		},
		"properties": {
				"FieldA": "blue",
				"Test": "blargh"
		}
	}`)
	lineString := &Feature{
		ID:         "abc",
		Bbox:       Bbox{100, 100, 0, 70, 70, 0},
		Properties: map[string]string{"FieldA": "blue", "Test": "blargh"},
		Geometry: &LineStringGeometry{
			Coordinates: LineStringCoords{
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

	if !reflect.DeepEqual(unmarshalResult, lineString) {
		t.Fatalf("Unmarshalled linestring feature should match!\n Expected: %v\n Actual:   %v", lineString, unmarshalResult)
	}

	expected := `{"type":"Feature","id":"abc","bbox":[100,100,0,70,70,0],"geometry":{"coordinates":[[90,90,0],[100,100,0],[70,70,0]],"type":"LineString"},"properties":{"FieldA":"blue","Test":"blargh"}}`
	marshalResult, err := lineString.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling Feature:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled linestring feature should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}

	geoJson = []byte(`{
		"type": "Feature",
		"id": "def",
		"bbox": [0, 0, 0, 0, 0, 0],
		"geometry": {
				"type": "LineString",
				"coordinates": []
		},
		"properties": {
				"FieldA": "yellow",
				"Test": "blargh"
		}
	}`)
	lineString = &Feature{
		ID:         "def",
		Bbox:       Bbox{0, 0, 0, 0, 0, 0},
		Properties: map[string]string{"FieldA": "yellow", "Test": "blargh"},
		Geometry: &LineStringGeometry{
			Coordinates: LineStringCoords{},
		},
	}

	unmarshalResult = &Feature{}
	err = json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling Feature:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, lineString) {
		t.Fatalf("Unmarshalled linestring feature should match!\n Expected: %v\n Actual:   %v", lineString, unmarshalResult)
	}

	expected = `{"type":"Feature","id":"def","bbox":[0,0,0,0,0,0],"geometry":{"coordinates":[],"type":"LineString"},"properties":{"FieldA":"yellow","Test":"blargh"}}`
	marshalResult, err = lineString.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling Feature:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled linestring feature should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}
}
