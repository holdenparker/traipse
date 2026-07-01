package filesystem

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	geo "github.com/holdenparker/traipse/golang/geo-types"
	ogc "github.com/holdenparker/traipse/golang/ogc-api/features/schemas"
)

func ptr[T any](v T) *T {
	return &v
}

var expectedDbCollections = []fileDbRow{
	{
		Collection: ogc.Collection{
			Id:          "city-parks",
			Title:       ptr("City Parks"),
			Description: ptr("Public park boundaries, points of interest, and managed recreational areas."),
			ItemType:    ptr("feature"),
			Crs:         &[]string{"http://opengis.net"},
			Extent: &ogc.Extent{
				Spatial: &struct {
					Bbox *[][]float32          `json:"bbox,omitempty"`
					Crs  *ogc.ExtentSpatialCrs `json:"crs,omitempty"`
				}{
					Bbox: &[][]float32{{-70.3225, 44.1102, -70.3195, 44.1121}},
				},
			},
			Links: []ogc.Link{
				{
					Href: "http://example.com",
					Rel:  ptr("self"),
					Type: ptr("application/json"),
				},
			},
		},
		FeatureCollection: geo.FeatureCollection{
			Features: []geo.Feature{
				{
					ID: "park-001",
					Geometry: &geo.PolygonGeometry{
						Coordinates: [][][3]float64{{
							{-70.3201, 44.1102, 0.0},
							{-70.3195, 44.1102, 0.0},
							{-70.3195, 44.1108, 0.0},
							{-70.3201, 44.1108, 0.0},
							{-70.3201, 44.1102, 0.0},
						}},
					},
					Properties: map[string]string{
						"name":  "Riverside Park",
						"acres": "15.5",
					},
				},
				{
					ID: "park-002",
					Geometry: &geo.PointGeometry{
						Coordinates: [3]float64{-70.3225, 44.1121, 0.0},
					},
					Properties: map[string]string{
						"name":  "Oak Grove",
						"acres": "5.2",
					},
				},
			},
		},
	},
	{
		Collection: ogc.Collection{
			Id:          "hiking-trails",
			Title:       ptr("Hiking Trails"),
			Description: ptr("Designated multi-use paths and recreational hiking trails."),
			ItemType:    ptr("feature"),
			Crs:         &[]string{"http://opengis.net"},
			Extent: &ogc.Extent{
				Spatial: &struct {
					Bbox *[][]float32          `json:"bbox,omitempty"`
					Crs  *ogc.ExtentSpatialCrs `json:"crs,omitempty"`
				}{
					Bbox: &[][]float32{{-70.325, 44.108, -70.3155, 44.1175}},
				},
			},
			Links: []ogc.Link{
				{
					Href: "http://example.com",
					Rel:  ptr("self"),
					Type: ptr("application/json"),
				},
			},
		},
		FeatureCollection: geo.FeatureCollection{
			Features: []geo.Feature{
				{
					ID: "trail-101",
					Geometry: &geo.LineStringGeometry{
						Coordinates: [][3]float64{
							{-70.325, 44.115, 0.0},
							{-70.324, 44.1165, 0.0},
							{-70.322, 44.1175, 0.0},
						},
					},
					Properties: map[string]string{
						"name":       "Riverbend Trail",
						"difficulty": "Easy",
					},
				},
				{
					ID: "trail-102",
					Geometry: &geo.LineStringGeometry{
						Coordinates: [][3]float64{
							{-70.318, 44.109, 0.0},
							{-70.317, 44.108, 0.0},
							{-70.3155, 44.1095, 0.0},
						},
					},
					Properties: map[string]string{
						"name":       "Bluff Loop",
						"difficulty": "Moderate",
					},
				},
			},
		},
	},
	{
		Collection: ogc.Collection{
			Id:          "city-bike-stations",
			Title:       ptr("City Bike Stations"),
			Description: ptr("Active micromobility and bike sharing docking stations."),
			ItemType:    ptr("feature"),
			Crs:         &[]string{"http://opengis.net"},
			Extent: &ogc.Extent{
				Spatial: &struct {
					Bbox *[][]float32          `json:"bbox,omitempty"`
					Crs  *ogc.ExtentSpatialCrs `json:"crs,omitempty"`
				}{
					Bbox: &[][]float32{{-70.3245, 44.1115, -70.319, 44.1142}},
				},
			},
			Links: []ogc.Link{
				{
					Href: "http://example.com",
					Rel:  ptr("self"),
					Type: ptr("application/json"),
				},
			},
		},
		FeatureCollection: geo.FeatureCollection{
			Features: []geo.Feature{
				{
					ID: "station-201",
					Geometry: &geo.PointGeometry{
						Coordinates: [3]float64{-70.319, 44.1115, 0.0},
					},
					Properties: map[string]string{
						"station_name": "Metro Hub East",
						"capacity":     "20",
					},
				},
				{
					ID: "station-202",
					Geometry: &geo.PointGeometry{
						Coordinates: [3]float64{-70.3245, 44.1142, 0.0},
					},
					Properties: map[string]string{
						"station_name": "Trailhead West",
						"capacity":     "12",
					},
				},
			},
		},
	},
}

func TestLoad(t *testing.T) {
	db := FileDB{
		Filename: "./test/does.not.exist.json",
	}
	if err := db.Load(); !errors.Is(err, FileDBError) && !os.IsNotExist(err) {
		t.Fatalf("Should fail to load a missing file!\n Error: %v\n", err)
	}
	db = FileDB{
		Filename: "./file_test.go",
	}
	if err := db.Load(); !errors.Is(err, FileDBRowNotParseError) {
		t.Fatalf("Should fail to parse file_test.go file as fileDbRow!\n Error: %v\n", err)
	}
	db = FileDB{
		Filename: "./test/test.collection.json",
	}
	if err := db.Load(); err != nil {
		t.Fatalf("Should successfully load relative path!\n Error: %v", err)
	}
	if !reflect.DeepEqual(db.collections, expectedDbCollections) {
		t.Fatalf("Expected db.collections to deep equal expected!\n Expected: %v\n Actual: %v\n", expectedDbCollections, db.collections)
	}

	absPath, err := filepath.Abs("./test/test.collection.json")
	if err != nil {
		t.Fatalf("Cannot generate absolute path!\n Error: %v\n", err)
	}
	db = FileDB{
		Filename: absPath,
	}
	if err := db.Load(); err != nil {
		t.Fatalf("Should successfully load absolute path!\n Error: %v\n", err)
	}
	if !reflect.DeepEqual(db.collections, expectedDbCollections) {
		t.Fatalf("Expected absolute path db.collections to deep equal expected!\n Expected: %v\n Actual: %v\n", expectedDbCollections, db.collections)
	}
}

func TestGetCollectionIds(t *testing.T) {
	db := FileDB{
		Filename: "./test/test.collection.json",
	}
	if _, err := db.GetCollectionIds(); !errors.Is(err, FileDBNotLoadedError) {
		t.Fatal("Should throw a FileDBNotLoadedError when db is not yet loaded!")
	}

	if err := db.Load(); err != nil {
		t.Fatalf("Should successfully load db!\n Error: %v\n", err)
	}

	expectedCollectionIds := make([]string, len(expectedDbCollections))
	for i := range expectedDbCollections {
		expectedCollectionIds[i] = expectedDbCollections[i].Id
	}
	collectionIds, err := db.GetCollectionIds()
	if err != nil {
		t.Fatalf("Unable to retrieve collectionIds!\n Error: %v", err)
	}
	if !reflect.DeepEqual(collectionIds, expectedCollectionIds) {
		t.Fatalf("Unexpected collectionIds!\n Expected: %v\n Actual: %v\n", expectedCollectionIds, collectionIds)
	}
}

func TestHasCollection(t *testing.T) {
	db := FileDB{
		Filename: "./test/test.collection.json",
	}
	expectedId := expectedDbCollections[1].Id
	if _, err := db.HasCollection(expectedId); !errors.Is(err, FileDBNotLoadedError) {
		t.Fatal("Should throw a FileDBNotLoadedError when db is not yet loaded!")
	}

	if err := db.Load(); err != nil {
		t.Fatalf("Should successfully load db!\n Error: %v\n", err)
	}

	hasId, err := db.HasCollection(expectedId)
	if err != nil {
		t.Fatalf("Unable to confirm valid collection id!\n Error: %v\n", err)
	}
	if !hasId {
		t.Fatalf("Expected to find collection Id: %v\n", expectedId)
	}

	hasId, err = db.HasCollection("bad-id")
	if err != nil {
		t.Fatalf("Unable to confirm invalid collection id!\n Error: %v", err)
	}
	if hasId {
		t.Fatal("Expected to not find collection id: bad-id")
	}
}

func TestGetItems(t *testing.T) {
	db := FileDB{
		Filename: "./test/test.collection.json",
	}
	collectionId := expectedDbCollections[2].Id
	expectedFC := &expectedDbCollections[2].FeatureCollection
	if _, err := db.GetItems(collectionId); !errors.Is(err, FileDBNotLoadedError) {
		t.Fatal("Should throw a FileDBNotLoadedError when db is not yet loaded!")
	}

	if err := db.Load(); err != nil {
		t.Fatalf("Should successfully load db!\n Error: %v\n", err)
	}

	items, err := db.GetItems(collectionId)
	if err != nil {
		t.Fatalf("Unable to retrieve valid collection items!\n Error: %v\n", err)
	}
	if !reflect.DeepEqual(items, expectedFC) {
		t.Fatalf("Unexpected items returned!\n Expected: %v\n Actual: %v\n", expectedFC, items)
	}

	if items, err = db.GetItems("bad-id"); !errors.Is(err, FileDBCollectionNotFoundError) {
		t.Fatalf("Expected to error on a bad collection id!\n Items: %v\n Error: %v", items, err)
	}
}

func TestGetItem(t *testing.T) {
	db := FileDB{
		Filename: "./test/test.collection.json",
	}
	collectionId := expectedDbCollections[2].Id
	expectedFeat := &expectedDbCollections[2].FeatureCollection.Features[1]
	if _, err := db.GetItem(collectionId, expectedFeat.ID); !errors.Is(err, FileDBNotLoadedError) {
		t.Fatal("Should throw a FileDBNotLoadedError when db is not yet loaded!")
	}

	if err := db.Load(); err != nil {
		t.Fatalf("Should successfully load db!\n Error: %v\n", err)
	}

	feature, err := db.GetItem(collectionId, expectedFeat.ID)
	if err != nil {
		t.Fatalf("Unable to retrieve valid collection and feature!\n $rror: %v\n", err)
	}
	if !reflect.DeepEqual(feature, expectedFeat) {
		t.Fatalf("Unexpected feature returned!\n Expected: %v\n Actual: %v\n", expectedFeat, feature)
	}

	if feature, err = db.GetItem("bad-id", expectedFeat.ID); !errors.Is(err, FileDBCollectionNotFoundError) {
		t.Fatalf("Expected to error on a bad collection id!\n Feature: %v\n Error: %v", feature, err)
	}

	if feature, err = db.GetItem(collectionId, "bad-id"); !errors.Is(err, FileDBItemNotFoundError) {
		t.Fatalf("Expected to error on a bad item id!\n Feature: %v\n Error: %v\n", feature, err)
	}
}
