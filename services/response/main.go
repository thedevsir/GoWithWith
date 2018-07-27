package response

import (
	"net/http"

	"github.com/labstack/echo"
)

type errorMessage struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"code"`
}

type errorResponse struct {
	error errorMessage
}

//Created response
// @Summary response created
// @Description for method POST, to show created data id
// @Tags response
// @Produce  json
func Created(e echo.Context, data string) error {
	return e.JSON(http.StatusCreated, map[string]string{"data": data})
}

//Ok response
// @Summary response ok
// @Description for method PUT, PATCH, to show successfully response message
// @Tags response
// @Produce  json
func Ok(e echo.Context, message string) error {
	return e.JSON(http.StatusOK, map[string]string{"message": message})
}

//Data response
// @Summary response data
// @Description for method GET, to show response data
// @Tags response
// @Produce  json
func Data(e echo.Context, data interface{}) error {
	return e.JSON(http.StatusOK, data)
}

//Error response
// @Summary response error
// @Description show error response with a description and internal error code
// @Tags response
// @Produce  json
func Error(text string, code int) error {
	response := &errorResponse{errorMessage{Message: text, ErrorCode: code}}
	return &echo.HTTPError{Code: http.StatusUnprocessableEntity, Message: response.error}
}
