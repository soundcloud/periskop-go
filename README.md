# periskop-go
Go client for Periskop

## Usage example


```go
package main

import (
	"encoding/json"
	"net/http"

	"github.com/soundcloud/periskop-go"
)

func faultyJSONParser() error {
	var dat map[string]interface{}
	// will return "unexpected end of JSON input"
	return json.Unmarshal([]byte(`{"id":`), &dat)
}

func main() {
	c := periskop.NewErrorCollector()

	// Without context
	c.Report(faultyJSONParser())

	// With HTTP context
	c.ReportWithContext(faultyJSONParser(), periskop.HTTPContext{
		RequestMethod:  "GET",
		RequestURL:     "http://example.com",
		RequestHeaders: map[string]string{"Cache-Control": "no-cache"},
	})

	// Call the exporter and HTTP handler to expose the
	// errors in /exceptions endpoints
	e := periskop.NewErrorExporter(&c)
	h := periskop.NewHandler(e)
	http.Handle("/exceptions", h)
	http.ListenAndServe(":8080", nil)
}
```
