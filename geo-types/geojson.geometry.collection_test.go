package geotypes

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func TestGeometryCollection(t *testing.T) {
	badData := []byte(`{"not": "real"}`)

	unmarshalResult := &GeometryCollection{}

	err := json.Unmarshal(badData, unmarshalResult)

	if err == nil {
		t.Fatal("We should error when trying to parse a bogus JSON as a GeometryCollection!")
	}
	if !errors.Is(err, GeometryCollectionUnmarshallingError) {
		t.Fatalf("We should be returning a GeometryCollectionUnmarshallingError!\nActual: %v\n", err)
	}

	badData = []byte(`{
		"type": "Feature",
		"geometries": [{
			"type": "Point",
			"coordinates": [90.0, 90.0, 0]
		}]
	}`)

	unmarshalResult = &GeometryCollection{}

	err = json.Unmarshal(badData, unmarshalResult)

	if err == nil {
		t.Fatal("We should error when trying to parse a Feature as a GeometryCollection!")
	}
	if !errors.Is(err, UnmarshallingGeometryCollectionTypeMismatch) {
		t.Fatalf("We should be returning an UnmarshallingGeometryCollectionTypeMismatch!\nActual: %v\n", err)
	}

	badData = []byte(`{
		"type": "GeometryCollection",
		"geometries": [
			{
				"type": "GeometryCollection",
				"geometries": [
					{
						"type": "MultiLineString",
						"coordinates": [[[50.0, 40, 0], [90.0, 90.0, 0]]]
					}
				]
			}
		]
	}`)

	unmarshalResult = &GeometryCollection{}

	err = json.Unmarshal(badData, unmarshalResult)

	if err == nil {
		t.Fatal("We should error when trying to parse a nested GeometryCollection!")
	}
	if !errors.Is(err, UnmarshallingGeometryCollectionUnsupportedGeometryType) {
		t.Fatalf("We should be returning an UnmarshallingGeometryCollectionUnsupportedGeometryType!\nActual: %v\n", err)
	}

	geoJson := []byte(`{
		"type": "GeometryCollection",
		"geometries": [
			{
				"type": "MultiLineString",
				"coordinates": [[[50.0, 40, 0], [90.0, 90.0, 0]]]
			}
		]
	}`)
	geometryCollection := &GeometryCollection{
		Geometries: []GeoJSON{
			MultiLineStringGeometry{
				Coordinates: MultiLineStringCoords{LineStringCoords{Position{50.0, 40.0, 0}, Position{90.0, 90.0, 0}}},
			},
		},
	}

	unmarshalResult = &GeometryCollection{}
	err = json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling GeometryCollection:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, geometryCollection) {
		t.Fatalf("Unmarshalled geometry collection should match!\n Expected: %v\n Actual:   %v", geometryCollection, unmarshalResult)
	}

	expected := `{"type":"GeometryCollection","geometries":[{"coordinates":[[[50,40,0],[90,90,0]]],"type":"MultiLineString"}],"bbox":[0,0,0,0,0,0]}`
	marshalResult, err := geometryCollection.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling Feature:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled GeometryCollection feature should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}
}

func TestGeometryCollectionFeatures(t *testing.T) {
	geoJson := []byte(`{
		"type": "Feature",
		"id": "abc",
		"bbox": [100, 100, 0, 70, 70, 0],
		"geometry": {
				"type": "GeometryCollection",
				"geometries": [
					{
						"type": "MultiLineString",
						"coordinates": [[[90.0, 90.0, 0], [100.0, 100.0, 0], [70.0, 70.0, 0]]]
					}
				]
		},
		"properties": {
				"FieldA": "blue",
				"Test": "blargh"
		}
	}`)
	geometryCollection := &Feature{
		ID:         "abc",
		Bbox:       Bbox{100, 100, 0, 70, 70, 0},
		Properties: map[string]string{"FieldA": "blue", "Test": "blargh"},
		Geometry: &GeometryCollection{
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
		},
	}

	unmarshalResult := &Feature{}
	err := json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling Feature:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, geometryCollection) {
		t.Fatalf("Unmarshalled geometrycollection feature should match!\n Expected: %v\n Actual:   %v", geometryCollection, unmarshalResult)
	}

	expected := `{"type":"Feature","id":"abc","bbox":[100,100,0,70,70,0],"geometry":{"type":"GeometryCollection","geometries":[{"coordinates":[[[90,90,0],[100,100,0],[70,70,0]]],"type":"MultiLineString"}],"bbox":[0,0,0,0,0,0]},"properties":{"FieldA":"blue","Test":"blargh"}}`
	marshalResult, err := geometryCollection.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling Feature:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled geometrycollection feature should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}

	geoJson = []byte(`{
		"type": "Feature",
		"id": "def",
		"bbox": [0, 0, 0, 0, 0, 0],
		"geometry": {
				"type": "GeometryCollection",
				"geometries": [
					{
						"type": "MultiLineString",
						"coordinates": []
					}
				]
		},
		"properties": {
				"FieldA": "yellow",
				"Test": "blargh"
		}
	}`)
	geometryCollection = &Feature{
		ID:         "def",
		Bbox:       Bbox{0, 0, 0, 0, 0, 0},
		Properties: map[string]string{"FieldA": "yellow", "Test": "blargh"},
		Geometry: &GeometryCollection{
			Geometries: []GeoJSON{
				MultiLineStringGeometry{
					Coordinates: MultiLineStringCoords{},
				},
			},
		},
	}

	unmarshalResult = &Feature{}
	err = json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling Feature:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, geometryCollection) {
		t.Fatalf("Unmarshalled linestring feature should match!\n Expected: %v\n Actual:   %v", geometryCollection, unmarshalResult)
	}

	expected = `{"type":"Feature","id":"def","bbox":[0,0,0,0,0,0],"geometry":{"type":"GeometryCollection","geometries":[{"coordinates":[],"type":"MultiLineString"}],"bbox":[0,0,0,0,0,0]},"properties":{"FieldA":"yellow","Test":"blargh"}}`
	marshalResult, err = geometryCollection.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling Feature:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled linestring feature should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}
}
