package geotypes

import (
	"errors"
	"fmt"
)

const LineStringType = "LineString"

type LineStringCoords = []PointCoords
type LineStringGeometry geometryBuilder[LineStringCoords]

var LineStringGeometryUnmarshallingError = errors.New("LineString UnmarshalJSON error")

func (ls LineStringGeometry) MarshalJSON() ([]byte, error) {
	return GeoJSONMarshalFactory(ls.Type(), ls.Coordinates)
}

func (ls LineStringGeometry) IsValid() bool {
	return len(ls.Coordinates) > 1
}

func (ls LineStringGeometry) Type() string {
	return LineStringType
}

func (ls *LineStringGeometry) UnmarshalJSON(data []byte) error {
	coords, err := GeoJSONUnmarshalFactory(ls.Type(), LineStringCoords{}, data)

	if err != nil {
		return fmt.Errorf("%w: %w", LineStringGeometryUnmarshallingError, err)
	}

	ls.Coordinates = coords

	return nil
}
