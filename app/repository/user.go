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

func FindUserByCredentials(username, password string) (model.User, error) {

	getUser := database.Connection.Model(model.UserCollection)
	user := &model.User{}
	var findStruct bson.M

	if strings.Index(username, "@") > -1 {
		findStruct = bson.M{"email": strings.ToLower(username)}
	} else {
		findStruct = bson.M{"username": strings.ToLower(username)}
	}

	err := getUser.FindOne(findStruct).Exec(user)

	if _, ok := err.(*mongodm.NotFoundError); ok {
		return model.User{}, errors.New("Credentials are invalid")
	} else if err != nil {
		panic(err)
	} else {
		match := encrypt.CheckPasswordHash(password, user.Password)
		if !match {
			return model.User{}, errors.New("Credentials are invalid")
		}
	}

	return *user, nil
}

func ChangeUserPassword(username, password string) {

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

func CheckUsername(username string) (model.User, error) {

	getUser := database.Connection.Model(model.UserCollection)
	user := &model.User{}
	err := getUser.FindOne(bson.M{"username": strings.ToLower(username)}).Exec(user)

	if _, ok := err.(*mongodm.NotFoundError); ok {
		return model.User{}, nil
	} else if err != nil {
		panic(err)
	} else {
		errs := validation.Errors{}
		errs["username"] = errors.New("The username exist")
		return *user, errs
	}
}

func CheckEmail(email string) (model.User, error) {

	getUser := database.Connection.Model(model.UserCollection)
	user := &model.User{}
	err := getUser.FindOne(bson.M{"email": strings.ToLower(email)}).Exec(user)

	if _, ok := err.(*mongodm.NotFoundError); ok {
		return model.User{}, nil
	} else if err != nil {
		panic(err)
	} else {
		errs := validation.Errors{}
		errs["email"] = errors.New("The email exist")
		return *user, errs
	}
}

func CheckEmailVerify(email string) (model.User, error) {

	getUser := database.Connection.Model(model.UserCollection)
	user := &model.User{}
	err := getUser.FindOne(bson.M{"email": strings.ToLower(email), "verifyEmail": false}).Exec(user)

	if err, ok := err.(*mongodm.NotFoundError); ok {
		return model.User{}, err
	} else if err != nil {
		panic(err)
	} else {
		return model.User{}, nil
	}
}

func CreateUser(username, password, email string) {

	getUser := database.Connection.Model(model.UserCollection)
	user := &model.User{}
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
