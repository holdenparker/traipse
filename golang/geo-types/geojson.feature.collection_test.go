package geotypes

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func TestFeatureCollection(t *testing.T) {
	badData := []byte(`{"not": "real"}`)

	unmarshalResult := &FeatureCollection{}

	err := json.Unmarshal(badData, unmarshalResult)

	if err == nil {
		t.Fatal("We should error when trying to parse a bogus JSON as a FeatureCollection!")
	}
	if !errors.Is(err, FeatureCollectionUnmarshallingError) {
		t.Fatalf("We should be returning a FeatureCollectionUnmarshallingError!\nActual: %v\n", err)
	}

	badData = []byte(`{
		"type": "Feature",
		"geometry": {
			"type": "Point",
			"coordinates": [90.0, 90.0, 0]
		}
	}`)

	unmarshalResult = &FeatureCollection{}

	err = json.Unmarshal(badData, unmarshalResult)

	if err == nil {
		t.Fatal("We should error when trying to parse a Feature as a FeatureCollection!")
	}
	if !errors.Is(err, UnmarshallingFeatureCollectionTypeMismatch) {
		t.Fatalf("We should be returning a UnmarshallingFeatureCollectionTypeMismatch!\nActual: %v\n", err)
	}

	geoJson := []byte(`{
		"type": "FeatureCollection",
		"features": [{
			"type": "Feature",
			"id": "abc",
			"bbox": [90, 90, 0, 90, 90, 0],
			"geometry": {
				"type": "Point",
				"coordinates": [90.0, 90.0, 0]
			},
			"properties": {
					"FieldA": "blue",
					"Test": "blargh"
			}
		},{
			"type": "Feature",
			"id": "def",
			"bbox": [100, 100, 0, 70, 70, 0],
			"geometry": {
				"type": "GeometryCollection",
				"geometries": [
					{
						"type": "MultiLineString",
						"coordinates": [[[50.0, 40, 0], [90.0, 90.0, 0]]]
					}
				]
			}
		}]
	}`)
	featureCollection := &FeatureCollection{
		Features: []Feature{
			{
				ID:         "abc",
				Bbox:       Bbox{90, 90, 0, 90, 90, 0},
				Properties: map[string]string{"FieldA": "blue", "Test": "blargh"},
				Geometry: &PointGeometry{
					Coordinates: Position{90, 90, 0},
				},
			},
			{
				ID:   "def",
				Bbox: Bbox{100, 100, 0, 70, 70, 0},
				Geometry: &GeometryCollection{
					Geometries: []GeoJSON{
						MultiLineStringGeometry{
							Coordinates: MultiLineStringCoords{
								LineStringCoords{
									Position{50, 40, 0},
									Position{90, 90, 0},
								},
							},
						},
					},
				},
			},
		},
	}

	unmarshalResult = &FeatureCollection{}
	err = json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling FeatureCollection:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, featureCollection) {
		t.Fatalf("Unmarshalled featurecollection feature should match!\n Expected: %v\n Actual:   %v", featureCollection, unmarshalResult)
	}

	expected := `{"type":"FeatureCollection","features":[{"type":"Feature","id":"abc","bbox":[90,90,0,90,90,0],"geometry":{"coordinates":[90,90,0],"type":"Point"},"properties":{"FieldA":"blue","Test":"blargh"}},{"type":"Feature","id":"def","bbox":[100,100,0,70,70,0],"geometry":{"type":"GeometryCollection","geometries":[{"coordinates":[[[50,40,0],[90,90,0]]],"type":"MultiLineString"}],"bbox":[0,0,0,0,0,0]},"properties":null}]}`
	marshalResult, err := featureCollection.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling FeatureCollection:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled featurecollection should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}

	geoJson = []byte(`{
		"type": "FeatureCollection",
		"features": [{
			"type": "Feature",
			"id": "abc",
			"bbox": [0, 0, 0, 0, 0, 0],
			"geometry": {
				"type": "Point",
				"coordinates": []
			},
			"properties": {
					"FieldA": "blue",
					"Test": "blargh"
			}
		},{
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
			}
		}]
	}`)
	featureCollection = &FeatureCollection{
		Features: []Feature{
			{
				ID:         "abc",
				Bbox:       Bbox{0, 0, 0, 0, 0, 0},
				Properties: map[string]string{"FieldA": "blue", "Test": "blargh"},
				Geometry:   &PointGeometry{},
			},
			{
				ID:   "def",
				Bbox: Bbox{0, 0, 0, 0, 0, 0},
				Geometry: &GeometryCollection{
					Geometries: []GeoJSON{
						MultiLineStringGeometry{
							Coordinates: MultiLineStringCoords{},
						},
					},
				},
			},
		},
	}

	unmarshalResult = &FeatureCollection{}
	err = json.Unmarshal(geoJson, unmarshalResult)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling FeatureCollection:\n Error: %v", err)
	}

	if !reflect.DeepEqual(unmarshalResult, featureCollection) {
		t.Fatalf("Unmarshalled featurecollection should match!\n Expected: %v\n Actual:   %v\n Expected Point: %v\n Actual Point:   %v\n Expected GeometryCollection: %v\n Actual GeometryCollection:   %v\n", featureCollection, unmarshalResult, featureCollection.Features[0].Geometry, unmarshalResult.Features[0].Geometry, featureCollection.Features[1].Geometry, unmarshalResult.Features[1].Geometry)
	}

	expected = `{"type":"FeatureCollection","features":[{"type":"Feature","id":"abc","bbox":[0,0,0,0,0,0],"geometry":{"coordinates":[0,0,0],"type":"Point"},"properties":{"FieldA":"blue","Test":"blargh"}},{"type":"Feature","id":"def","bbox":[0,0,0,0,0,0],"geometry":{"type":"GeometryCollection","geometries":[{"coordinates":[],"type":"MultiLineString"}],"bbox":[0,0,0,0,0,0]},"properties":null}]}`
	marshalResult, err = featureCollection.MarshalJSON()

	if err != nil {
		t.Fatalf("Unexpected error marshalling FeatureCollection:\n Error: %v", err)
	}

	if expected != string(marshalResult) {
		t.Fatalf("Marshalled featurecollection should match!\n Expected: %v\n Actual:   %v", expected, string(marshalResult))
	}
}
