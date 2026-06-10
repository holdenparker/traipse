module github.com/holdenparker/traipse/golang/ogc-api/features

go 1.25.0

replace github.com/holdenparker/traipse/golang/algorithms => ../../algorithms

replace github.com/holdenparker/traipse/golang/geo-types => ../../geo-types

replace github.com/holdenparker/traipse/golang/ogc-api/features/schemas => ./schemas

require (
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/oapi-codegen/runtime v1.4.1 // indirect
)

require github.com/holdenparker/traipse/golang/algorithms v0.0.0-00010101000000-000000000000 // indirect

require github.com/holdenparker/traipse/golang/geo-types v0.0.0-00010101000000-000000000000

require (
	github.com/danielgtaylor/huma/v2 v2.38.0
	github.com/holdenparker/traipse/golang/ogc-api/features/schemas v0.0.0-00010101000000-000000000000
)
