package geotypes

import (
	"encoding/json"
	"errors"
	"fmt"
)

type PointCoords = Position
type PointGeometry geometryBuilder[PointCoords]

var (
	UnmarshallingPointTypeMismatch  = errors.New("Expecting to unmarshal to Point geometry type")
	PointGeometryUnmarshallingError = errors.New("PointGeometry UnmarshalJSON error")
)

func (p PointGeometry) IsValid() bool {
	return true
}

func (p PointGeometry) MarshalJSON() ([]byte, error) {
	// This type "rename" prevents circular unmarshalling
	type Pg PointGeometry
	return json.Marshal(struct {
		Pg
		Type string `json:"type"`
	}{
		Type: PointType,
		Pg:   Pg{Coordinates: p.Coordinates},
	})
}

func (p PointGeometry) Type() string {
	return PointType
}

func (p *PointGeometry) UnmarshalJSON(data []byte) error {
	// This type "rename" prevents circular unmarshalling
	type Pg PointGeometry
	result := &struct {
		Pg
		Type string `json:"type"`
	}{}
	err := json.Unmarshal(data, result)
	if err != nil {
		return fmt.Errorf("%w: %w", PointGeometryUnmarshallingError, err)
	}
	if result.Type != PointType {
		return fmt.Errorf("%w: %w, found %s", PointGeometryUnmarshallingError, UnmarshallingPointTypeMismatch, result.Type)
	}

	p.Coordinates = result.Coordinates

	return nil
}
