package jwt

import (
	"fmt"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

func ParseJWT(token string, secret []byte) (*jwt.Token, error) {

	var err error
	token = strings.Replace(token, "Bearer ", "", 1)
	tokenParsed := new(jwt.Token)
	tokenParsed, err = jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Token Parse Error")
		}
		return secret, nil
	})

	if err == nil && tokenParsed.Valid {
		return tokenParsed, nil
	}

	return &jwt.Token{}, err
}
