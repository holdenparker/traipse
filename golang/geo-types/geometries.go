package geotypes

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	PointType              = "Point"
	MultiPointType         = "MultiPoint"
	LineStringType         = "LineString"
	MultiLineStringType    = "MultiLineString"
	PolygonType            = "Polygon"
	MultiPolygonType       = "MultiPolygon"
	GeometryCollectionType = "GeometryCollection"
)

var UnmarshallingTypeMismatch = errors.New("Unmarshal type field mismatch")

type geometryBuilder[CoordinateShape any] struct {
	Coordinates CoordinateShape `json:"coordinates"`
}

func GeoJSONMarshalFactory[coordsType any](typeString string, coords coordsType) ([]byte, error) {
	// This type "rename" prevents circular unmarshalling
	type GeomType geometryBuilder[coordsType]
	return json.Marshal(struct {
		GeomType
		Type string `json:"type"`
	}{
		Type:     typeString,
		GeomType: GeomType{Coordinates: coords},
	})
}

func GeoJSONUnmarshalFactory[coordsType any](typeString string, dflt coordsType, data []byte) (coordsType, error) {
	type GeomType geometryBuilder[coordsType]
	result := &struct {
		GeomType
		Type string `json:"type"`
	}{}
	err := json.Unmarshal(data, result)
	if err != nil {
		return dflt, err
	}
	if result.Type != typeString {
		return dflt, fmt.Errorf("%w: Expected %v, found %v", UnmarshallingTypeMismatch, typeString, result.Type)
	}

	return result.Coordinates, nil
}

type GeoJSON interface {
	Type() string
}

type Position = [3]float64
type Bbox = [6]float64
