package repository

import (
	"errors"
	"strings"

	"github.com/Gommunity/GoWithWith/app/model"
	"github.com/Gommunity/GoWithWith/config/database"
	"github.com/Gommunity/GoWithWith/services/encrypt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2/bson"
)

type userModel struct {
	*model.User
}

func (u userModel) FindUserByCredentials(username, password string) (*userModel, error) {

	getUser := database.Connection.Model(model.UserCollection)
	user := &u
	var findStruct bson.M

	if strings.Index(username, "@") > -1 {
		findStruct = bson.M{"email": strings.ToLower(username)}
	} else {
		findStruct = bson.M{"username": strings.ToLower(username)}
	}

	err := getUser.FindOne(findStruct).Exec(user)

	if _, ok := err.(*mongodm.NotFoundError); ok {
		return user, errors.New("Credentials are invalid")
	} else if err != nil {
		panic(err)
	} else {
		match := encrypt.CheckPasswordHash(password, user.Password)
		if !match {
			return user, errors.New("Credentials are invalid")
		}
	}

	return user, nil
}

func (u userModel) ChangeUserPassword(username, password string) {

	getUser := database.Connection.Model(model.UserCollection)
	hash, _ := encrypt.HashPassword(password)
	update := bson.M{
		"$set": bson.M{
			"password": hash,
		},
	}

	err := getUser.Update(bson.M{"username": strings.ToLower(username)}, update)

	if err != nil {
		panic(err)
	}
}

func (u userModel) CheckUsername(username string) (*userModel, error) {

	getUser := database.Connection.Model(model.UserCollection)
	user := &u
	err := getUser.FindOne(bson.M{"username": strings.ToLower(username)}).Exec(user)

	if _, ok := err.(*mongodm.NotFoundError); ok {
		return user, nil
	} else if err != nil {
		panic(err)
	} else {
		errs := validation.Errors{}
		errs["username"] = errors.New("The username exist")
		return user, errs
	}
}

func (u userModel) CheckEmail(email string) (*userModel, error) {

	getUser := database.Connection.Model(model.UserCollection)
	user := &u
	err := getUser.FindOne(bson.M{"email": strings.ToLower(email)}).Exec(user)

	if _, ok := err.(*mongodm.NotFoundError); ok {
		return user, nil
	} else if err != nil {
		panic(err)
	} else {
		errs := validation.Errors{}
		errs["email"] = errors.New("The email exist")
		return user, errs
	}
}

func (u userModel) CheckEmailVerify(email string) (*userModel, error) {

	getUser := database.Connection.Model(model.UserCollection)
	user := &u
	err := getUser.FindOne(bson.M{"email": strings.ToLower(email), "verifyEmail": false}).Exec(user)

	if err, ok := err.(*mongodm.NotFoundError); ok {
		return user, err
	} else if err != nil {
		panic(err)
	} else {
		return user, nil
	}
}

func (u userModel) CreateUser(username, password, email string) {

	getUser := database.Connection.Model(model.UserCollection)
	user := &u
	getUser.New(user)
	hash, _ := encrypt.HashPassword(password)

	user.Username = strings.ToLower(username)
	user.Password = hash
	user.Email = strings.ToLower(email)
	user.VerifyEmail = false

	err := user.Save()

	if err != nil {
		panic(err)
	}
}
