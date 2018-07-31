package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJWTUserToken(t *testing.T) {

	composer := &JWTUserToken{
		Session:    "session",
		SID:        "SID",
		ID:         "ID",
		Username:   "Username",
		SigningKey: []byte("secret"),
	}

	assert.NotPanics(t, func() { composer.Create() })
}
