package repository

import (
	"errors"
	"time"

	"github.com/Gommunity/GoWithWith/app/model"
	"github.com/Gommunity/GoWithWith/config/database"
	"github.com/Gommunity/GoWithWith/services/encrypt"
	"github.com/Gommunity/GoWithWith/services/paginate"
	"github.com/satori/go.uuid"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2/bson"
)

func GetUserSessions(userID string, page, limit int) (paginate.Pagination, error) {

	getModel := database.Connection.Model(model.SessionCollection)
	sessions := []*model.Session{}
	dbResult := getModel.Find(bson.M{"userId": userID}).Select(bson.M{"key": 0, "userId": 0}).Sort("updatedAt").
		Skip((page - 1) * limit).Limit(limit)
	count, _ := dbResult.Count()
	err := dbResult.Exec(&sessions)

	if _, ok := err.(*mongodm.NotFoundError); ok {
		return paginate.Pagination{}, errors.New("no record")
	} else if err != nil {
		panic(err)
	}

	pagination := paginate.GeneratePagination(sessions, count, page, limit)

	return pagination, nil
}

func CreateSession(userID, ip, userAgent string) (string, string) {

	getSession := database.Connection.Model(model.SessionCollection)
	session := &model.Session{}
	getSession.New(session)
	sess := uuid.Must(uuid.NewV4(), nil).String() // nil
	hash, _ := encrypt.HashPassword(sess)

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

	getSession := database.Connection.Model(model.SessionCollection)
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

	match := encrypt.CheckPasswordHash(sess, session.Key)
	if !match {
		return errors.New("Credentials are invalid")
	}

	return nil
}

func SessionFindByID(SID string) (model.Session, error) {

	getSession := database.Connection.Model(model.SessionCollection)
	session := &model.Session{}

	err := getSession.FindId(bson.ObjectIdHex(SID)).Exec(session)
	if _, ok := err.(*mongodm.NotFoundError); ok {
		return model.Session{}, errors.New("Session not found")
	} else if err != nil {
		panic(err)
	}

	return *session, nil
}

func DeleteSession(id string) {

	getSession := database.Connection.Model(model.SessionCollection)
	err := getSession.RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		panic(err)
	}
}
