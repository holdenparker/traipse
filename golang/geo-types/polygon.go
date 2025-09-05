package geotypes

import (
	"errors"
	"fmt"
	"reflect"
)

const PolygonType = "Polygon"

type PolygonCoords = []LineStringCoords
type PolygonGeometry geometryBuilder[PolygonCoords]

var PolygonGeometryUnmarshallingError = errors.New("Polygon UnmarshalJSON error")

func (p PolygonGeometry) MarshalJSON() ([]byte, error) {
	return GeoJSONMarshalFactory(p.Type(), p.Coordinates)
}

func (p PolygonGeometry) IsValid() bool {
	for _, lr := range p.Coordinates {
		// linear rings require 4 or more positions
		if len(lr) < 4 {
			return false
		}
		// The first and last element of a linear ring must match
		if !reflect.DeepEqual(lr[0], lr[len(lr)-1]) {
			return false
		}
		// TODO: The first linear ring must describe the bounds of the surface,
		// with the following linear rings defining holes within that surface.
		// In other words, all following linear rings must be completely contained
		// within the first linear ring.
	}
	return true
}

func (p PolygonGeometry) Type() string {
	return PolygonType
}

func (p *PolygonGeometry) UnmarshalJSON(data []byte) error {
	coords, err := GeoJSONUnmarshalFactory(p.Type(), PolygonCoords{}, data)

	if err != nil {
		return fmt.Errorf("%w: %w", PolygonGeometryUnmarshallingError, err)
	}

	p.Coordinates = coords

	return nil
}
