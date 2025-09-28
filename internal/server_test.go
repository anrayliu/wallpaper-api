package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

// allows loading of env file before testing

func TestMain(m *testing.M) {
	server := HTTPServer{}

	go func() {
		server.StartServer()
	}()

	time.Sleep(time.Second)

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestBadRequests(t *testing.T) {
	type testCase struct {
		path string
		code int
	}

	tests := []testCase{
		{
			"/",
			404,
		},
		{
			"/bad/path",
			404,
		},
		{
			"/badpath/",
			404,
		},
		{
			"/get",
			400,
		},
		{
			"/get/",
			400,
		},
		{
			"/get/does-not-exist.jpg",
			404,
		},
	}

	for _, test := range tests {
		resp, err := http.Get("http://localhost:8000" + test.path)
		if err != nil {
			t.Errorf("Request failed: %s", err)
		}
		defer resp.Body.Close()

		var data map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			t.Errorf("Error decoding response body: %s", err)
		}

		fmt.Println(data)

		_, exists := data["Code"]
		if !exists {
			t.Error("Did not receive an error.")
		}

	}

}

func TestGoodRequest(t *testing.T) {
	resp, err := http.Get("http://localhost:8000/get/moon.jpg")
	if err != nil {
		t.Errorf("Request failed: %s", err)
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	len := resp.Header.Get("Content-Length")
	intLen, err := strconv.Atoi(len)
	if err != nil || intLen != 2056737 || contentType != "image/jpeg" {
		t.Error("Bad response")
	}

}
