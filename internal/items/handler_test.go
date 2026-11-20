package items

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func do(mux http.Handler, method, path, body string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(method, path, strings.NewReader(body)))
	return rec
}

func TestItems(t *testing.T) {
	mux := http.NewServeMux()
	NewHandler(NewStore()).Register(mux)

	tests := []struct {
		name, method, path, body string
		want                     int
	}{
		{"create", http.MethodPost, "/api/items", `{"name":"first"}`, http.StatusCreated},
		{"create without name", http.MethodPost, "/api/items", `{"name":" "}`, http.StatusBadRequest},
		{"bad json", http.MethodPost, "/api/items", `{`, http.StatusBadRequest},
		{"get existing", http.MethodGet, "/api/items/1", "", http.StatusOK},
		{"get missing", http.MethodGet, "/api/items/2", "", http.StatusNotFound},
		{"list", http.MethodGet, "/api/items", "", http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if rec := do(mux, tt.method, tt.path, tt.body); rec.Code != tt.want {
				t.Fatalf("%s %s = %d, want %d (body %s)", tt.method, tt.path, rec.Code, tt.want, rec.Body)
			}
		})
	}
}
