package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDocumentation(t *testing.T) {
	mux := http.NewServeMux()
	Docs(mux)
	for _, path := range []string{"/v3/api-docs", "/swagger-ui/", "/swagger-ui/swagger-ui-bundle.js", "/swagger-ui/swagger-ui.css"} {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		if w.Code != 200 {
			t.Fatalf("%s returned %d", path, w.Code)
		}
		if path == "/v3/api-docs" {
			var s struct {
				Paths map[string]map[string]json.RawMessage
			}
			if err := json.Unmarshal(w.Body.Bytes(), &s); err != nil || s.Paths["/categories"]["get"] == nil {
				t.Fatalf("invalid category contract: %v", err)
			}
		}
	}
}
