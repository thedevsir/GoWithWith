package model

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Gommunity/GoWithWith/config/database"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2/bson"
)

const AuthAttemptCollection = "AuthAttempt"

type (
	AuthAttempt struct {
		mongodm.DocumentBase `json:",inline" bson:",inline"`

		IP       string `json:"ip" bson:"ip"`
		Username string `json:"username" bson:"username"`
	}
	AbuseDetected struct {
		MaxIP            string
		MaxIPAndUsername string
	}
)

func (a *AbuseDetected) Check(ip, username string) error {

	anHourLater := time.Now().Add(-1 * time.Hour)
	authAttemptModel := database.Connection.Model("AuthAttempt")

	numIP, _ := authAttemptModel.Find(
		bson.M{
			"ip": ip,
			"createdAt": bson.M{
				"$gt": anHourLater,
			},
		},
	).Count()
	numIPUsername, _ := authAttemptModel.Find(
		bson.M{
			"ip":       ip,
			"username": username,
			"createdAt": bson.M{
				"$gt": anHourLater,
			},
		},
	).Count()

	maxIP, _ := strconv.Atoi(a.MaxIP)
	maxIPUsername, _ := strconv.Atoi(a.MaxIPAndUsername)

	ipLimitReached := numIP >= maxIP
	ipUsernameLimitReached := numIPUsername >= maxIPUsername

	if ipLimitReached || ipUsernameLimitReached {
		return errors.New("Maximum number of auth attempts reached")
	}

	return nil
}

func (a *AuthAttempt) Create() error {

	authAttemptModel := database.Connection.Model("AuthAttempt")
	attempt := &AuthAttempt{}
	authAttemptModel.New(attempt)

	attempt.IP = a.IP
	attempt.Username = strings.ToLower(a.Username)

	err := attempt.Save()

	if err != nil {
		panic(err)
	}

	return nil
}
