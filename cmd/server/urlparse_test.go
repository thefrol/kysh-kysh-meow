package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlParsing(t *testing.T) {
	type metric struct {
		t     string //type
		name  string
		value string
	}

	tests := []struct {
		name        string
		url         string
		output      metric
		raisesError bool
	}{
		{
			name:        "counter #1",
			url:         "/update/counter/test1/20111",
			output:      metric{"counter", "test1", "20111"},
			raisesError: false,
		},
		{
			name:        "gauge #1",
			url:         "/update/gauge/test1/2.11",
			output:      metric{"gauge", "test1", "2.11"},
			raisesError: false,
		},
		{
			name:        "gauge #2",
			url:         "/update/gauge/test1/2.11e6",
			output:      metric{"gauge", "test1", "2.11e6"},
			raisesError: false,
		},
		{
			name:        "negative #1",
			url:         "/gauge/test1/2.11",
			output:      metric{}, //fail, needs memory
			raisesError: true,
		},
		{
			name:        "negative #2",
			url:         "/update/counter/unknown/22/33",
			output:      metric{},
			raisesError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := ParseURL(tt.url)
			if tt.raisesError {
				assert.NotNil(t, err)
				return
			}
			assert.Equal(t, tt.output.t, value.Type())
			assert.Equal(t, tt.output.name, value.Name())
			assert.Equal(t, tt.output.value, value.Value())
		})
	}
}
