package models

import (
	"errors"
	"time"

	db "../database"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2/bson"
)

type (
	// Session Model
	Session struct {
		mongodm.DocumentBase `json:",inline" bson:",inline"`

		IP        string `json:"ip" bson:"ip"`
		Key       string `json:"key" bson:"key"`
		UserID    string `json:"userId" bson:"userId"`
		UserAgent string `json:"userAgent" bson:"userAgent"`
	}
	// JWTUserEmbed for CreateJWToken func
	JWTUserEmbed struct {
		ID       string `json:"userId"`
		Username string `json:"username"`
	}
	// JWTClaims for CreateJWToken func
	JWTClaims struct {
		Session string `json:"session"`
		SID     string `json:"sid"`
		JWTUserEmbed
		jwt.StandardClaims
	}
)

// GetUserSessions ...
func GetUserSessions(userID string, page, limit int) (Pagination, error) {

	getModel := db.Connection.Model("Session")

	sessions := []*Session{}

	dbResult := getModel.Find(bson.M{"userId": userID}).Select(bson.M{"key": 0, "userId": 0}).Sort("updatedAt").
		Skip((page - 1) * limit).Limit(limit)

	count, _ := dbResult.Count()
	err := dbResult.Exec(&sessions)

	if _, ok := err.(*mongodm.NotFoundError); ok {
		return Pagination{}, errors.New("no record")
	} else if err != nil {
		panic(err)
	}

	pagination := GeneratePagination(sessions, count, page, limit)

	return pagination, nil
}

// CreateSession ...
func CreateSession(userID, ip, userAgent string) (string, string, error) {

	getSession := db.Connection.Model("Session")

	session := &Session{}
	getSession.New(session)

	sess := uuid.Must(uuid.NewV4()).String()
	hash, _ := HashPassword(sess)

	session.IP = ip
	session.Key = hash
	session.UserID = userID
	session.UserAgent = userAgent

	err := session.Save()

	if err != nil {
		panic(err)
	}

	return session.GetId().Hex(), sess, nil
}

// UpdateLastActive ...
func UpdateLastActive(SID string) error {

	getSession := db.Connection.Model("Session")

	update := bson.M{
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}

	err := getSession.UpdateId(bson.ObjectIdHex(SID), update)

	if err != nil {
		panic(err)
	}

	return nil
}

// SessionFindByCredentials ...
func SessionFindByCredentials(sess, SID string) error {

	session, err := SessionFindByID(SID)

	if err != nil {
		return err
	}

	match := CheckPasswordHash(sess, session.Key)
	if !match {
		return errors.New("Credentials are invalid")
	}

	return nil
}

// SessionFindByID ...
func SessionFindByID(SID string) (Session, error) {

	getSession := db.Connection.Model("Session")

	session := &Session{}

	err := getSession.FindId(bson.ObjectIdHex(SID)).Exec(session)

	if _, ok := err.(*mongodm.NotFoundError); ok {
		return Session{}, errors.New("Session not found")
	} else if err != nil {
		panic(err)
	}

	return *session, nil
}

// CreateJWToken ...
func CreateJWToken(session, SID, username, userID string, signingKey []byte) (string, error) {

	var SignedToken string
	var err error

	claims := JWTClaims{
		session,
		SID,
		JWTUserEmbed{
			userID,
			username,
		},
		jwt.StandardClaims{},
	}

	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if SignedToken, err = tokenJWT.SignedString(signingKey); err != nil {
		panic(err)
	}

	SignedToken = "Bearer " + SignedToken

	return SignedToken, nil
}

// DeleteSession ...
func DeleteSession(id string) error {

	getSession := db.Connection.Model("Session")
	err := getSession.RemoveId(bson.ObjectIdHex(id))

	if err != nil {
		panic(err)
	}

	return nil
}
