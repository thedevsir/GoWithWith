package structs

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// SignupStruct | /user/signup
type SignupStruct struct {
	Username string `form:"username"`
	Password string `form:"password"`
	Email    string `form:"email"`
}

// Joi | /user/signup
func (s SignupStruct) Joi() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Username, validation.Required, validation.Length(3, 50), is.Alphanumeric),
		validation.Field(&s.Password, validation.Required, validation.Length(3, 50)),
		validation.Field(&s.Email, validation.Required, is.Email),
	)
}

// LoginStruct | /user/login
type LoginStruct struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

// Authorization | /user/login
type Authorization struct {
	Authorization string `json:"authorization"`
}

// Joi | /user/login
func (l LoginStruct) Joi() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Username, validation.Required),
		validation.Field(&l.Password, validation.Required),
	)
}

// LogoutStruct | /user/auth/logout
type LogoutStruct struct {
	ID string `form:"id"`
}

// Joi | /user/auth/logout
func (l LogoutStruct) Joi() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.ID, validation.Required),
	)
}

// ForgotStruct | /user/forgot
type ForgotStruct struct {
	Email string `form:"email"`
}

// Joi | /user/forgot
func (f ForgotStruct) Joi() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.Email, validation.Required, is.Email),
	)
}

// ResetStruct | /user/reset
type ResetStruct struct {
	Token    string `form:"token"`
	Password string `form:"password"`
}

// Joi | /user/reset
func (r ResetStruct) Joi() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Token, validation.Required),
		validation.Field(&r.Password, validation.Required, validation.Length(3, 50)),
	)
}
