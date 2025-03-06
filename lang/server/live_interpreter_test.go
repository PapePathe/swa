package server_test

import (
	"net/http"
	"swahili/lang/server"
	"testing"
)

func TestLiveInterpreter_ServeHTTP(t *testing.T) {
	tests := []struct {
		name string
		w    http.ResponseWriter
		r    *http.Request
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mux server.LiveInterpreter

			mux.ServeHTTP(tt.w, tt.r)
		})
	}
}
