package geotypes

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	FeatureType           = "Feature"
	FeatureCollectionType = "FeatureCollection"
)

var (
	FeatureUnmarshallingError                   = errors.New("Feature UnmarshalJSON error")
	UnmarshallingFeatureTypeMismatch            = errors.New("Expecting to unmarshal to Feature type")
	UnmarshallingFeatureMissingGeometryType     = errors.New("Geometry is missing type field")
	UnmarshallingFeatureUnsupportedGeometryType = errors.New("Provided geometry type is not supported")
)

type GeometryType = PointGeometry

type Feature struct {
	ID         string            `json:"id,omitempty"`
	Bbox       Bbox              `json:"bbox,omitempty"`
	Properties map[string]string `json:"properties"`
	Geometry   GeoJSON           `json:"geometry"`
}

func (f Feature) Type() string {
	return FeatureType
}

func (f Feature) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       string            `json:"type"`
		ID         string            `json:"id,omitempty"`
		Bbox       Bbox              `json:"bbox,omitempty"`
		Geometry   any               `json:"geometry"`
		Properties map[string]string `json:"properties"`
	}{
		Type:       f.Type(),
		ID:         f.ID,
		Bbox:       f.Bbox,
		Geometry:   f.Geometry,
		Properties: f.Properties,
	})
}

func (f *Feature) UnmarshalJSON(data []byte) error {
	result := &struct {
		Type       string            `json:"type"`
		ID         string            `json:"id,omitempty"`
		Bbox       Bbox              `json:"bbox,omitempty"`
		Geometry   json.RawMessage   `json:"geometry"`
		Properties map[string]string `json:"properties"`
	}{}
	err := json.Unmarshal(data, result)
	if err != nil {
		return fmt.Errorf("%w: %w", FeatureUnmarshallingError, err)
	}
	if result.Type != FeatureType {
		return fmt.Errorf("%w: %w", FeatureUnmarshallingError, UnmarshallingFeatureTypeMismatch)
	}
	partialGeom := &struct {
		Type string `json:"type"`
	}{}
	err = json.Unmarshal(result.Geometry, partialGeom)
	if err != nil {
		return fmt.Errorf("%w: %w - %w", FeatureUnmarshallingError, UnmarshallingFeatureMissingGeometryType, err)
	}
	var geom GeoJSON
	switch partialGeom.Type {
	case PointType:
		geom = &PointGeometry{}
		err = json.Unmarshal(result.Geometry, geom)
	case MultiPointType:
		geom = &MultiPointGeometry{}
		err = json.Unmarshal(result.Geometry, geom)
	case LineStringType:
		geom = &LineStringGeometry{}
		err = json.Unmarshal(result.Geometry, geom)
	case MultiLineStringType:
		geom = &MultiLineStringGeometry{}
		err = json.Unmarshal(result.Geometry, geom)
	case PolygonType:
		geom = &PolygonGeometry{}
		err = json.Unmarshal(result.Geometry, geom)
	case MultiPolygonType:
		geom = &MultiPolygonGeometry{}
		err = json.Unmarshal(result.Geometry, geom)
	default:
		err = fmt.Errorf("%w: %w", FeatureUnmarshallingError, UnmarshallingFeatureUnsupportedGeometryType)
	}
	if err != nil {
		return fmt.Errorf("%w: %w", FeatureUnmarshallingError, err)
	}

	f.ID = result.ID
	f.Bbox = result.Bbox
	f.Properties = result.Properties
	f.Geometry = geom

	return nil
}

type FeatureCollection struct {
	ID       string    `json:"id,omitempty"`
	Features []Feature `json:"features"`
	Bbox     Bbox      `json:"Bbox,omitempty"`
}

func (fc FeatureCollection) Type() string {
	return FeatureCollectionType
}

func (fc FeatureCollection) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type     string    `json:"type"`
		ID       string    `json:"id,omitempty"`
		Features []Feature `json:"features"`
	}{
		Type:     fc.Type(),
		ID:       fc.ID,
		Features: fc.Features,
	})
}
