package geotypes

import (
	"errors"
	"fmt"
)

const MultiLineStringType = "MultiLineString"

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

func (mls *MultiLineStringGeometry) UnmarshalJSON(data []byte) error {
	coords, err := GeoJSONUnmarshalFactory(mls.Type(), MultiLineStringCoords{}, data)

	if err != nil {
		return fmt.Errorf("%w: %w", MultiLineStringGeometryUnmarshallingError, err)
	}

	mls.Coordinates = coords

	return nil
}
