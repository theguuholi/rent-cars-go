// Package api serves the public contract and the interactive API documentation.
package api

import (
	_ "embed"
	swaggerfiles "github.com/swaggo/files/v2"
	"net/http"
)

//go:embed openapi.json
var spec []byte

func Docs(mux *http.ServeMux) {
	mux.HandleFunc("GET /v3/api-docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(spec)
	})
	mux.HandleFunc("GET /swagger-ui.html", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/swagger-ui/", http.StatusFound) })
	mux.HandleFunc("GET /swagger-ui/{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<!doctype html><html lang="en"><head><meta charset="utf-8"><title>Rent Cars API</title><link rel="stylesheet" href="swagger-ui.css"></head><body><div id="swagger-ui"></div><script src="swagger-ui-bundle.js"></script><script>SwaggerUIBundle({url:'/v3/api-docs',dom_id:'#swagger-ui'})</script></body></html>`))
	})
	mux.Handle("GET /swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.FS(swaggerfiles.FS))))
}
