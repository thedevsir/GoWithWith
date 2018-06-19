package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	mail "../../gomail"
	models "../../models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	gomail "gopkg.in/gomail.v2"
)

// Global ...
var (
	ModeliCheckUsername         = models.CheckUsername
	ModeliCheckEmail            = models.CheckEmail
	ModeliCreateUser            = models.CreateUser
	ModeliAbuseDetected         models.AbuseDetected
	ModeliAbuseDetectedCheck    = ModeliAbuseDetected.Check
	ModeliFindUserByCredentials = models.FindUserByCredentials
	ModeliCreateSession         = models.CreateSession
	ModeliCreateJWToken         = models.CreateJWToken
	ModeliDeleteSession         = models.DeleteSession
	ModeliSessionFindByID       = models.SessionFindByID
	ModeliGetUserSessions       = models.GetUserSessions
	ModeliChangeUserPassword    = models.ChangeUserPassword
	ModeliMakeEmailToken        = MakeEmailToken
	ModeliSendResetMail         = SendResetMail
	ModeliParseJWT              = ParseJWT
)

// JoiError ...
type JoiError struct {
	Code    int
	Message map[string]string
}

// JoiString ...
type JoiString struct {
	Code    int
	Message string
}

// EmailToken ...
type EmailToken struct {
	Username string
	Email    string
	jwt.StandardClaims
}

// InitConfig ...
func InitConfig() {
	ModeliAbuseDetected = models.AbuseDetected{
		MaxIP:            os.Getenv("AuthAttemptsForIp"),
		MaxIPAndUsername: os.Getenv("AuthAttemptsForIpAndUser"),
	}
}

// MakeEmailToken ...
func MakeEmailToken(username, email string, secret []byte) (string, error) {

	var token string
	var err error

	claims := EmailToken{
		username,
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if token, err = tokenJWT.SignedString(secret); err != nil {
		panic(err)
	}

	return token, nil
}

// SendResetMail ...
func SendResetMail(username, email, token string) error {

	forgotBody := mail.Forgot{
		Username:     username,
		EmailAddress: email,
		Token:        token,
	}

	emailBody, emailText := mail.GenerateTemplate(forgotBody.Email())

	// Send mail
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EmailFrom"))
	m.SetHeader("To", forgotBody.EmailAddress)
	m.SetHeader("Subject", "Reset your password")
	m.SetBody("text/plain", emailText)
	m.AddAlternative("text/html", emailBody)

	// Send the email
	if err := mail.Driver.DialAndSend(m); err != nil {
		panic(err)
	}

	return nil
}

// MakeReq ...
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

// ParseJWT ...
func ParseJWT(token string, secret []byte) (*jwt.Token, error) {

	var err error
	token = strings.Replace(token, "Bearer ", "", 1)
	tokenParsed := new(jwt.Token)
	tokenParsed, err = jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return secret, nil
	})

	if err == nil && tokenParsed.Valid {
		return tokenParsed, nil
	}

	return &jwt.Token{}, err
}

// PaginationSettings ...
func PaginationSettings(c echo.Context) (int, int) {

	var page, limit = 1, 10

	queryPage := c.QueryParam("page")
	queryLimit := c.QueryParam("limit")

	if queryPage != "" {
		p, err := strconv.Atoi(queryPage)
		if err == nil {
			page = p
		}
	}

	if queryLimit != "" {
		l, err := strconv.Atoi(queryLimit)
		if err == nil {
			limit = l
		}
	}

	return page, limit
}
