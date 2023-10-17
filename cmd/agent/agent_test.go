package main

import "testing"

func TestEndpoint(t *testing.T) {

	tests := []struct {
		name string
		cfg  config
		want string
	}{
		{
			name: "positive",
			cfg:  config{Addr: "localhost"},
			want: "http://localhost/updates",
		},
		{
			name: "positive 2",
			cfg:  config{Addr: ":8080"},
			want: "http://:8080/updates",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Endpoint(tt.cfg); got != tt.want {
				t.Errorf("Endpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
