package main

import (
	"github.com/go-errors/errors"
	"testing"
)

func TestCollector_Report(t *testing.T) {
	c := Collector{}
	err := errors.New("testing")
	c.Report(&err.Err)

	if len(c.exceptions) != 1 {
		t.Errorf("expected one element")
	}

	if c.exceptions[0].Error.Cause != &err.Err {
		t.Errorf("failed")
	}

	if len(c.exceptions[0].Error.Stacktrace) == 0 {
		t.Errorf("expected a collected stack trace")
	}
}
