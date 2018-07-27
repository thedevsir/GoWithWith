package model

import (
	"github.com/zebresel-com/mongodm"
)

const SessionCollection = "Session"

type SignupStruct struct {
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Password string `json:"password validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
}

type LogoutStruct struct {
	ID string `json:"id" validate:"required"`
}

type LoginStruct struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Authorization struct {
	Authorization string `json:"authorization"`
}

type ResetStruct struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=3,max=50"`
}

type Session struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`

	IP        string `json:"ip" bson:"ip"`
	Key       string `json:"key" bson:"key"`
	UserID    string `json:"userId" bson:"userId"`
	UserAgent string `json:"userAgent" bson:"userAgent"`
}
