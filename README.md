# periskop-go
Go client for Periskop

## Usage example


```go
package main

func faultJsonParsing() error {
    var dat map[string]interface{}
    // will return "unexpected end of JSON input"
	return json.Unmarshal([]byte(`{"num":`), &dat)
}

func main() {
    c := NewErrorCollector()

    // Without context
    c.Report(faultJsonParsing())
    
    // With HTTP context
	c.ReportWithContext(faultJsonParsing(), HTTPContext{
		RequestMethod:  "GET",
		RequestURL:     "http://example.com",
		RequestHeaders: map[string]string{"Cache-Control": "no-cache"},
    })
    
    // Call the exporter and HTTP handler to expose the 
    // errors in /exceptions endpoints
	e := NewErrorExporter(&c)
	h := NewHandler(e)
	http.Handle("/exceptions", h)
	http.ListenAndServe(":8080", nil)
}
```
