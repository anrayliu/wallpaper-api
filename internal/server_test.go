package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

// allows loading of env file before testing

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	server := HttpServer{}

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

	content_type := resp.Header.Get("Content-Type")
	len := resp.Header.Get("Content-Length")
	int_len, err := strconv.Atoi(len)
	if err != nil || int_len != 2056737 || content_type != "image/jpeg" {
		t.Error("Bad response")
	}

}
