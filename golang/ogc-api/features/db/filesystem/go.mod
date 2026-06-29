module github.com/holdenparker/traipse/golang/ogc-api/features/db/filesystem

go 1.25.0

replace github.com/holdenparker/traipse/golang/algorithms => ../../../../algorithms

replace github.com/holdenparker/traipse/golang/geo-types => ../../../../geo-types

replace github.com/holdenparker/traipse/golang/ogc-api/features/schemas => ../../schemas

require github.com/holdenparker/traipse/golang/algorithms v0.0.0-00010101000000-000000000000 // indirect

require github.com/holdenparker/traipse/golang/geo-types v0.0.0-00010101000000-000000000000

require github.com/holdenparker/traipse/golang/ogc-api/features/schemas v0.0.0-20260610021616-caae20528f5a

require (
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/oapi-codegen/runtime v1.4.2 // indirect
)
