package request

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/labstack/echo"
)

func MakeReq(method string, form url.Values, data bool, authorization string) (echo.Context, *httptest.ResponseRecorder) {

	e := echo.New()
	var req *http.Request
	req = httptest.NewRequest(method, "/route/path", strings.NewReader(form.Encode()))

	if data != true {
		req = httptest.NewRequest(method, "/route/path", nil)
	}

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	if authorization != "" {
		req.Header.Set(echo.HeaderAuthorization, authorization)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return c, rec
}
