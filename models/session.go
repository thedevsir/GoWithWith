package models

import (
	"errors"
	"time"

	"github.com/Gommunity/GoWithWith/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2/bson"
)

type (
	Session struct {
		mongodm.DocumentBase `json:",inline" bson:",inline"`

		IP        string `json:"ip" bson:"ip"`
		Key       string `json:"key" bson:"key"`
		UserID    string `json:"userId" bson:"userId"`
		UserAgent string `json:"userAgent" bson:"userAgent"`
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

func GetUserSessions(userID string, page, limit int) (Pagination, error) {

	getModel := database.Connection.Model("Session")
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

func CreateSession(userID, ip, userAgent string) (string, string) {

	getSession := database.Connection.Model("Session")
	session := &Session{}
	getSession.New(session)
	sess := uuid.Must(uuid.NewV4(), nil).String() // nil
	hash, _ := HashPassword(sess)

	session.IP = ip
	session.Key = hash
	session.UserID = userID
	session.UserAgent = userAgent

	err := session.Save()
	if err != nil {
		panic(err)
	}

	return session.GetId().Hex(), sess
}

func UpdateLastActive(SID string) error {

	getSession := database.Connection.Model("Session")
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

func SessionFindByID(SID string) (Session, error) {

	getSession := database.Connection.Model("Session")
	session := &Session{}

	err := getSession.FindId(bson.ObjectIdHex(SID)).Exec(session)
	if _, ok := err.(*mongodm.NotFoundError); ok {
		return Session{}, errors.New("Session not found")
	} else if err != nil {
		panic(err)
	}

	return *session, nil
}

func CreateJWToken(session, SID, username, userID string, signingKey []byte) string {

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
	return SignedToken
}

func DeleteSession(id string) {

	getSession := database.Connection.Model("Session")
	err := getSession.RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		panic(err)
	}
}
