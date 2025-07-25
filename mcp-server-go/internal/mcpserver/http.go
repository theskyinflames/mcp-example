package mcpserver

import (
	"net/http"

	_ "embed"
)

//go:embed openapi.json
var openAPISpec []byte

func startOpenAPIServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /docs/openapi.json", openaiSpecHndler())

	return http.ListenAndServe(":8091", mux)
}

func openaiSpecHndler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(openAPISpec)
	}
}
