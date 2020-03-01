package periskop

import (
	"strings"

	"github.com/soundcloud/periskop-go/errutils"
)

// ErrorCollector collects all the aggregated errors
type ErrorCollector struct {
	aggregatedErrors map[string]*AggregatedError
}

// NewErrorCollector creates new ErrorCollector
func NewErrorCollector() ErrorCollector {
	return ErrorCollector{
		aggregatedErrors: make(map[string]*AggregatedError),
	}
}

// Report is used to add an error to map of aggregated errors
func (c *ErrorCollector) Report(err error) {
	c.addError(err, HTTPContext{})
}

// ReportWithContext is used to add an error (with HTTPContext) to map of aggregated errors
func (c *ErrorCollector) ReportWithContext(err error, httpCtx HTTPContext) {
	c.addError(err, httpCtx)
}

func getStackTrace(err error) []string {
	e := errutils.New(err)
	trace := string(e.Stack())
	s := strings.Split(trace, "\n")
	return s
}

func (c *ErrorCollector) getAggregatedErrors() Payload {
	aggregatedErrors := make([]AggregatedError, 0)
	for _, aggregateError := range c.aggregatedErrors {
		aggregatedErrors = append(aggregatedErrors, *aggregateError)
	}
	return Payload{AggregatedErrors: aggregatedErrors}
}

func (c *ErrorCollector) addError(err error, httpCtx HTTPContext) {
	errorInstance := NewErrorInstance(err, getStackTrace(err))
	errorWithContext := NewErrorWithContext(errorInstance, SeverityError, httpCtx)
	aggregationKey := errorWithContext.aggregationKey()
	if aggregatedError, ok := c.aggregatedErrors[aggregationKey]; ok {
		aggregatedError.addError(errorWithContext)
	} else {
		aggregatedError := NewErrorAggregate(aggregationKey, SeverityError)
		aggregatedError.addError(errorWithContext)
		c.aggregatedErrors[aggregationKey] = &aggregatedError
	}
}
