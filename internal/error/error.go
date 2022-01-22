package error

import "errors"

type MyError struct {
	Cause      error
	Message    string
	StatusCode int
}

var JsonMarshalError = MyError{
	Cause:      errors.New("encode response error"),
	Message:    "Unable to convert an object to application/json type",
	StatusCode: 500,
}

var InvalidPara = MyError{
	Cause:      errors.New("invalid parameter"),
	Message:    "The Parameter is invalid.",
	StatusCode: 400,
}

func NewError(err error, msg string, code int) *MyError {
	return &MyError{
		Cause:      err,
		Message:    msg,
		StatusCode: code,
	}
}
