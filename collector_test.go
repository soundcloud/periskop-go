package main

import (
	"testing"

	"github.com/go-errors/errors"
)

func getErrorAggregate(exceptions map[string]AggregatedError) AggregatedError {
	for _, errorAggregate := range exceptions {
		return errorAggregate
	}
	return AggregatedError{}
}

func TestCollector_addError(t *testing.T) {
	c := NewErrorCollector()
	err := errors.New("testing")
	c.addError(err, HTTPContext{})

	if len(c.aggregatedErrors) != 1 {
		t.Errorf("expected one element")
	}

	c.addError(err, HTTPContext{})
	if len(c.aggregatedErrors) != 2 {
		t.Errorf("expected two element")
	}
}

func TestCollector_Report(t *testing.T) {
	c := NewErrorCollector()
	err := errors.New("testing")
	c.Report(err.Err)

	if len(c.aggregatedErrors) != 1 {
		t.Errorf("expected one element")
	}

	errorWithContext := getErrorAggregate(c.aggregatedErrors).LatestErrors[0]
	if errorWithContext.Error.Cause != err.Err.Error() {
		t.Errorf("expected a propagated error")
	}

	if len(errorWithContext.Error.Stacktrace) == 0 {
		t.Errorf("expected a collected stack trace")
	}
}

func TestCollector_ReportWithContext(t *testing.T) {
	c := NewErrorCollector()
	err := errors.New("testing")
	httpContext := HTTPContext{
		RequestMethod:  "GET",
		RequestURL:     "http://example.com",
		RequestHeaders: map[string]string{"Cache-Control": "no-cache"},
	}
	c.ReportWithContext(err.Err, httpContext)

	if len(c.aggregatedErrors) != 1 {
		t.Errorf("expected one element")
	}

	errorWithContext := getErrorAggregate(c.aggregatedErrors).LatestErrors[0]
	if errorWithContext.HTTPContext.RequestMethod != "GET" {
		t.Errorf("expected HTTP method GET")
	}
}
