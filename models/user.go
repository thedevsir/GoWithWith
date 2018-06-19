package models

import (
	"errors"
	"strings"

	db "../database"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/zebresel-com/mongodm"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

// User Model
type User struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`

	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Email    string `json:"email" bson:"email"`
}

// HashPassword with bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash ...
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// FindUserByCredentials ...
func FindUserByCredentials(username, password string) (User, error) {

	getUser := db.Connection.Model("User")

	user := &User{}

	var findStruct bson.M

	if strings.Index(username, "@") > -1 {
		findStruct = bson.M{"email": strings.ToLower(username)}
	} else {
		findStruct = bson.M{"username": strings.ToLower(username)}
	}

	err := getUser.FindOne(findStruct).Exec(user)

	if _, ok := err.(*mongodm.NotFoundError); ok {
		return User{}, errors.New("Credentials are invalid")
	} else if err != nil {
		panic(err)
	} else {
		match := CheckPasswordHash(password, user.Password)
		if !match {
			return User{}, errors.New("Credentials are invalid")
		}
	}

	return *user, nil
}

// ChangeUserPassword ...
func ChangeUserPassword(username, password string) error {

	getUser := db.Connection.Model("User")

	hash, _ := HashPassword(password)

	update := bson.M{
		"$set": bson.M{
			"password": hash,
		},
	}

	err := getUser.Update(bson.M{"username": strings.ToLower(username)}, update)

	if err != nil {
		panic(err)
	}

	return nil
}

// CheckUsername ...
func CheckUsername(username string) (User, error) {

	getUser := db.Connection.Model("User")

	user := &User{}
	err := getUser.FindOne(bson.M{"username": strings.ToLower(username)}).Exec(user)

	if _, ok := err.(*mongodm.NotFoundError); ok {
		return User{}, nil
	} else if err != nil {
		panic(err)
	} else {
		errs := validation.Errors{}
		errs["username"] = errors.New("The username exist")
		return *user, errs
	}

	return User{}, nil
}

// CheckEmail ...
func CheckEmail(email string) (User, error) {

	getUser := db.Connection.Model("User")

	user := &User{}
	err := getUser.FindOne(bson.M{"email": strings.ToLower(email)}).Exec(user)

	if _, ok := err.(*mongodm.NotFoundError); ok {
		return User{}, nil
	} else if err != nil {
		panic(err)
	} else {
		errs := validation.Errors{}
		errs["email"] = errors.New("The email exist")
		return *user, errs
	}

	return User{}, nil
}

// CreateUser ...
func CreateUser(username, password, email string) error {

	getUser := db.Connection.Model("User")

	user := &User{}
	getUser.New(user)

	hash, _ := HashPassword(password)

	user.Username = strings.ToLower(username)
	user.Password = hash
	user.Email = strings.ToLower(email)

	err := user.Save()

	if err != nil {
		panic(err)
	}

	return nil
}
