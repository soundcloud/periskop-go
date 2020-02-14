package main

import (
	crashy "github.com/go-errors/errors"
	"strings"
)

type Collector struct {
	exceptions []ErrorWithContext
}

func getStackTrace(err *error) []string {
	crash := crashy.New(err)
	trace := string(crash.Stack())
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
	errorWithContext := NewErrorWithContext(errorInstance, SeverityError)
	c.exceptions = append(c.exceptions, errorWithContext)
}
