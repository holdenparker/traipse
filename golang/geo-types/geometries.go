package geotypes

import (
	"encoding/json"
	"reflect"
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

type geometryBuilder[CoordinateShape any] struct {
	Coordinates CoordinateShape `json:"coordinates"`
}

type GeoJSON interface {
	Type() string
}

type Position = [3]float64
type Bbox = [6]float64
