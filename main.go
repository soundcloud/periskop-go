package main

import (
	"encoding/json"
	"net/http"
)

func faultyFunction() error {
	var dat map[string]interface{}
	return json.Unmarshal([]byte(`{"id":`), &dat)
}

func main() {
	c := NewErrorCollector()
	c.Report(faultyFunction())
	c.Report(faultyFunction())
	c.Report(faultyFunction())

	e := NewErrorExporter(&c)
	h := NewHandler(e)
	http.Handle("/exceptions", h)
	http.ListenAndServe(":8080", nil)
}
