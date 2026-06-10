package features

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

type OAPI struct {
	Title        string
	Version      string
	OpenAPIPath  string
	ServeDomain  string
	router       *http.ServeMux
	api          huma.API
}

func (oa *OAPI) addTransformers(transformers []huma.Transformer) {
	oa.transformers = append(oa.transformers, transformers...)
}

func (oa *OAPI) BuildAPI() {
	oa.router = http.NewServeMux()
	apiConfig := huma.DefaultConfig(oa.Title, oa.Version)
	apiConfig.OpenAPIPath = oa.OpenAPIPath
	oa.api = humago.New(oa.router, apiConfig)
}

func (oa *OAPI) StartAPI() {
	http.ListenAndServe(oa.ServeDomain, oa.router)
}
