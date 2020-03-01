package periskop

import (
	"fmt"
	"net/http"
)

// NewHandler receives a Periskop Error Exporter and returns
// a handler with the exported errors in json format
func NewHandler(e ErrorExporter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		j, err := e.Export()
		if err != nil {
			fmt.Printf("error exporting Periskop errors %s", err)
		}
		w.Write([]byte(j))
	})
}
