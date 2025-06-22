package geotypes

import (
	"errors"
	"fmt"
	"reflect"
)

type MultiPolygonCoords = []PolygonCoords
type MultiPolygonGeometry geometryBuilder[MultiPolygonCoords]

var MultiPolygonGeometryUnmarshallingError = errors.New("MultiPolygon UnmarshalJSON error")

func (mp MultiPolygonGeometry) MarshalJSON() ([]byte, error) {
	return GeoJSONMarshalFactory(mp.Type(), mp.Coordinates)
}

func (mp MultiPolygonGeometry) IsValid() bool {
	// see PolygonGeometry.IsValid for complete notes
	for _, p := range mp.Coordinates {
		for _, lr := range p {
			if len(lr) < 4 {
				return false
			}
			if !reflect.DeepEqual(lr[0], lr[len(lr)-1]) {
				return false
			}
		}
	}
	return true
}

func (mp MultiPolygonGeometry) Type() string {
	return MultiPolygonType
}

func (mp *MultiPolygonGeometry) UnmarshalJSON(data []byte) error {
	coords, err := GeoJSONUnmarshalFactory(mp.Type(), MultiPolygonCoords{}, data)

	if err != nil {
		return fmt.Errorf("%w: %w", MultiPolygonGeometryUnmarshallingError, err)
	}

	mp.Coordinates = coords

	return nil
}
