package internal

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	type testCase struct {
		path    string
		success bool
		code    int
	}

	tests := []testCase{
		{
			"/",
			false,
			404,
		},
		{
			"/bad/path",
			false,
			404,
		},
		{
			"/badpath/",
			false,
			404,
		},
		{
			"/get",
			false,
			400,
		},
		{
			"/get/",
			false,
			400,
		},
		{
			"/get/test",
			true,
			0,
		},
	}

	server := HttpServer{}

	go func() {
		server.StartServer()
	}()

	time.Sleep(time.Second)

	for _, test := range tests {
		resp, err := http.Get("http://localhost:8000" + test.path)
		if err != nil {
			t.Errorf("Request failed: %s", err)
		}
		defer resp.Body.Close()

		var data map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			t.Errorf("Error decoding response body: %s", err)
		}

		_, exists := data["Code"]
		if exists == test.success {
			t.Error("Path had incorrect success state")
		}

	}

}
