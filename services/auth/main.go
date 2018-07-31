package auth

import (
	"github.com/dgrijalva/jwt-go"
)

type (
	JWTUserToken struct {
		Session    string
		SID        string
		ID         string
		Username   string
		SigningKey []byte
	}
	JWTUserEmbed struct {
		ID       string `json:"userId"`
		Username string `json:"username"`
	}
	JWTClaims struct {
		Session string `json:"session"`
		SID     string `json:"sid"`
		JWTUserEmbed
		jwt.StandardClaims
	}
)

func (j *JWTUserToken) Create() string {

	var SignedToken string
	var err error

	claims := JWTClaims{
		j.Session,
		j.SID,
		JWTUserEmbed{
			j.ID,
			j.Username,
		},
		jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if SignedToken, err = token.SignedString(j.SigningKey); err != nil {
		panic(err)
	}

	return "Bearer " + SignedToken
}
