package helpers

type JoiError struct {
	Code    int   `json:"Code"`
	Message error `json:"Message"`
}

type JoiString struct {
	Code    int    `json:"Code"`
	Message string `json:"Message"`
}

func Throw(err error) interface{} {
	return &JoiError{
		Code:    400,
		Message: err,
	}
}

func ThrowString(err error) interface{} {
	return &JoiString{
		Code:    400,
		Message: err.Error(),
	}
}

func SayOk(message string) interface{} {
	return &JoiString{
		Code:    200,
		Message: message,
	}
}
