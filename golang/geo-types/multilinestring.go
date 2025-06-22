package geotypes

import (
	"errors"
	"fmt"
)

type MultiLineStringCoords = []LineStringCoords
type MultiLineStringGeometry geometryBuilder[MultiLineStringCoords]

var MultiLineStringGeometryUnmarshallingError = errors.New("MultiLineString UnmarshalJSON error")

func (mls MultiLineStringGeometry) MarshalJSON() ([]byte, error) {
	return GeoJSONMarshalFactory(mls.Type(), mls.Coordinates)
}

func (mls MultiLineStringGeometry) IsValid() bool {
	for _, ls := range mls.Coordinates {
		if len(ls) < 2 {
			return false
		}
	}
	return true
}

func (mls MultiLineStringGeometry) Type() string {
	return MultiLineStringType
}

func (ls *MultiLineStringGeometry) UnmarshalJSON(data []byte) error {
	coords, err := GeoJSONUnmarshalFactory(ls.Type(), MultiLineStringCoords{}, data)

	if err != nil {
		return fmt.Errorf("%w: %w", MultiLineStringGeometryUnmarshallingError, err)
	}

	ls.Coordinates = coords

	return nil
}
