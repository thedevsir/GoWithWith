package repository

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Gommunity/GoWithWith/config/database"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2/bson"
)

type AuthAttempt struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`

	IP       string `json:"ip" bson:"ip"`
	Username string `json:"username" bson:"username"`
}

type AbuseDetected struct {
	MaxIP            string
	MaxIPAndUsername string
}

func (Config *AbuseDetected) Check(ip, username string) error {

	lastHour := time.Now().Add(-1 * time.Hour)
	getAuthAttempt := database.Connection.Model("AuthAttempt")

	numIP, _ := getAuthAttempt.Find(
		bson.M{
			"ip": ip,
			"createdAt": bson.M{
				"$gt": lastHour,
			},
		},
	).Count()
	numIPAndUsername, _ := getAuthAttempt.Find(
		bson.M{
			"ip":       ip,
			"username": username,
			"createdAt": bson.M{
				"$gt": lastHour,
			},
		},
	).Count()

	maxIP, _ := strconv.Atoi(Config.MaxIP)
	maxIPAndUsername, _ := strconv.Atoi(Config.MaxIPAndUsername)

	ipLimitReached := numIP >= maxIP
	ipUserLimitReached := numIPAndUsername >= maxIPAndUsername

	if ipLimitReached || ipUserLimitReached {
		return errors.New("Maximum number of auth attempts reached")
	}

	return nil
}

func AttemptCreate(ip, username string) error {

	getAuthAttempt := database.Connection.Model("AuthAttempt")
	attempt := &AuthAttempt{}
	getAuthAttempt.New(attempt)

	attempt.IP = ip
	attempt.Username = strings.ToLower(username)

	err := attempt.Save()

	if err != nil {
		panic(err)
	}

	return nil
}
