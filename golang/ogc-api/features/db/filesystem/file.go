package filesystem

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	geo "github.com/holdenparker/traipse/golang/geo-types"
	ogc "github.com/holdenparker/traipse/golang/ogc-api/features/schemas"
)

var (
	FileDBError                   = errors.New("FileDB Error:")
	FileDBRowNotParseError        = fmt.Errorf("%w: %w", FileDBError, errors.New("Unable to parse file into []fileDbRow"))
	FileDBNotLoadedError          = fmt.Errorf("%w: %w", FileDBError, errors.New("FileDB has not successfully loaded"))
	FileDBCollectionNotFoundError = fmt.Errorf("%w: %w", FileDBError, errors.New("Collection was not found"))
	FileDBItemNotFoundError       = fmt.Errorf("%w: %w", FileDBError, errors.New("Item was not found"))
)

type fileDbRow struct {
	ogc.Collection
	FeatureCollection geo.FeatureCollection
}

type FileDB struct {
	Filename    string
	collections []fileDbRow
	loaded      bool
}

func (fdb *FileDB) Load() error {
	file, err := os.Open(fdb.Filename)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("ERROR: File %v does not exist!\n", fdb.Filename)
			return fmt.Errorf("%w: %w", FileDBError, err)
		}
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&fdb.collections); err != nil {
		return fmt.Errorf("%w: %w", FileDBRowNotParseError, err)
	}
	fdb.loaded = true
	return nil
}

func (fdb *FileDB) GetCollectionIds() ([]string, error) {
	if !fdb.loaded {
		return []string{}, FileDBNotLoadedError
	}
	result := make([]string, len(fdb.collections))
	for col := range fdb.collections {
		result[col] = fdb.collections[col].Id
	}
	return result, nil
}

func (fdb *FileDB) HasCollection(collectionId string) (bool, error) {
	if !fdb.loaded {
		return false, FileDBNotLoadedError
	}
	for col := range fdb.collections {
		if fdb.collections[col].Id == collectionId {
			return true, nil
		}
	}
	return false, nil
}

func (fdb *FileDB) GetItems(collectionId string) (*geo.FeatureCollection, error) {
	if !fdb.loaded {
		return &geo.FeatureCollection{}, FileDBNotLoadedError
	}
	for col := range fdb.collections {
		if fdb.collections[col].Id == collectionId {
			return &fdb.collections[col].FeatureCollection, nil
		}
	}
	return nil, FileDBCollectionNotFoundError
}

func (fdb *FileDB) GetItem(collectionId string, itemId string) (*geo.Feature, error) {
	if !fdb.loaded {
		return &geo.Feature{}, FileDBNotLoadedError
	}
	items, err := fdb.GetItems(collectionId)
	if err != nil {
		return &geo.Feature{}, err
	}
	for i := range items.Features {
		if items.Features[i].ID == itemId {
			return &items.Features[i], nil
		}
	}
	return &geo.Feature{}, FileDBItemNotFoundError
}
