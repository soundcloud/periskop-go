package main

import (
	"strings"

	"github.com/go-errors/errors"
)

type Collector struct {
	exceptions []ErrorWithContext
}

func getStackTrace(err *error) []string {
	e := errors.New(err)
	trace := string(e.Stack())
	return strings.Split(trace, "\n")
}

func (c *Collector) Report(err *error) {
	c.report(err, HTTPContext{})
}

func (c *Collector) ReportWithContext(err *error, httpCtx HTTPContext) {
	c.report(err, httpCtx)
}

func (c *Collector) report(err *error, httpCtx HTTPContext) {
	errorInstance := NewErrorInstance(err, getStackTrace(err))
	errorWithContext := NewErrorWithContext(errorInstance, SeverityError, httpCtx)
	c.exceptions = append(c.exceptions, errorWithContext)
}
