package main

import (
	"testing"

	"github.com/go-errors/errors"
)

var aggregationKeyCases = []struct {
	expectedAggregationKey string
	stacktrace             []string
}{
	{"testingError@811c9dc5", []string{""}},
	{"testingError@8de5e669", []string{"line 0:", "division by zero"}},
	{"testingError@9d610c3f", []string{"line 0:", "division by zero", "line 1:", "test()", "line 4:", "checkTest()"}},
	{"testingError@9b5eca82", []string{"line 0:", "division by zero", "line 1:", "test()", "line 5:", "checkTest()"}},
}

func TestException_aggregationKey(t *testing.T) {
	for _, tt := range aggregationKeyCases {
		t.Run(tt.expectedAggregationKey, func(t *testing.T) {
			errorInstance := NewErrorInstance(errors.New("testingError"), tt.stacktrace)
			errorInstanceWithContext := NewErrorWithContext(errorInstance, SeverityError, HTTPContext{})
			resultAggregationKey := errorInstanceWithContext.aggregationKey()
			if resultAggregationKey != tt.expectedAggregationKey {
				t.Errorf("error in aggregationKey, expected: %s, got %s", tt.expectedAggregationKey, resultAggregationKey)
			}
		})
	}
}
