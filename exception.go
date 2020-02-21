package main

import (
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Severity string

const (
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
	maxTraces                = 5
)

type ErrorAggregate struct {
	AggregationKey string             `json:"aggregation_key"`
	TotalCount     int                `json:"total_count"`
	Severity       string             `json:"severity"`
	LatestErrors   []ErrorWithContext `json:"latest_errors"`
}

type ErrorWithContext struct {
	Error       ErrorInstance `json:"error"`
	UUID        uuid.UUID     `json:"uuid"`
	Timestamp   time.Time     `json:"timestamp"`
	Severity    Severity      `json:"severity"`
	HTTPContext HTTPContext   `json:"http_context"`
}

type ErrorInstance struct {
	Class      string   `json:"class"`
	Message    string   `json:"message"`
	Stacktrace []string `json:"stacktrace"`
	Cause      string   `json:"cause"`
}

type HTTPContext struct {
	RequestMethod  string            `json:"request_method"`
	RequestURL     string            `json:"request_url"`
	RequestHeaders map[string]string `json:"request_headers"`
}

func NewErrorInstance(err error, stacktrace []string) ErrorInstance {
	return ErrorInstance{
		Cause:      err.Error(),
		Class:      err.Error(),
		Stacktrace: stacktrace,
	}
}

func NewErrorWithContext(errorInstance ErrorInstance, severity Severity, httpCtx HTTPContext) ErrorWithContext {
	return ErrorWithContext{
		Error:       errorInstance,
		UUID:        uuid.New(),
		Timestamp:   time.Now().UTC(),
		Severity:    severity,
		HTTPContext: httpCtx,
	}
}

func hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum32())
}

func (e *ErrorWithContext) aggregationKey() string {
	stacktraceHead := e.Error.Stacktrace
	if len(e.Error.Stacktrace) > maxTraces {
		stacktraceHead = stacktraceHead[:maxTraces]
	}
	stacktraceHeadHash := hash(strings.Join(stacktraceHead[:], ""))
	return fmt.Sprintf("%s@%s", e.Error.Class, stacktraceHeadHash)
}
