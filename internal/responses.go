package internal

// Google's JSON style guide

type ResponseError struct {
	Code    int
	Message string
}

type ResponseSuccess struct {
	TestData string
}

func NewResponse404() ResponseError {
	return ResponseError{
		Code:    404,
		Message: "Page not found.",
	}
}

func NewResponse400() ResponseError {
	return ResponseError{
		Code:    400,
		Message: "Target file not specified.",
	}
}
