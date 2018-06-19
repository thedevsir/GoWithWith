package helpers

// JoiError ...
type JoiError struct {
	Code    int   `json:"Code"`
	Message error `json:"Message"`
}

// JoiString ...
type JoiString struct {
	Code    int    `json:"Code"`
	Message string `json:"Message"`
}

// Throw Error ...
func Throw(err error) interface{} {

	return &JoiError{
		Code:    400,
		Message: err,
	}
}

// ThrowString Error ...
func ThrowString(err error) interface{} {

	return &JoiString{
		Code:    400,
		Message: err.Error(),
	}
}

// SayOk String ...
func SayOk(message string) interface{} {

	return &JoiString{
		Code:    200,
		Message: message,
	}
}
