package main

import (
	"testing"

	"github.com/thefrol/kysh-kysh-meow/internal/collector"
	"github.com/thefrol/kysh-kysh-meow/internal/config"
)

func TestEndpoint(t *testing.T) {

	tests := []struct {
		name string
		cfg  config.Agent
		want string
	}{
		{
			name: "positive",
			cfg:  config.Agent{Addr: "localhost"},
			want: "http://localhost/updates",
		},
		{
			name: "positive 2",
			cfg:  config.Agent{Addr: ":8080"},
			want: "http://:8080/updates",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := collector.Endpoint(tt.cfg.Addr, updateRoute); got != tt.want {
				t.Errorf("Endpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
