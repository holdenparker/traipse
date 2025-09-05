package geotypes

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestPointGeometry(t *testing.T) {
	geoJson := []byte(`{
		"type": "Point",
		"coordinates": [80.0, 80, 0]
	}`)
	point := &PointGeometry{
		Coordinates: Position{80.0, 80.0, 0},
	}

	unmarshalResult := &PointGeometry{}
	err := json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling PointGeometry:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, point) {
		t.Fatalf("Unmarshalled point geometry should match!\n Expected: %v\n Actual:   %v", point, unmarshalResult)
	}
}

func TestPointFeatures(t *testing.T) {
	geoJson := []byte(`{
		"type": "Feature",
		"id": "123",
		"bbox": [90, 90, 0, 90, 90, 0],
		"geometry": {
				"type": "Point",
				"coordinates": [90.0, 90.0, 0]
		},
		"properties": {
				"FieldA": "yellow",
				"Test": "blargh"
		}
	}`)
	point := &Feature{
		ID:         "123",
		Bbox:       Bbox{90, 90, 0, 90, 90, 0},
		Properties: map[string]string{"FieldA": "yellow", "Test": "blargh"},
		Geometry: &PointGeometry{
			Coordinates: Position{90, 90, 0},
		},
	}

	unmarshalResult := &Feature{}
	err := json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling Feature:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, point) {
		t.Fatalf("Unmarshalled point feature should match!\n Expected: %v\n Actual:   %v", point, unmarshalResult)
	}

	expected := `{"type":"Feature","id":"123","bbox":[90,90,0,90,90,0],"geometry":{"coordinates":[90,90,0],"type":"Point"},"properties":{"FieldA":"yellow","Test":"blargh"}}`
	marshalResult, err := point.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling Feature:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled point feature should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}

	geoJson = []byte(`{
		"type": "Feature",
		"id": "123",
		"bbox": [0, 0, 0, 0, 0, 0],
		"geometry": {
				"type": "Point",
				"coordinates": []
		},
		"properties": {
				"FieldA": "yellow",
				"Test": "blargh"
		}
	}`)
	point = &Feature{
		ID:         "123",
		Bbox:       Bbox{0, 0, 0, 0, 0, 0},
		Properties: map[string]string{"FieldA": "yellow", "Test": "blargh"},
		Geometry: &PointGeometry{
			Coordinates: Position{0, 0, 0},
		},
	}

	err = json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling Feature:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, point) {
		t.Fatalf("Unmarshalled point feature should match!\n Expected: %v\n Actual:   %v", point, unmarshalResult)
	}

	expected = `{"type":"Feature","id":"123","bbox":[0,0,0,0,0,0],"geometry":{"coordinates":[0,0,0],"type":"Point"},"properties":{"FieldA":"yellow","Test":"blargh"}}`
	marshalResult, err = point.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling Feature:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled point feature should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}
}

func TestPointGeometryCollection(t *testing.T) {
	geoJson := []byte(`{
		"type": "GeometryCollection",
		"bbox": [90, 90, 0, 90, 90, 0],
		"geometries": [{
				"type": "Point",
				"coordinates": [90.0, 90.0, 0]
		}]
	}`)
	point := &GeometryCollection{
		Bbox: Bbox{90, 90, 0, 90, 90, 0},
		Geometries: []GeoJSON{
			PointGeometry{
				Coordinates: Position{90, 90, 0},
			},
		},
	}

	unmarshalResult := &GeometryCollection{}
	err := json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling GeometryCollection:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, point) {
		t.Fatalf("Unmarshalled point geometrycollection should match!\n Expected: %v\n Actual:   %v", point, unmarshalResult)
	}

	expected := `{"type":"GeometryCollection","geometries":[{"coordinates":[90,90,0],"type":"Point"}],"bbox":[90,90,0,90,90,0]}`
	marshalResult, err := point.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling GeometryCollection:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled point geometrycollection should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}

	geoJson = []byte(`{
		"type": "GeometryCollection",
		"bbox": [0, 0, 0, 0, 0, 0],
		"geometries": [{
				"type": "Point",
				"coordinates": []
		}]
	}`)
	point = &GeometryCollection{
		Bbox: Bbox{0, 0, 0, 0, 0, 0},
		Geometries: []GeoJSON{
			PointGeometry{
				Coordinates: Position{0, 0, 0},
			},
		},
	}

	err = json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling GeometryCollection:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, point) {
		t.Fatalf("Unmarshalled point geometrycollection should match!\n Expected: %v\n Actual:   %v", point, unmarshalResult)
	}

	expected = `{"type":"GeometryCollection","geometries":[{"coordinates":[0,0,0],"type":"Point"}],"bbox":[0,0,0,0,0,0]}`
	marshalResult, err = point.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling GeometryCollection:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled point geometrycollection should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}
}
