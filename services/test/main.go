package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/Gommunity/GoWithWith/config/database"
	"github.com/Gommunity/GoWithWith/services/validation"
	"github.com/labstack/echo"
	validator "gopkg.in/go-playground/validator.v9"
)

func DBComposer(localPath string) database.Composer {

	db := &database.Composer{
		Locals:   localPath,
		Addrs:    []string{os.Getenv("DBAddrs")},
		Database: os.Getenv("DBName") + "_test",
		Username: os.Getenv("DBUsername"),
		Password: os.Getenv("DBPassword"),
		Source:   os.Getenv("DBSource") + "_test",
	}

	return *db
}

func MakeRequest(method, data string, haveData bool, authorization string) (echo.Context, *httptest.ResponseRecorder) {

	e := echo.New()
	e.Validator = &validation.DataValidator{ValidatorData: validator.New()}

	var req *http.Request
	req = httptest.NewRequest(method, "/", strings.NewReader(data))

	if haveData != true {
		req = httptest.NewRequest(method, "/", nil)
	}

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if authorization != "" {
		req.Header.Set(echo.HeaderAuthorization, authorization)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return c, rec
}
