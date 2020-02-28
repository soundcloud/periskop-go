package main

import (
	"net/http"
)

// NewHandler receives a Periskop Error Exporter and returns
// a handler with the exported errors in json format
func NewHandler(e ErrorExporter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		j, _ := e.Export()
		w.Write([]byte(j))
	})
}
