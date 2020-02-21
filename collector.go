package main

import (
	"strings"

	"github.com/go-errors/errors"
)

type Collector struct {
	exceptions map[string]ErrorAggregate
}

func NewCollector() Collector {
	return Collector{
		exceptions: make(map[string]ErrorAggregate),
	}
}

func getStackTrace(err error) []string {
	e := errors.New(err)
	trace := string(e.Stack())
	return strings.Split(trace, "\n")
}

func (c *Collector) Report(err error) {
	c.addError(err, HTTPContext{})
}

func (c *Collector) ReportWithContext(err error, httpCtx HTTPContext) {
	c.addError(err, httpCtx)
}

func (c *Collector) getExceptionAggregate() ExceptionAggregate {
	var errorAggregates []ErrorAggregate
	for _, errorAggregate := range c.exceptions {
		errorAggregates = append(errorAggregates, errorAggregate)
	}
	return ExceptionAggregate{ErrorAggregates: errorAggregates}
}

func (c *Collector) addError(err error, httpCtx HTTPContext) {
	errorInstance := NewErrorInstance(err, getStackTrace(err))
	errorWithContext := NewErrorWithContext(errorInstance, SeverityError, httpCtx)
	if errorAggregate, ok := c.exceptions[errorWithContext.aggregationKey()]; ok {
		errorAggregate.addError(errorWithContext)
	} else {
		errorAggregate = NewErrorAggregate(errorWithContext.aggregationKey(), SeverityError)
		errorAggregate.addError(errorWithContext)
		c.exceptions[errorWithContext.aggregationKey()] = errorAggregate
	}
}
