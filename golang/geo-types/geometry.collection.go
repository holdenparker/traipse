package geotypes

import (
	"encoding/json"
	"errors"
	"fmt"
)

const GeometryCollectionType = "GeometryCollection"

var (
	GeometryCollectionUnmarshallingError                   = errors.New("GeometryCollection UnmarshalJSON error")
	UnmarshallingGeometryCollectionTypeMismatch            = errors.New("Expecting to unmarshal to GeometryCollection type")
	UnmarshallingGeometryCollectionUnsupportedGeometryType = errors.New("Provided geometry type is not supported")
)

type GeometryCollection struct {
	Geometries []GeoJSON `json:"geometries,omitempty"`
	Bbox       Bbox      `json:"bbox,omitempty"`
}

func (g GeometryCollection) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string    `json:"type"`
		Geometries []GeoJSON `json:"geometries,omitempty"`
		Bbox       Bbox      `json:"bbox,omitempty"`
	}{
		Type:       g.Type(),
		Geometries: g.Geometries,
		Bbox:       g.Bbox,
	})
}

func (g GeometryCollection) Type() string {
	return GeometryCollectionType
}

func (g *GeometryCollection) UnmarshalJSON(data []byte) error {
	result := &struct {
		Type       string            `json:"type"`
		Geometries []json.RawMessage `json:"geometries,omitempty"`
		Bbox       Bbox              `json:"bbox,omitempty"`
	}{}
	err := json.Unmarshal(data, result)
	if err != nil {
		return fmt.Errorf("%w: %w", GeometryCollectionUnmarshallingError, err)
	}
	if result.Type != GeometryCollectionType {
		return fmt.Errorf("%w: %w", GeometryCollectionUnmarshallingError, UnmarshallingGeometryCollectionTypeMismatch)
	}
	geometries := make([]GeoJSON, len(result.Geometries))
	// g.Geometries = make([]GeoJSON, len(result.Geometries))
	for i, geomBin := range result.Geometries {
		partialGeom := &struct {
			Type string `json:"type"`
		}{}
		err = json.Unmarshal(geomBin, partialGeom)
		if err != nil {
			return fmt.Errorf("%w: %w - %w", GeometryCollectionUnmarshallingError, UnmarshallingFeatureMissingGeometryType, err)
		}
		switch partialGeom.Type {
		case PointType:
			point := &PointGeometry{}
			err = json.Unmarshal(geomBin, point)
			geometries[i] = *point
		case MultiPointType:
			multiPoint := &MultiPointGeometry{}
			err = json.Unmarshal(geomBin, multiPoint)
			geometries[i] = *multiPoint
		case LineStringType:
			line := &LineStringGeometry{}
			err = json.Unmarshal(geomBin, line)
			geometries[i] = *line
		case MultiLineStringType:
			multiLine := &MultiLineStringGeometry{}
			err = json.Unmarshal(geomBin, multiLine)
			geometries[i] = *multiLine
		case PolygonType:
			polygon := &PolygonGeometry{}
			err = json.Unmarshal(geomBin, polygon)
			geometries[i] = *polygon
		case MultiPolygonType:
			multiPolygon := &MultiPolygonGeometry{}
			err = json.Unmarshal(geomBin, multiPolygon)
			geometries[i] = *multiPolygon
		default:
			err = fmt.Errorf("%w: %w", GeometryCollectionUnmarshallingError, UnmarshallingGeometryCollectionUnsupportedGeometryType)
		}
		if err != nil {
			return fmt.Errorf("%w: %w", GeometryCollectionUnmarshallingError, err)
		}
	}

	g.Bbox = result.Bbox
	g.Geometries = geometries

	return nil
}
