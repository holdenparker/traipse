package geotypes

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestMultiPolygonGeometry(t *testing.T) {
	geoJson := []byte(`{
		"type": "MultiPolygon",
		"coordinates": [[[[50.0, 40, 0], [90.0, 90.0, 0]]]]
	}`)
	multiPolygon := &MultiPolygonGeometry{
		Coordinates: MultiPolygonCoords{PolygonCoords{LineStringCoords{Position{50.0, 40.0, 0}, Position{90.0, 90.0, 0}}}},
	}

	unmarshalResult := &MultiPolygonGeometry{}
	err := json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling MultiPolygon:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, multiPolygon) {
		t.Fatalf("Unmarshalled point geometry should match!\n Expected: %v\n Actual:   %v", multiPolygon, unmarshalResult)
	}

	expected := `{"coordinates":[[[[50,40,0],[90,90,0]]]],"type":"MultiPolygon"}`
	marshalResult, err := multiPolygon.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling Feature:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled multipolygon feature should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}
}

func TestMultiPolygonFeatures(t *testing.T) {
	geoJson := []byte(`{
		"type": "Feature",
		"id": "abc",
		"bbox": [100, 100, 0, 70, 70, 0],
		"geometry": {
				"type": "MultiPolygon",
				"coordinates": [[[[90.0, 90.0, 0], [100.0, 100.0, 0], [70.0, 70.0, 0]]]]
		},
		"properties": {
				"FieldA": "blue",
				"Test": "blargh"
		}
	}`)
	multiPolygon := &Feature{
		ID:         "abc",
		Bbox:       Bbox{100, 100, 0, 70, 70, 0},
		Properties: map[string]string{"FieldA": "blue", "Test": "blargh"},
		Geometry: &MultiPolygonGeometry{
			Coordinates: MultiPolygonCoords{
				PolygonCoords{
					LineStringCoords{
						Position{90, 90, 0},
						Position{100, 100, 0},
						Position{70, 70, 0},
					},
				},
			},
		},
	}

	unmarshalResult := &Feature{}
	err := json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling Feature:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, multiPolygon) {
		t.Fatalf("Unmarshalled multipolygon feature should match!\n Expected: %v\n Actual:   %v", multiPolygon, unmarshalResult)
	}

	expected := `{"type":"Feature","id":"abc","bbox":[100,100,0,70,70,0],"geometry":{"coordinates":[[[[90,90,0],[100,100,0],[70,70,0]]]],"type":"MultiPolygon"},"properties":{"FieldA":"blue","Test":"blargh"}}`
	marshalResult, err := multiPolygon.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling Feature:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled multipolygon feature should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}

	geoJson = []byte(`{
		"type": "Feature",
		"id": "def",
		"bbox": [0, 0, 0, 0, 0, 0],
		"geometry": {
				"type": "MultiPolygon",
				"coordinates": []
		},
		"properties": {
				"FieldA": "yellow",
				"Test": "blargh"
		}
	}`)
	multiPolygon = &Feature{
		ID:         "def",
		Bbox:       Bbox{0, 0, 0, 0, 0, 0},
		Properties: map[string]string{"FieldA": "yellow", "Test": "blargh"},
		Geometry: &MultiPolygonGeometry{
			Coordinates: MultiPolygonCoords{},
		},
	}

	unmarshalResult = &Feature{}
	err = json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling Feature:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, multiPolygon) {
		t.Fatalf("Unmarshalled linestring feature should match!\n Expected: %v\n Actual:   %v", multiPolygon, unmarshalResult)
	}

	expected = `{"type":"Feature","id":"def","bbox":[0,0,0,0,0,0],"geometry":{"coordinates":[],"type":"MultiPolygon"},"properties":{"FieldA":"yellow","Test":"blargh"}}`
	marshalResult, err = multiPolygon.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling Feature:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled linestring feature should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}
}

func TestMultiPolygonGeometryCollection(t *testing.T) {
	geoJson := []byte(`{
		"type": "GeometryCollection",
		"bbox": [100, 100, 0, 70, 70, 0],
		"geometries": [{
				"type": "MultiPolygon",
				"coordinates": [[[[90.0, 90.0, 0], [100.0, 100.0, 0], [70.0, 70.0, 0]]]]
		}]
	}`)
	multiPolygon := &GeometryCollection{
		Bbox: Bbox{100, 100, 0, 70, 70, 0},
		Geometries: []GeoJSON{
			MultiPolygonGeometry{
				Coordinates: MultiPolygonCoords{
					PolygonCoords{
						LineStringCoords{
							Position{90, 90, 0},
							Position{100, 100, 0},
							Position{70, 70, 0},
						},
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

	if !reflect.DeepEqual(unmarshalResult, multiPolygon) {
		t.Fatalf("Unmarshalled multipolygon geometrycollection should match!\n Expected: %v\n Actual:   %v", multiPolygon, unmarshalResult)
	}

	expected := `{"type":"GeometryCollection","geometries":[{"coordinates":[[[[90,90,0],[100,100,0],[70,70,0]]]],"type":"MultiPolygon"}],"bbox":[100,100,0,70,70,0]}`
	marshalResult, err := multiPolygon.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling GeometryCollection:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled multipolygon geometrycollection should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}

	geoJson = []byte(`{
		"type": "GeometryCollection",
		"bbox": [0, 0, 0, 0, 0, 0],
		"geometries": [{
				"type": "MultiPolygon",
				"coordinates": []
		}]
	}`)
	multiPolygon = &GeometryCollection{
		Bbox: Bbox{0, 0, 0, 0, 0, 0},
		Geometries: []GeoJSON{
			MultiPolygonGeometry{
				Coordinates: MultiPolygonCoords{},
			},
		},
	}

	unmarshalResult = &GeometryCollection{}
	err = json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling GeometryCollection:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, multiPolygon) {
		t.Fatalf("Unmarshalled linestring geometrycollection should match!\n Expected: %v\n Actual:   %v", multiPolygon, unmarshalResult)
	}

	expected = `{"type":"GeometryCollection","geometries":[{"coordinates":[],"type":"MultiPolygon"}],"bbox":[0,0,0,0,0,0]}`
	marshalResult, err = multiPolygon.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling GeometryCollection:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled linestring geometrycollection should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}
}
