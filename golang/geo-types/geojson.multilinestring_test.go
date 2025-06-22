package geotypes

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestMultiLineStringGeometry(t *testing.T) {
	geoJson := []byte(`{
		"type": "MultiLineString",
		"coordinates": [[[50.0, 40, 0], [90.0, 90.0, 0]]]
	}`)
	multiLineString := &MultiLineStringGeometry{
		Coordinates: MultiLineStringCoords{LineStringCoords{Position{50.0, 40.0, 0}, Position{90.0, 90.0, 0}}},
	}

	unmarshalResult := &MultiLineStringGeometry{}
	err := json.Unmarshal(geoJson, unmarshalResult)

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
