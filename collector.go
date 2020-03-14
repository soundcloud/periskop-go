package periskop

import (
	"strings"
	"sync"

	"github.com/soundcloud/periskop-go/errutils"
)

// ErrorCollector collects all the aggregated errors
type ErrorCollector struct {
	aggregatedErrors sync.Map
}

// NewErrorCollector creates new ErrorCollector
func NewErrorCollector() ErrorCollector {
	return ErrorCollector{}
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

func (c *ErrorCollector) getAggregatedErrors() payload {
	aggregatedErrors := make([]aggregatedError, 0)
	c.aggregatedErrors.Range(func(key, value interface{}) bool {
		aggregatedErr, _ := value.(aggregatedError)
		aggregatedErrors = append(aggregatedErrors, aggregatedErr)
		return true
	})
	return payload{aggregatedErrors}
}

func (c *ErrorCollector) addError(err error, httpCtx HTTPContext) {
	errorInstance := newErrorInstance(err, getStackTrace(err))
	errorWithContext := newErrorWithContext(errorInstance, SeverityError, httpCtx)
	aggregationKey := errorWithContext.aggregationKey()
	if aggregatedErr, ok := c.aggregatedErrors.Load(aggregationKey); ok {
		aggregatedErr, _ := aggregatedErr.(*aggregatedError)
		aggregatedErr.addError(errorWithContext)
	} else {
		aggregatedErr := newErrorAggregate(aggregationKey, SeverityError)
		aggregatedErr.addError(errorWithContext)
		c.aggregatedErrors.Store(aggregationKey, &aggregatedErr)
	}
}
