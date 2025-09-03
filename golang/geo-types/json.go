package geotypes

import (
	"encoding/json"
	"errors"
)

type partialUnmarshal struct {
	Type        string            `json:"type"`
	Id          string            `json:"id,omitempty"`
	Bbox        Bbox              `json:"bbox,omitempty"`
	Coordinates json.RawMessage   `json:"coordinates,omitempty"`
	Features    []json.RawMessage `json:"features,omitempty"`
	Geometry    json.RawMessage   `json:"geometry,omitempty"`
	Geometries  []json.RawMessage `json:"geometries,omitempty"`
	Properties  map[string]string `json:"properties,omitempty"`
}

var (
	GeometryTypeNotKnownError            = errors.New("GeoJSON Geometry type not known!")
	NestedGeometryCollectionNotSupported = errors.New("Nested GeometryCollection Type Not Supported")
	GeoJSONTypeNotKnownError             = errors.New("GeoJSON type not known!")
)

func UnmarshalGeoJSON(b []byte) (GeoJSON, error) {
	geoj := &partialUnmarshal{}
	err := json.Unmarshal(b, geoj)
	if err != nil {
		return Feature{}, err
	}

	switch geoj.Type {
	case FeatureType:
		return UnmarshalGeoJSONFeature(geoj)
	case FeatureCollectionType:
		return UnmarshalGeoJSONFeatureCollection(geoj)
	case GeometryCollectionType:
		return UnmarshalGeoJSONGeometryCollection(geoj)
	default:
		return UnmarshalGeoJSONGeometry(geoj)
	}
}

func UnmarshalGeoJSONFeature(f *partialUnmarshal) (*Feature, error) {
	geom := &partialUnmarshal{}
	err := json.Unmarshal(f.Geometry, geom)
	if err != nil {
		return &Feature{}, err
	}
	var geometry GeoJSON
	if geom.Type == GeometryCollectionType {
		geometry, err = UnmarshalGeoJSONGeometryCollection(geom)
	} else {
		geometry, err = UnmarshalGeoJSONGeometry(geom)
	}
	if err != nil {
		return &Feature{}, err
	}
	return &Feature{
		ID:         f.Id,
		Bbox:       f.Bbox,
		Properties: f.Properties,
		Geometry:   geometry,
	}, nil
}

func UnmarshalGeoJSONFeatureCollection(f *partialUnmarshal) (*FeatureCollection, error) {
	result := FeatureCollection{
		ID:       f.Id,
		Features: make([]Feature, len(f.Features)),
		Bbox:     f.Bbox,
	}
	for i, featBin := range f.Features {
		partialFeat := &partialUnmarshal{}
		err := json.Unmarshal(featBin, partialFeat)
		if err != nil {
			return &FeatureCollection{}, err
		}
		feat, err := UnmarshalGeoJSONFeature(partialFeat)
		if err != nil {
			return &FeatureCollection{}, err
		}
		result.Features[i] = *feat
	}
	return &result, nil
}

func UnmarshalGeoJSONGeometryCollection(geom *partialUnmarshal) (*GeometryCollection, error) {
	result := GeometryCollection{
		Geometries: make([]GeoJSON, len(geom.Geometries)),
		Bbox:       geom.Bbox,
	}
	for i, geomBin := range geom.Geometries {
		partialGeom := &partialUnmarshal{}
		err := json.Unmarshal(geomBin, partialGeom)
		if err != nil {
			return &GeometryCollection{}, err
		}
		result.Geometries[i], err = UnmarshalGeoJSONGeometry(partialGeom)
	}
	return &result, nil
}

func UnmarshalGeoJSONGeometry(geom *partialUnmarshal) (GeoJSON, error) {
	switch geom.Type {
	case PointType:
		var coords PointCoords
		err := json.Unmarshal(geom.Coordinates, &coords)
		if err != nil {
			return PointGeometry{}, err
		}

		return PointGeometry{
			Coordinates: coords,
		}, nil
	case MultiPointType:
		var coords MultiPointCoords
		err := json.Unmarshal(geom.Coordinates, &coords)
		if err != nil {
			return MultiPointGeometry{}, err
		}

		return MultiPointGeometry{
			Coordinates: coords,
		}, nil
	case LineStringType:
		var coords LineStringCoords
		err := json.Unmarshal(geom.Coordinates, &coords)
		if err != nil {
			return &LineStringGeometry{}, err
		}

		return LineStringGeometry{
			Coordinates: coords,
		}, nil
	case MultiLineStringType:
		var coords MultiLineStringCoords
		err := json.Unmarshal(geom.Coordinates, &coords)
		if err != nil {
			return MultiLineStringGeometry{}, err
		}

		return MultiLineStringGeometry{
			Coordinates: coords,
		}, nil
	case PolygonType:
		var coords PolygonCoords
		err := json.Unmarshal(geom.Coordinates, &coords)
		if err != nil {
			return PolygonGeometry{}, err
		}

		return PolygonGeometry{
			Coordinates: coords,
		}, nil
	case MultiPolygonType:
		var coords MultiPolygonCoords
		err := json.Unmarshal(geom.Coordinates, &coords)
		if err != nil {
			return MultiPolygonGeometry{}, err
		}

		return MultiPolygonGeometry{
			Coordinates: coords,
		}, nil
	case GeometryCollectionType:
		return PointGeometry{}, NestedGeometryCollectionNotSupported
	}
	return PointGeometry{}, GeometryTypeNotKnownError
}
