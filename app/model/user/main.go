package user

import (
	"errors"
	"strings"

	"github.com/Gommunity/GoWithWith/config/database"
	"github.com/Gommunity/GoWithWith/services/encrypt"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2/bson"
)

const UserCollection = "User"

type User struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`

	Username    string `json:"username" bson:"username"`
	Password    string `json:"password" bson:"password"`
	Email       string `json:"email" bson:"email"`
	VerifyEmail bool   `json:"verifyEmail" bson:"verifyEmail"`
}

func (u User) FindUserByCredentials() (User, error) {

	userModel := database.Connection.Model(UserCollection)
	user := &User{}
	var findStruct bson.M

	if strings.Index(u.Username, "@") > -1 {
		findStruct = bson.M{"email": strings.ToLower(u.Username)}
	} else {
		findStruct = bson.M{"username": strings.ToLower(u.Username)}
	}

	err := userModel.FindOne(findStruct).Exec(user)

	_, ok := err.(*mongodm.NotFoundError)
	switch {
	case ok:
		return User{}, errors.New("Your username has not been found")
	case err != nil:
		panic(err)
	default:
		{
			match := encrypt.CheckHash(u.Password, user.Password)
			if !match {
				return User{}, errors.New("Credentials are invalid")
			}
		}
	}

	return *user, nil
}

func (u User) ChangePassword() {

	userModel := database.Connection.Model(UserCollection)
	hash, _ := encrypt.Hash(u.Password)
	update := bson.M{
		"$set": bson.M{
			"password": hash,
		},
	}

	var err error
	if u.Email != "" {
		err = userModel.Update(bson.M{"email": strings.ToLower(u.Email)}, update)
	} else {
		err = userModel.UpdateId(bson.ObjectIdHex(u.Username), update)
	}

	if err != nil {
		panic(err)
	}
}

func (u User) Activation() {

	userModel := database.Connection.Model(UserCollection)
	update := bson.M{
		"$unset": bson.M{
			"verifyEmail": nil,
		},
	}

	err := userModel.Update(bson.M{"email": strings.ToLower(u.Email)}, update)

	if err != nil {
		panic(err)
	}
}

func (u User) CheckUsername(username string) (User, error) {

	userModel := database.Connection.Model(UserCollection)
	user := &User{}
	err := userModel.FindOne(bson.M{"username": strings.ToLower(username)}).Exec(user)

	_, ok := err.(*mongodm.NotFoundError)
	switch {
	case ok:
		return User{}, nil
	case err != nil:
		panic(err)
	default:
		return *user, errors.New("The username has already been used")
	}
}

func (u User) CheckEmail(email string) (User, error) {

	userModel := database.Connection.Model(UserCollection)
	user := &User{}
	err := userModel.FindOne(bson.M{"email": strings.ToLower(email)}).Exec(user)

	_, ok := err.(*mongodm.NotFoundError)
	switch {
	case ok:
		return User{}, nil
	case err != nil:
		panic(err)
	default:
		return *user, errors.New("The email address has already been used")
	}
}

func (u User) CreateUser() {

	userModel := database.Connection.Model(UserCollection)
	user := &User{}
	userModel.New(user)
	hash, _ := encrypt.Hash(u.Password)

	user.Username = strings.ToLower(u.Username)
	user.Password = hash
	user.Email = strings.ToLower(u.Email)
	user.VerifyEmail = true

	err := user.Save()

	if err != nil {
		panic(err)
	}
}
