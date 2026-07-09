package features

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	geotypes "github.com/holdenparker/traipse/geo-types"
	"github.com/holdenparker/traipse/ogc-api/features/schemas"
)

type TestCoreDB struct {
	collections []geotypes.FeatureCollection
}

func (tcdb *TestCoreDB) GetCollectionIds() ([]string, error) {
	colls := make([]string, len(tcdb.collections))
	for i, coll := range tcdb.collections {
		colls[i] = coll.ID
	}
	return colls, nil
}

func (tcdb *TestCoreDB) HasCollection(collectionId string) (bool, error) {
	result := false
	for _, coll := range tcdb.collections {
		if coll.ID == collectionId {
			result = true
		}
	}
	return result, nil
}

func (tcdb *TestCoreDB) GetItems(collectionId string) (*geotypes.FeatureCollection, error) {
	var result *geotypes.FeatureCollection
	for i := 0; i < len(tcdb.collections); i++ {
		if tcdb.collections[i].ID == collectionId {
			result = &tcdb.collections[i]
		}
	}
	if result != nil {
		return result, nil
	}
	return nil, errors.New("Invalid collection id")
}

func (tcdb *TestCoreDB) GetItem(collectionId string, itemId string) (*geotypes.Feature, error) {
	var result *geotypes.Feature
	coll, err := tcdb.GetItems(collectionId)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(coll.Features); i++ {
		if coll.Features[i].ID == itemId {
			result = &coll.Features[i]
		}
	}
	if result != nil {
		return result, nil
	}
	return nil, errors.New("Invalid feature id")
}

func getExampleService() OgcApiFeaturesCore {
	ogcapi := OgcApiFeaturesCore{
		OAPI: OAPI{
			Title:       "This is a test",
			Version:     "7.8.9",
			OpenAPIPath: "/api",
			ServeDomain: "http://traipse.other.example.local",
		},
		Description: "The awesome API",
		HostDomain:  "http://traipse.host.example.local",
		Db: &TestCoreDB{
			collections: []geotypes.FeatureCollection{
				{
					ID: "Stuff",
					Features: []geotypes.Feature{
						{
							ID: "one",
							Properties: map[string]string{
								"Hello": "world",
							},
							Geometry: geotypes.PointGeometry{
								Coordinates: geotypes.PointCoords{0.0, 0.0, 0.0},
							},
						},
					},
				},
			},
		},
	}

	ogcapi.BuildAPI()

	return ogcapi
}

func TestOgcApiFeaturesCoreGetLandingPage(t *testing.T) {
	ogcapi := getExampleService()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	ogcapi.router.ServeHTTP(res, req)

	expected := &schemas.LandingPage{
		Title:       strPtr("This is a test"),
		Description: strPtr("The awesome API"),
		Links: []schemas.Link{
			{
				Href:  "http://traipse.host.example.local",
				Rel:   Self.ptr(),
				Type:  ApplicationGeoJSON.ptr(),
				Title: strPtr("This document"),
			},
			{
				Href:  "http://traipse.host.example.local/api.json",
				Rel:   ServiceDesc.ptr(),
				Type:  ApplicationOAIJson.ptr(),
				Title: strPtr("the API definition"),
			},
			{
				Href:  "http://traipse.host.example.local/api.html",
				Rel:   ServiceDoc.ptr(),
				Type:  TextHTML.ptr(),
				Title: strPtr("the API documentation"),
			},
			{
				Href:  "http://traipse.host.example.local/conformance",
				Rel:   Conformance.ptr(),
				Type:  ApplicationJSON.ptr(),
				Title: strPtr("OGC API conformance classes implemented by this server"),
			},
			{
				Href:  "http://traipse.host.example.local/collections",
				Rel:   Data.ptr(),
				Type:  ApplicationJSON.ptr(),
				Title: strPtr("Information about the feature collections"),
			},
		},
	}
	unmarshalResult := &schemas.LandingPage{}
	err := json.Unmarshal(res.Body.Bytes(), unmarshalResult)
	if err != nil {
		t.Fatalf("Unexpected error unmarshalling LandingPate:\n Error: %v\n Body: %v", err, res.Body.String())
	}
	if !reflect.DeepEqual(unmarshalResult, expected) {
		t.Fatalf("GET / response does not match expected value!\n Expected: %v\n Actual: %v\n", expected, unmarshalResult)
	}
}

func TestOgcApiFeaturesCoreGetConformanceDeclarations(t *testing.T) {
	ogcapi := getExampleService()

	req := httptest.NewRequest(http.MethodGet, "/conformance", nil)
	res := httptest.NewRecorder()

	ogcapi.router.ServeHTTP(res, req)

	expected := &schemas.ConformanceDeclaration{
		ConformsTo: []string{
			"http://www.opengis.net/spec/ogcapi-features-1/1.0/conf/core",
			"http://www.opengis.net/spec/ogcapi-features-1/1.0/conf/oas30",
			"http://www.opengis.net/spec/ogcapi-features-1/1.0/conf/geojson",
		},
	}
	unmarshalResult := &schemas.ConformanceDeclaration{}
	err := json.Unmarshal(res.Body.Bytes(), unmarshalResult)
	if err != nil {
		t.Fatalf("Unexpected error unmarshalling LandingPate:\n Error: %v\n Body: %v", err, res.Body.String())
	}
	if !reflect.DeepEqual(unmarshalResult, expected) {
		t.Fatalf("GET / response does not match expected value!\n Expected: %v\n Actual: %v\n", expected, unmarshalResult)
	}
}

func TestOgcApiFeaturesCoreGetItems(t *testing.T) {
	ogcapi := getExampleService()

	req := httptest.NewRequest(http.MethodGet, "/collections/Stuff/items", nil)
	res := httptest.NewRecorder()

	ogcapi.router.ServeHTTP(res, req)

	expected := &geotypes.FeatureCollection{
		ID: "Stuff",
		Features: []geotypes.Feature{
			{
				ID: "one",
				Properties: map[string]string{
					"Hello": "world",
				},
				Geometry: &geotypes.PointGeometry{
					Coordinates: geotypes.PointCoords{0.0, 0.0, 0.0},
				},
			},
		},
	}
	unmarshalResult := &geotypes.FeatureCollection{}
	err := json.Unmarshal(res.Body.Bytes(), unmarshalResult)
	if err != nil {
		t.Fatalf("Unexpected error unmarshalling Feature:\n Error: %v\n Body: %v", err, res.Body.String())
	}
	if !reflect.DeepEqual(unmarshalResult, expected) {
		t.Fatalf("Response does not match expected value!\n Expected: %v\n Actual: %v\n", expected, unmarshalResult)
	}
}

func TestOgcApiFeaturesCoreGetCollections(t *testing.T) {
	ogcapi := getExampleService()

	req := httptest.NewRequest(http.MethodGet, "/collections", nil)
	res := httptest.NewRecorder()

	ogcapi.router.ServeHTTP(res, req)

	expected := &schemas.Collections{
		Collections: []schemas.Collection{
			{
				Id:       "Stuff",
				ItemType: strPtr("feature"),
				Links: []schemas.Link{
					{
						Href: "http://traipse.host.example.local/collections/Stuff",
						Rel:  Self.ptr(),
					},
					{
						Href: "http://traipse.host.example.local/collections/Stuff/items",
						Rel:  Items.ptr(),
						Type: ApplicationGeoJSON.ptr(),
					},
				},
			},
		},
		Links: []schemas.Link{
			{
				Href: "http://traipse.host.example.local/collections",
				Rel:  Self.ptr(),
				Type: ApplicationJSON.ptr(),
			},
		},
	}
	unmarshalResult := &schemas.Collections{}
	err := json.Unmarshal(res.Body.Bytes(), unmarshalResult)
	if err != nil {
		t.Fatalf("Unexpected error unmarshalling Feature:\n Error: %v\n Body: %v", err, res.Body.String())
	}
	if !reflect.DeepEqual(unmarshalResult, expected) {
		t.Fatalf("Response does not match expected value!\n Expected: %v\n Actual: %v\n", expected, unmarshalResult)
	}
}

func TestOgcApiFeaturesCoreGetCollection(t *testing.T) {
	ogcapi := getExampleService()

	req := httptest.NewRequest(http.MethodGet, "/collections/Stuff", nil)
	res := httptest.NewRecorder()

	ogcapi.router.ServeHTTP(res, req)

	expected := &schemas.Collection{
		Id:       "Stuff",
		ItemType: strPtr("feature"),
		Links: []schemas.Link{
			{
				Href: "http://traipse.host.example.local/collections/Stuff",
				Rel:  Self.ptr(),
			},
			{
				Href: "http://traipse.host.example.local/collections/Stuff/items",
				Rel:  Items.ptr(),
				Type: ApplicationGeoJSON.ptr(),
			},
		},
	}
	unmarshalResult := &schemas.Collection{}
	err := json.Unmarshal(res.Body.Bytes(), unmarshalResult)
	if err != nil {
		t.Fatalf("Unexpected error unmarshalling Feature:\n Error: %v\n Body: %v", err, res.Body.String())
	}
	if !reflect.DeepEqual(unmarshalResult, expected) {
		t.Fatalf("Response does not match expected value!\n Expected: %v\n Actual: %v\n", expected, unmarshalResult)
	}
}

func TestOgcApiFeaturesCoreGetItem(t *testing.T) {
	ogcapi := getExampleService()

	req := httptest.NewRequest(http.MethodGet, "/collections/Stuff/items/one", nil)
	res := httptest.NewRecorder()

	ogcapi.router.ServeHTTP(res, req)

	expected := &geotypes.Feature{
		ID: "one",
		Properties: map[string]string{
			"Hello": "world",
		},
		Geometry: &geotypes.PointGeometry{
			Coordinates: geotypes.PointCoords{0.0, 0.0, 0.0},
		},
	}
	unmarshalResult := &geotypes.Feature{}
	err := json.Unmarshal(res.Body.Bytes(), unmarshalResult)
	if err != nil {
		t.Fatalf("Unexpected error unmarshalling Feature:\n Error: %v\n Body: %v", err, res.Body.String())
	}
	if !reflect.DeepEqual(unmarshalResult, expected) {
		t.Fatalf("Response does not match expected value!\n Expected: %v\n Actual: %v\n", expected, unmarshalResult)
	}
}
