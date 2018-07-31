package response

import (
	"github.com/labstack/echo"
)

type (
	Composer struct {
		Context echo.Context
	}
	Message struct {
		Message   string `json:"Message"`
		ErrorCode int    `json:"Code"`
	}
	Throw struct {
		error Message
	}
)

// JSON response
// @Summary throw response
// @Description show response with a description and internal error code
// @Tags response
// @Produce  json
func (c Composer) JSON(Code int, Text string) error {
	message := &Throw{Message{Message: Text, ErrorCode: Code}}
	return c.Context.JSON(Code, message.error)
}
