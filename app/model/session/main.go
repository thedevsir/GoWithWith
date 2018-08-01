package session

import (
	"errors"
	"time"

	"github.com/Gommunity/GoWithWith/config/database"
	"github.com/Gommunity/GoWithWith/services/encrypt"
	"github.com/Gommunity/GoWithWith/services/paginate"
	"github.com/satori/go.uuid"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2/bson"
)

const SessionCollection = "Session"

type Session struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`

	IP        string `json:"ip" bson:"ip"`
	Key       string `json:"key" bson:"key"`
	UserID    string `json:"userId" bson:"userId"`
	UserAgent string `json:"userAgent" bson:"userAgent"`
}

func (s *Session) Create() (string, string) {

	sessionModel := database.Connection.Model(SessionCollection)
	session := &Session{}
	sessionModel.New(session)
	uuid := uuid.Must(uuid.NewV4(), nil).String() // nil
	hash, _ := encrypt.Hash(uuid)

	session.IP = s.IP
	session.Key = hash
	session.UserID = s.UserID
	session.UserAgent = s.UserAgent

	err := session.Save()
	if err != nil {
		panic(err)
	}

	return session.GetId().Hex(), uuid
}

func SessionFindByCredentials(Session, SID string) error {

	session, err := SessionFindByID(SID)

	if err != nil {
		return err
	}

	match := encrypt.CheckHash(Session, session.Key)
	if !match {
		return errors.New("Credentials are invalid")
	}

	return nil
}

func SessionFindByID(SID string) (Session, error) {

	sessionModel := database.Connection.Model(SessionCollection)
	session := &Session{}

	err := sessionModel.FindId(bson.ObjectIdHex(SID)).Exec(session)
	_, ok := err.(*mongodm.NotFoundError)
	switch {
	case ok:
		return Session{}, errors.New("Session was not found")
	case err != nil:
		panic(err)
	default:
		return *session, nil
	}
}

func SessionUpdateLastActivity(SID string) {

	sessionModel := database.Connection.Model(SessionCollection)
	update := bson.M{
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}

	err := sessionModel.UpdateId(bson.ObjectIdHex(SID), update)
	if err != nil {
		panic(err)
	}
}

func GetUserSessions(userID string, page, limit int) (paginate.Paginate, error) {

	sessionModel := database.Connection.Model(SessionCollection)
	sessions := []*Session{}
	result := sessionModel.Find(bson.M{"userId": userID}).Select(bson.M{"key": 0, "userId": 0}).Sort("updatedAt").Skip((page - 1) * limit).Limit(limit)

	count, _ := result.Count()
	err := result.Exec(&sessions)

	_, ok := err.(*mongodm.NotFoundError)
	switch {
	case ok:
		return paginate.Paginate{}, errors.New("No data was found")
	case err != nil:
		panic(err)
	}

	pagination := paginate.Generate(sessions, count, page, limit)
	return pagination, nil
}

func DeleteSession(ID string) {

	sessionModel := database.Connection.Model(SessionCollection)
	err := sessionModel.RemoveId(bson.ObjectIdHex(ID))

	if err != nil {
		panic(err)
	}
}
