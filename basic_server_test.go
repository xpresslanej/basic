package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_helloWorldHandler(t *testing.T) {
	sr := strings.NewReader("{\"name\":\"Jeff\"}")
	req := httptest.NewRequest("GET", "http://localhost:8080/helloworld", sr)
	rw := httptest.NewRecorder()
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Jeff", args: args{rw, req}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testreq helloWorldRequest
			decoder := json.NewDecoder(tt.args.r.Body)
			err := decoder.Decode(&testreq)
			if err != nil {
				t.Error("Missing body")
			}
			if tt.name != testreq.Name {
				t.Errorf("Got %v, wanted %v ", testreq.Name, tt.name)
			}
			helloWorldHandler(tt.args.w, tt.args.r)
		})
	}
}
