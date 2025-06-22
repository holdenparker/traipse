package geotypes

import (
	"errors"
	"fmt"
)

type MultiPointCoords = []PointCoords
type MultiPointGeometry geometryBuilder[MultiPointCoords]

var MultiPointGeometryUnmarshallingError = errors.New("LineString UnmarshalJSON error")

func (mp MultiPointGeometry) MarshalJSON() ([]byte, error) {
	return GeoJSONMarshalFactory(mp.Type(), mp.Coordinates)
}

func (p MultiPointGeometry) IsValid() bool {
	return true
}

func (mp MultiPointGeometry) Type() string {
	return MultiPointType
}

func (mp *MultiPointGeometry) UnmarshalJSON(data []byte) error {
	coords, err := GeoJSONUnmarshalFactory(mp.Type(), MultiPointCoords{}, data)

	if err != nil {
		return fmt.Errorf("%w: %w", MultiPointGeometryUnmarshallingError, err)
	}

	mp.Coordinates = coords

	return nil
}
