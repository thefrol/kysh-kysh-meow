package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJob_Elapsed(t *testing.T) {
	type fields struct {
		lastCallAgo time.Duration
		interval    time.Duration
		//function func()
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "positive #1",
			fields: fields{lastCallAgo: 3 * time.Second, interval: 10 * time.Second},
			want:   false,
		},
		{
			name:   "negative #1",
			fields: fields{lastCallAgo: 1 * time.Second, interval: 2 * time.Second},
			want:   false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				lastCall: time.Now().Add(-tt.fields.lastCallAgo),
				interval: tt.fields.interval,
			}
			assert.Equal(t, tt.want, j.Elapsed())
		})
	}
}
