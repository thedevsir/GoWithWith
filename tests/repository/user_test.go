package repository

import (
	"testing"

	"github.com/Gommunity/GoWithWith/config/database"
	"github.com/Gommunity/GoWithWith/services/encrypt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/zebresel-com/mongodm"
)

var dbUser *mongodm.Model

func userBeforeTest() {

	err := godotenv.Load("../.env")
	if err != nil {
		panic(":main:init: ErrorLoading.EnvFile")
	}

	Models := make(map[string]mongodm.IDocumentBase)
	Models["user"] = &User{}

	database.Initial(Models, true)

	dbUser = database.Connection.Model("user")
	dbUser.RemoveAll(nil)
}

func userAfterTest() {
	dbUser.RemoveAll(nil)
}

func TestBcrypt(t *testing.T) {

	var str string
	var err error

	password := "12345"

	t.Run("HashPassword", func(t *testing.T) {
		str, err = encrypt.HashPassword(password)
		assert.IsType(t, String, str)
		assert.Nil(t, err)
	})

	t.Run("CheckPasswordHash", func(t *testing.T) {
		err := encrypt.CheckPasswordHash(password, str)
		assert.True(t, err)
	})
}

func TestUser(t *testing.T) {

	userBeforeTest()

	person := &User{
		Username: "Irani",
		Password: "12345",
		Email:    "freshmanlimited@gmail.com",
	}

	CreateUser(person.Username, person.Password, person.Email)

	t.Run("CheckEmail", func(t *testing.T) {
		ce, err := CheckEmail(person.Email)
		ce1, err1 := CheckEmail("fake@service.domain")
		assert.Error(t, err)
		assert.IsType(t, User{}, ce)
		assert.IsType(t, User{}, ce1)
		assert.Nil(t, err1)
	})

	// t.Run("CheckEmailVerify", func(t *testing.T) {
	// 	_, err := CheckEmailVerify(person.Email)
	// 	_, err1 := CheckEmail("fake@aaservice.domain")
	// 	assert.Error(t, err1)
	// 	assert.Nil(t, err)
	// })

	t.Run("CheckUsername", func(t *testing.T) {
		cu, err := CheckUsername(person.Username)
		cu1, err1 := CheckUsername("fakeuser")
		assert.Error(t, err)
		assert.IsType(t, User{}, cu)
		assert.IsType(t, User{}, cu1)
		assert.Nil(t, err1)
	})

	t.Run("FindUserByCredentials", func(t *testing.T) {
		user, err := FindUserByCredentials(person.Username, person.Password)
		user1, err1 := FindUserByCredentials("fakeuser", person.Password)
		assert.IsType(t, User{}, user)
		assert.Nil(t, err)
		assert.IsType(t, User{}, user1)
		assert.Error(t, err1)
	})

	ChangeUserPassword(person.Username, "newPassword")

	userAfterTest()
}
