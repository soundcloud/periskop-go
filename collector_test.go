package main

import (
	"testing"

	"github.com/go-errors/errors"
)

func TestCollector_Report(t *testing.T) {
	c := Collector{}
	err := errors.New("testing")
	c.Report(err.Err)

	if len(c.exceptions) != 1 {
		t.Errorf("expected one element")
	}

	if c.exceptions[0].Error.Cause != err.Err {
		t.Errorf("expected a propagated error")
	}

	if len(c.exceptions[0].Error.Stacktrace) == 0 {
		t.Errorf("expected a collected stack trace")
	}
}

func TestCollector_ReportWithContext(t *testing.T) {
	c := Collector{}
	err := errors.New("testing")
	httpContext := HTTPContext{
		RequestMethod:  "GET",
		RequestURL:     "http://example.com",
		RequestHeaders: map[string]string{"Cache-Control": "no-cache"},
	}
	c.ReportWithContext(err.Err, httpContext)

	if len(c.exceptions) != 1 {
		t.Errorf("expected one element")
	}

	if c.exceptions[0].HTTPContext.RequestMethod != "GET" {
		t.Errorf("expected HTTP method GET")
	}
}
