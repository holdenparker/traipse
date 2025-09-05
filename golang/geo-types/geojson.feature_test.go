package geotypes

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestFeatures(t *testing.T) {
	badData := []byte(`{"not": "real"}`)

	unmarshalResult := &Feature{}

	err := json.Unmarshal(badData, unmarshalResult)

	if err == nil {
		t.Fatal("We should error when trying to parse a bogus JSON as a Feature!")
	}
	if !errors.Is(err, FeatureUnmarshallingError) {
		t.Fatalf("We should be returning a FeatureUnmarshallingError!\nActual: %v\n", err)
	}

	badData = []byte(`{
		"type": "GeometryCollection",
		"geometry": {
			"type": "Point",
			"coordinates": [90.0, 90.0, 0]
		}
	}`)

	unmarshalResult = &Feature{}

	err = json.Unmarshal(badData, unmarshalResult)

	if err == nil {
		t.Fatal("We should error when trying to parse a GeometryCollection as a Feature!")
	}
	if !errors.Is(err, UnmarshallingFeatureTypeMismatch) {
		t.Fatalf("We should be returning an UnmarshallingFeatureTypeMismatch!\nActual: %v\n", err)
	}

	badData = []byte(`{
		"type": "Feature",
		"geometry": {
			"type": "Feature",
			"geometry": {
				"type": "MultiLineString",
				"coordinates": [[[50.0, 40, 0], [90.0, 90.0, 0]]]
			}
		}
	}`)

	unmarshalResult = &Feature{}

	err = json.Unmarshal(badData, unmarshalResult)

	if err == nil {
		t.Fatal("We should error when trying to parse a nested Feature!")
	}
	if !errors.Is(err, UnmarshallingFeatureUnsupportedGeometryType) {
		t.Fatalf("We should be returning an UnmarshallingFeatureUnsupportedGeometryType!\nActual: %v\n", err)
	}

}
