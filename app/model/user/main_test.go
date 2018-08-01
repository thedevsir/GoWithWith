package user

import (
	"testing"

	"github.com/Gommunity/GoWithWith/config/database"
	"github.com/Gommunity/GoWithWith/services/test"
	"github.com/Gommunity/GoWithWith/services/utility"
	"github.com/stretchr/testify/assert"
	"github.com/zebresel-com/mongodm"
)

var userCollection *mongodm.Model

func userBeforeTest() {

	utility.LoadEnvironmentVariables("../../../.env")

	db := test.DBComposer("../../../resource/locals/locals.json")
	db.Shoot(map[string]mongodm.IDocumentBase{
		"users": &User{},
	})

	userCollection = database.Connection.Model(UserCollection)
	userCollection.RemoveAll(nil)
}

func userAfterTest() {
	userCollection.RemoveAll(nil)
}

func TestUser(t *testing.T) {

	userBeforeTest()

	user := &User{
		Username: "Irani",
		Password: "12345",
		Email:    "freshmanlimited@gmail.com",
	}

	t.Run("CreateUser", func(t *testing.T) {
		assert.NotPanics(t, func() { user.CreateUser() })
	})

	t.Run("Activation", func(t *testing.T) {
		assert.NotPanics(t, func() { user.Activation() })
	})

	t.Run("CheckEmail", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			result, err := user.CheckEmail("fake@service.domain")
			assert.Nil(t, err)
			assert.IsType(t, User{}, result)
		})

		t.Run("NotFound", func(t *testing.T) {
			result, err := user.CheckEmail(user.Email)
			assert.Error(t, err)
			assert.IsType(t, User{}, result)
		})
	})

	t.Run("CheckUsername", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {
			result, err := user.CheckUsername("fakeuser")
			assert.IsType(t, User{}, result)
			assert.Nil(t, err)
		})

		t.Run("NotFound", func(t *testing.T) {
			result, err := user.CheckUsername(user.Username)
			assert.Error(t, err)
			assert.IsType(t, User{}, result)
		})
	})

	t.Run("ChangePassword", func(t *testing.T) {

		user := &User{
			Password: "12345678",
			Email:    "freshmanlimited@gmail.com",
		}

		assert.NotPanics(t, func() { user.ChangePassword() })
	})

	t.Run("FindUserByCredentials", func(t *testing.T) {

		t.Run("Success", func(t *testing.T) {

			user := &User{
				Username: "Irani",
				Password: "12345678",
			}

			result, err := user.FindUserByCredentials()
			if assert.NoError(t, err) {
				assert.IsType(t, User{}, result)
				assert.Equal(t, "irani", result.Username)
				assert.Equal(t, "freshmanlimited@gmail.com", result.Email)
			}
		})

		t.Run("UsernameNotFound", func(t *testing.T) {

			user := &User{
				Username: "Irani@gmail.com",
				Password: "12345",
			}
			_, err := user.FindUserByCredentials()
			if assert.Error(t, err) {
				assert.Equal(t, "Your username has not been found", err.Error())
			}
		})

		t.Run("PasswordHasNotMatch", func(t *testing.T) {

			user := &User{
				Username: "Irani",
				Password: "123",
			}
			_, err := user.FindUserByCredentials()
			if assert.Error(t, err) {
				assert.Equal(t, "Credentials are invalid", err.Error())
			}
		})
	})

	userAfterTest()
}
