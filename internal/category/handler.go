package category

import (
	"context"
	"encoding/json"
	"github.com/theguuholi/rent-cars-go/internal/db/store"
	"log/slog"
	"net/http"
)

type Lister interface {
	ListCategories(context.Context) ([]store.ListCategoriesRow, error)
}

type response struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

func List(categories Lister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := categories.ListCategories(r.Context())
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			slog.ErrorContext(r.Context(), "list categories", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Unable to load categories"})
			return
		}
		result := make([]response, 0, len(rows))
		for _, row := range rows {
			result = append(result, response{Name: row.Name, Description: row.Description})
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			slog.ErrorContext(r.Context(), "write categories response", "error", err)
		}
	}
}
