package periskop

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/soundcloud/periskop-go/errutils"
)

// ErrorCollector collects all the aggregated errors
type ErrorCollector struct {
	aggregatedErrors map[string]*aggregatedError
	mux              sync.RWMutex
	uuid             uuid.UUID
}

// NewErrorCollector creates new ErrorCollector
func NewErrorCollector() ErrorCollector {
	return ErrorCollector{
		aggregatedErrors: make(map[string]*aggregatedError),
		uuid:             uuid.New(),
	}
}

// Report adds an error with severity Error to map of aggregated errors
func (c *ErrorCollector) Report(err error, errKey ...string) {
	c.ReportWithSeverity(err, SeverityError, errKey...)
}

// ReportWithSeverity adds an error with given severity to map of aggregated errors
func (c *ErrorCollector) ReportWithSeverity(err error, severity Severity, errKey ...string) {
	c.addError(err, severity, nil, errKey...)
}

// ReportWithHTTPContext adds an error with severity Error (with HTTPContext) to map of aggregated errors
func (c *ErrorCollector) ReportWithHTTPContext(err error, httpCtx *HTTPContext, errKey ...string) {
	c.ReportWithHTTPContextAndSeverity(err, SeverityError, httpCtx, errKey...)
}

// ReportWithHTTPContextAndSeverity adds an error with given severity (with HTTPContext) to map of aggregated errors
func (c *ErrorCollector) ReportWithHTTPContextAndSeverity(err error, severity Severity, httpCtx *HTTPContext,
	errKey ...string) {
	c.addError(err, severity, httpCtx, errKey...)
}

// ReportWithHTTPRequest adds and error with severity Error  (with HTTPContext from http.Request) to map
// of aggregated errors
func (c *ErrorCollector) ReportWithHTTPRequest(err error, r *http.Request, errKey ...string) {
	c.ReportWithHTTPRequestAndSeverity(err, SeverityError, r, errKey...)
}

// ReportWithHTTPRequestAndSeverity adds and error with given severity (with HTTPContext from http.Request) to
// map of aggregated errors
func (c *ErrorCollector) ReportWithHTTPRequestAndSeverity(err error, severity Severity, r *http.Request,
	errKey ...string) {
	c.addError(err, severity,
		&HTTPContext{
			RequestMethod:  r.Method,
			RequestURL:     r.URL.String(),
			RequestHeaders: getAllHeaders(r.Header),
			RequestBody:    getBody(r.Body),
		}, errKey...)
}

// ReportErrorWithContext adds a manually generated errorWithContext with an specific to map of aggregated errors
func (c *ErrorCollector) ReportErrorWithContext(errWithContext ErrorWithContext, severity Severity, errKey ...string) {
	c.addErrorWithContext(errWithContext, severity, errKey...)
}

// getBody reads io.Reader request body and returns either body converted to a string or a nil
func getBody(body io.Reader) *string {
	if body == nil {
		return nil
	}
	r, err := ioutil.ReadAll(body)
	bodyAsString := string(r)
	if err != nil {
		return nil
	}
	return &bodyAsString
}

// getAllHeaders gets all the headers of HTTP Request
func getAllHeaders(h http.Header) map[string]string {
	headersMap := make(map[string]string)
	for name, values := range h {
		for _, value := range values {
			headersMap[name] = value
		}
	}
	return headersMap
}

// getStackTrace gets the trace of the reported error
func getStackTrace(err error) []string {
	e := errutils.New(err)
	// get all the traces produced by the error skipping those
	// traces generated by this package.
	trace := string(e.Stack("periskop-go"))
	s := strings.FieldsFunc(trace, func(c rune) bool { return c == '\n' })
	return s
}

func (c *ErrorCollector) getAggregatedErrors() payload {
	c.mux.RLock()
	defer c.mux.RUnlock()
	aggregatedErrors := make([]aggregatedError, 0)
	for _, value := range c.aggregatedErrors {
		aggregatedErrors = append(aggregatedErrors, *value)
	}
	return payload{aggregatedErrors, c.uuid}
}

// getAggregationKey gets the aggregation key of the error
// Specifying 'errKey' you bypass the default aggregation method
func getAggregationKey(errorWithContext ErrorWithContext, errKey ...string) string {
	if len(errKey) > 0 {
		// aggregate also by error type
		return fmt.Sprintf("%s@%s", errorWithContext.Error.Class, errKey[0])
	}
	return errorWithContext.aggregationKey()
}

func (c *ErrorCollector) addError(err error, severity Severity, httpCtx *HTTPContext, errKey ...string) {
	errorInstance := newErrorInstance(err, reflect.TypeOf(err).String(), getStackTrace(err))
	errWithContext := NewErrorWithContext(errorInstance, severity, httpCtx)
	c.addErrorWithContext(errWithContext, severity, errKey...)
}

func (c *ErrorCollector) addErrorWithContext(errWithContext ErrorWithContext, severity Severity, errKey ...string) {
	aggregationKey := getAggregationKey(errWithContext, errKey...)
	c.mux.Lock()
	defer c.mux.Unlock()
	if aggregatedErr, ok := c.aggregatedErrors[aggregationKey]; ok {
		aggregatedErr.addError(errWithContext)
	} else {
		aggregatedErr := newAggregatedError(aggregationKey, severity)
		aggregatedErr.addError(errWithContext)
		c.aggregatedErrors[aggregationKey] = &aggregatedErr
	}
}
