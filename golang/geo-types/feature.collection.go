package geotypes

import (
	"encoding/json"
	"errors"
	"fmt"
)

const FeatureCollectionType = "FeatureCollection"

var (
	FeatureCollectionUnmarshallingError        = errors.New("FeatureCollection UnmarshalJSON error")
	UnmarshallingFeatureCollectionTypeMismatch = errors.New("Expecting to unmarshal to FeatureCollection type")
)

type FeatureCollection struct {
	ID       string    `json:"id,omitempty"`
	Features []Feature `json:"features"`
	Bbox     Bbox      `json:"bbox,omitempty"`
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

func (fc *FeatureCollection) UnmarshalJSON(data []byte) error {
	result := &struct {
		ID       string    `json:"id,omitempty"`
		Type     string    `json:"type"`
		Features []Feature `json:"features,omitempty"`
		Bbox     Bbox      `json:"bbox,omitempty"`
	}{}
	err := json.Unmarshal(data, result)
	if err != nil {
		return fmt.Errorf("%w: %w", FeatureCollectionUnmarshallingError, err)
	}
	if result.Type != FeatureCollectionType {
		return fmt.Errorf("%w: %w", FeatureCollectionUnmarshallingError, UnmarshallingFeatureCollectionTypeMismatch)
	}
	fc.ID = result.ID
	fc.Features = result.Features
	fc.Bbox = result.Bbox

	return nil
}
