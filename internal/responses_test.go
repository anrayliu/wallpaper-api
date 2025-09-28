package internal

import (
	"testing"
)

func TestResponseError(t *testing.T) {
	type testCase struct {
		callback func() ResponseError
		code     int
		message  string
	}

	tests := []testCase{
		{
			NewResponse400,
			400,
			"Target file not specified.",
		},
		{
			NewResponse404,
			404,
			"Page not found.",
		},
	}

	for _, test := range tests {
		response := test.callback()
		if response.Code != test.code {
			t.Errorf("Expected code %d, found %d.", test.code, response.Code)
		}
		if response.Message != test.message {
			t.Errorf("Expected message \"%s\", found \"%s\".", test.message, response.Message)
		}
	}
}
