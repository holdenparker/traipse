package geotypes

import (
	"errors"
	"fmt"
)

type PointCoords = Position
type PointGeometry geometryBuilder[PointCoords]

var PointGeometryUnmarshallingError = errors.New("PointGeometry UnmarshalJSON error")

func (p PointGeometry) IsValid() bool {
	return true
}

func (p PointGeometry) MarshalJSON() ([]byte, error) {
	return GeoJSONMarshalFactory(PointType, p.Coordinates)
}

func (p PointGeometry) Type() string {
	return PointType
}

func (p *PointGeometry) UnmarshalJSON(data []byte) error {
	coords, err := GeoJSONUnmarshalFactory(PointType, PointCoords{}, data)

	if err != nil {
		return fmt.Errorf("%w: %w", PointGeometryUnmarshallingError, err)
	}
	p.Coordinates = coords

	return nil
}
