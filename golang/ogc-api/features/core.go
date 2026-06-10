package features

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	geo "github.com/holdenparker/traipse/golang/geo-types"
	ogc "github.com/holdenparker/traipse/golang/ogc-api/features/schemas"
)

type LinkRelations string

func (lr LinkRelations) ptr() *string {
	return strPtr(string(lr))
}

const (
	Alternate   LinkRelations = "alternate"
	Collection  LinkRelations = "collection"
	DescribedBy LinkRelations = "describedby"
	Item        LinkRelations = "item"
	Next        LinkRelations = "next"
	License     LinkRelations = "license"
	Prev        LinkRelations = "prev"
	Self        LinkRelations = "self"
	ServiceDesc LinkRelations = "service-desc"
	ServiceDoc  LinkRelations = "service-doc"
	Items       LinkRelations = "items"
	Conformance LinkRelations = "conformance"
	Data        LinkRelations = "data"
)

type OgcApiFeaturesCoreDB interface {
	GetCollectionIds() ([]string, error)
	HasCollection(collectionId string) (bool, error)
	GetItems(collectionId string) (*geo.FeatureCollection, error)
	GetItem(collectionId string, itemId string) (*geo.Feature, error)
}

type OgcApiFeaturesCore struct {
	OAPI
	Description  string
	HostDomain   string
	Db           OgcApiFeaturesCoreDB
	landingLinks []ogc.Link
	conformances []string
}

type CollectionsRequest struct {
	CollectionId string `path:"collectionId"`
}

type FeatureRequest struct {
	CollectionId string `path:"collectionId"`
	FeatureId    string `path:"itemId"`
}

func OgcApiFeaturesCoreTransform(ctx huma.Context, status string, v any) (any, error) {
	switch feat := v.(type) {
	case geo.FeatureCollection:
		feat.FCType = geo.FeatureCollectionType
		return &feat, nil
	case *geo.FeatureCollection:
		feat.FCType = geo.FeatureCollectionType
		return feat, nil
	case geo.Feature:
		feat.FType = geo.FeatureType
		return &feat, nil
	case *geo.Feature:
		feat.FType = geo.FeatureType
		return feat, nil
	}
	return v, nil
}

func (ofc *OgcApiFeaturesCore) addToLandingLinks(links []ogc.Link) {
	ofc.landingLinks = append(ofc.landingLinks, links...)
}

func (ofc *OgcApiFeaturesCore) addToConformanceDeclarations(declarations []string) {
	ofc.conformances = append(ofc.conformances, declarations...)
}

func (ofc *OgcApiFeaturesCore) BuildCollection(collectionId string) *ogc.Collection {
	return &ogc.Collection{
		Id:       collectionId,
		ItemType: strPtr("feature"),
		Links: []ogc.Link{
			{
				Href: ofc.HostDomain + "/collections/" + collectionId,
				Rel:  Self.ptr(),
			},
			{
				Href: ofc.HostDomain + "/collections/" + collectionId + "/items",
				Rel:  Items.ptr(),
				Type: ApplicationGeoJSON.ptr(),
			},
		},
	}
}

func (ofc *OgcApiFeaturesCore) BuildAPI() {
	ofc.OAPI.addTransformers([]huma.Transformer{
		OgcApiFeaturesCoreTransform,
	})
	ofc.OAPI.BuildAPI()

	ofc.addToLandingLinks([]ogc.Link{
		{
			Href:  ofc.HostDomain,
			Rel:   Self.ptr(),
			Type:  ApplicationGeoJSON.ptr(),
			Title: strPtr("This document"),
		},
		{
			Href:  ofc.HostDomain + "/api.json",
			Rel:   ServiceDesc.ptr(),
			Type:  ApplicationOAIJson.ptr(),
			Title: strPtr("the API definition"),
		},
		{
			Href:  ofc.HostDomain + "/api.html",
			Rel:   ServiceDoc.ptr(),
			Type:  TextHTML.ptr(),
			Title: strPtr("the API documentation"),
		},
		{
			Href:  ofc.HostDomain + "/conformance",
			Rel:   Conformance.ptr(),
			Type:  ApplicationJSON.ptr(),
			Title: strPtr("OGC API conformance classes implemented by this server"),
		},
		{
			Href:  ofc.HostDomain + "/collections",
			Rel:   Data.ptr(),
			Type:  ApplicationJSON.ptr(),
			Title: strPtr("Information about the feature collections"),
		},
	})

	ofc.addToConformanceDeclarations([]string{
		"http://www.opengis.net/spec/ogcapi-features-1/1.0/conf/core",
		"http://www.opengis.net/spec/ogcapi-features-1/1.0/conf/oas30",
		"http://www.opengis.net/spec/ogcapi-features-1/1.0/conf/geojson",
	})

	huma.Get(ofc.api, "/", func(ctx context.Context, input *struct{}) (*responseBuilder[ogc.LandingPage], error) {
		return &responseBuilder[ogc.LandingPage]{
			Body: &ogc.LandingPage{
				Description: &ofc.Description,
				Links:       ofc.landingLinks,
				Title:       &ofc.Title,
			},
		}, nil
	})

	huma.Get(ofc.api, "/conformance", func(ctx context.Context, input *struct{}) (*responseBuilder[ogc.ConformanceDeclaration], error) {
		return &responseBuilder[ogc.ConfClasses]{
			Body: &ogc.ConformanceDeclaration{
				ConformsTo: ofc.conformances,
			},
		}, nil
	})

	huma.Get(ofc.api, "/collections", func(ctx context.Context, input *struct{}) (*responseBuilder[ogc.Collections], error) {
		collIds, err := ofc.Db.GetCollectionIds()
		colls := make([]ogc.Collection, len(collIds))
		for i, collId := range collIds {
			colls[i] = *ofc.BuildCollection(collId)
		}
		if err != nil {
			return &responseBuilder[ogc.Collections]{}, err
		}
		return &responseBuilder[ogc.Collections]{
			Body: &ogc.Collections{
				Collections: colls,
				Links: []ogc.Link{
					{
						Href: ofc.HostDomain + "/collections",
						Rel:  Self.ptr(),
						Type: ApplicationJSON.ptr(),
					},
				},
			},
		}, nil
	})

	huma.Get(ofc.api, "/collections/{collectionId}", func(ctx context.Context, input *CollectionsRequest) (*responseBuilder[ogc.Collection], error) {
		hasColl, err := ofc.Db.HasCollection(input.CollectionId)
		if err != nil {
			return &responseBuilder[ogc.Collection]{}, err
		}
		if hasColl != true {
			return &responseBuilder[ogc.Collection]{}, huma.Error404NotFound("Collection does not exist.")
		}
		return &responseBuilder[ogc.Collection]{
			Body: ofc.BuildCollection(input.CollectionId),
		}, nil
	})

	huma.Get(ofc.api, "/collections/{collectionId}/items", func(ctx context.Context, input *CollectionsRequest) (*responseBuilder[geo.FeatureCollection], error) {
		fc, err := ofc.Db.GetItems(input.CollectionId)
		if err != nil {
			return &responseBuilder[geo.FeatureCollection]{}, err
		}
		return &responseBuilder[geo.FeatureCollection]{
			Body: fc,
		}, nil
	})

	huma.Get(ofc.api, "/collections/{collectionId}/items/{itemId}", func(ctx context.Context, input *FeatureRequest) (*responseBuilder[geo.Feature], error) {
		feat, err := ofc.Db.GetItem(input.CollectionId, input.FeatureId)
		if err != nil {
			return nil, nil
		}
		return &responseBuilder[geo.Feature]{
			Body: feat,
		}, nil
	})
}
