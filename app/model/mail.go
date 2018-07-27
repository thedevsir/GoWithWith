package model

import jwt "github.com/dgrijalva/jwt-go"

type ResendEmailStruct struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotStruct struct {
	Email string `json:"email" validate:"required,email"`
}

type EmailToken struct {
	Action   string
	Username string
	Email    string
	jwt.StandardClaims
}
