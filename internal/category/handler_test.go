package category

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/theguuholi/rent-cars-go/internal/db/store"
	"net/http/httptest"
	"reflect"
	"testing"
)

type listFunc func(context.Context) ([]store.ListCategoriesRow, error)

func (f listFunc) ListCategories(ctx context.Context) ([]store.ListCategoriesRow, error) {
	return f(ctx)
}

func TestList(t *testing.T) {
	for _, tc := range []struct {
		name   string
		rows   []store.ListCategoriesRow
		err    error
		status int
		json   string
	}{
		{"empty catalog", nil, nil, 200, "[]"},
		{"nullable description", []store.ListCategoriesRow{{Name: "VAN"}}, nil, 200, `[{"name":"VAN","description":null}]`},
		{"database failure", nil, errors.New("private database details"), 500, `{"error":"Unable to load categories"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/categories", nil)
			w := httptest.NewRecorder()
			List(listFunc(func(ctx context.Context) ([]store.ListCategoriesRow, error) {
				if ctx != r.Context() {
					t.Error("lost request context")
				}
				return tc.rows, tc.err
			}))(w, r)
			if w.Code != tc.status || w.Header().Get("Content-Type") != "application/json" {
				t.Fatalf("unexpected response: %d %v", w.Code, w.Header())
			}
			var got, want any
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal([]byte(tc.json), &want); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("got %s, want %s", w.Body, tc.json)
			}
		})
	}
}
