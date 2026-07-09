package features

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOAPI(t *testing.T) {
	oapi := OAPI{
		Title:       "Testing 1-2, 1-2",
		Version:     "4.5.6",
		OpenAPIPath: "/test",
		ServeDomain: "http://traipse.example.local",
	}

	oapi.BuildAPI()

	req := httptest.NewRequest(http.MethodGet, "/test.json", nil)
	res := httptest.NewRecorder()

	oapi.router.ServeHTTP(res, req)

	if res.Body.String() != `{"components":{"schemas":{}},"info":{"title":"Testing 1-2, 1-2","version":"4.5.6"},"openapi":"3.1.0"}` {
		t.Fatalf("Unexpected body response of: %v\n", res.Body)
	}
}
