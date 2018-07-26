package controller

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Gommunity/GoWithWith/config/mail"
	"github.com/Gommunity/GoWithWith/resource/views/mail/viewmail"

	"github.com/Gommunity/GoWithWith/app/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gopkg.in/gomail.v2"
)

var (
	ModeliCheckUsername         = repository.CheckUsername
	ModeliCheckEmail            = repository.CheckEmail
	ModeliCreateUser            = repository.CreateUser
	ModeliAbuseDetected         repository.AbuseDetected
	ModeliAbuseDetectedCheck    = ModeliAbuseDetected.Check
	ModeliFindUserByCredentials = repository.FindUserByCredentials
	ModeliCreateSession         = repository.CreateSession
	ModeliCreateJWToken         = repository.CreateJWToken
	ModeliDeleteSession         = repository.DeleteSession
	ModeliSessionFindByID       = repository.SessionFindByID
	ModeliGetUserSessions       = repository.GetUserSessions
	ModeliChangeUserPassword    = repository.ChangeUserPassword
	ModeliCheckEmailVerify      = repository.CheckEmailVerify
	ModeliMakeEmailToken        = MakeEmailToken
	ModeliSendResetMail         = SendResetMail
	ModeliSendVerficationMail   = SendVerficationMail
	ModeliParseJWT              = ParseJWT
)

type JoiError struct {
	Code    int
	Message map[string]string
}

type JoiString struct {
	Code    int
	Message string
}

type EmailToken struct {
	Action   string
	Username string
	Email    string
	jwt.StandardClaims
}

func InitConfig() {
	ModeliAbuseDetected = repository.AbuseDetected{
		MaxIP:            os.Getenv("AuthAttemptsForIp"),
		MaxIPAndUsername: os.Getenv("AuthAttemptsForIpAndUser"),
	}
}

func MakeEmailToken(action, username, email string, secret []byte) string {

	var token string
	var err error

	claims := EmailToken{
		action,
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

	return token
}

func SendVerficationMail(username, email, token string) {

	verifyBody := viewmail.Verify{
		Username:     username,
		EmailAddress: email,
		Token:        token,
	}
	emailBody, emailText := viewmail.GenerateTemplate(verifyBody.Email())

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EmailFrom"))
	m.SetHeader("To", verifyBody.EmailAddress)
	m.SetHeader("Subject", "Confirm your account")
	m.SetBody("text/plain", emailText)
	m.AddAlternative("text/html", emailBody)

	if err := mail.Driver.DialAndSend(m); err != nil {
		panic(err)
	}
}

func SendResetMail(username, email, token string) {

	forgotBody := viewmail.Forgot{
		Username:     username,
		EmailAddress: email,
		Token:        token,
	}
	emailBody, emailText := viewmail.GenerateTemplate(forgotBody.Email())

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EmailFrom"))
	m.SetHeader("To", forgotBody.EmailAddress)
	m.SetHeader("Subject", "Reset your password")
	m.SetBody("text/plain", emailText)
	m.AddAlternative("text/html", emailBody)

	if err := mail.Driver.DialAndSend(m); err != nil {
		panic(err)
	}
}

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
