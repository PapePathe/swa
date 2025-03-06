package server_test

import (
	"net/http"
	"swahili/lang/server"
	"testing"
)

func TestLiveInterpreter_ServeHTTP(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		w http.ResponseWriter
		r *http.Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var mux server.LiveInterpreter

			mux.ServeHTTP(tt.w, tt.r)
		})
	}
}
