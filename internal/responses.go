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

func NewResponse502() ResponseError {
	return ResponseError{
		Code:    502,
		Message: "Failed to fetch image from upstream server.",
	}
}

func NewResponse500() ResponseError {
	return ResponseError{
		Code:    500,
		Message: "Internal server error.",
	}
}
